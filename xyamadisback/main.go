package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis"
)

const TaskPool string = "transaction_distribution"
const UpStream string = ":message_from_emv_kernel"
const DownStream string = ":message_from_smart_card"
const Terminate string = "TERMINATE"
const commandIdStart uint32 = 0
const VisaAID string = "A0000000031010"
const MastercardAID string = "A0000000041010"

var taskId string

type Transaction struct {
	TransactionId uint64  	       `json:"transactionId"`
	SupportedCardNetwork []string  `json:"supportedCardNetwork"`
	PpseResponse string 	       `json:"ppseResponse"`
	AidResponse string             `json:"selectAidResponse"`
	InputTransactionData  []Tlv    `json:"inputTransactionData"`
	ExpectedTransactionTags []uint32 `json:"expectedTransactionTags"`
}

type Message struct {
	ActionName string   `json:"actionName"`
	ActionData string   `json:"actionData"`
}

type Apdus struct {
	Apdus []Apdu    `json:"apdus"`
}

type Apdu struct {
	CommandId uint32    `json:"commandId"`
	ApduCommand string `json:"apduCommand"`
}

type KernelError struct {
	ErrorCode string   `json:"errorCode"`
}

type Tlvs struct {
	Tlvs []Tlv      `json:"tlvs"`
}

type Tlv struct {
	Tag uint32       `json:"tag"`
	Value string    `json:"value"`
}

func main() {

	// Define a boolean flag
	isMastercard := flag.Bool("mastercard", false, "Enable mastercard mode")

	// Parse the command-line arguments
	flag.Parse()

	// Create a new Redis client using the default options.
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378", //"44.230.27.125:6379", //"localhost:6378", //"xyredis.kkntoc.ng.0001.use2.cache.amazonaws.com:6379",
		Password: "",               // leave blank if not using authentication
		DB:       0,                // use default database
	})

	// Ping the Redis server to ensure that the connection is working.
	err := client.Ping().Err()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Redis is available")
	}

	// Create a channel to receive signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Loop forever until a signal is received
	for {
		fmt.Println("Looping forever...")
		select {
		case <-sigs:
			fmt.Println("Got signal, terminating loop...")
			client.Close()
			return
		default:
			initData, _ := client.BRPop(10*time.Second, TaskPool).Result()
			// actually one BRPop only returns one message
			if len(initData) > 0 {
				finalInitData := strings.TrimSpace(initData[1])
				if len(finalInitData) > 0 {
					fmt.Println("Receive a task: ", finalInitData)
					doTransaction(client, finalInitData, *isMastercard)
					time.Sleep(200 * time.Microsecond)
				}
			}
		}
	}
}

func doTransaction(client *redis.Client, finalInitData string, isMastercard bool) {
	// JSON format need parse
	var transaction Transaction
	json.Unmarshal([]byte(finalInitData), &transaction)
	fmt.Printf("Transaction Id is: %d\n", transaction.TransactionId)

	if strings.Contains(transaction.PpseResponse, MastercardAID) {
		isMastercard = true
		fmt.Println("Mastercard case")
	} else {
		fmt.Println("Visa card case")
	}
	fmt.Println()
	taskId = strconv.FormatUint(transaction.TransactionId, 10)

	upstreamChannel := taskId + UpStream
	downstreamChannel := taskId + DownStream

	//simulate handle select PPSE response
	time.Sleep(12 * time.Microsecond)

	if strings.TrimSpace(transaction.AidResponse) == "" {
		// return with select AID command
		if isMastercard {
			client.LPush(upstreamChannel, composeApduCommand("00A4040007" + MastercardAID + "00"))
		} else {
			client.LPush(upstreamChannel, composeApduCommand("00A4040007" + VisaAID + "00"))
		}
		//client.LPush(upstreamChannel, "00A4040007A000000003101000")
		client.Expire(upstreamChannel, 10*time.Second).Result()

		// wait for select AID command response, need finish fast
		selectAidResp, _ := client.BRPop(1*time.Second, downstreamChannel).Result()
		if len(selectAidResp) == 0 {
			fmt.Println("time out during select AID command")
			client.LPush(upstreamChannel, composeErrorCommand(Terminate))
			return
		}
		finalSelectAidResp := strings.TrimSpace(selectAidResp[1])
		if len(finalSelectAidResp) == 0 {
			client.LPush(upstreamChannel, composeErrorCommand(Terminate))
			fmt.Println("no select AID response data found")
			return
		}

		// No need parse
		fmt.Println(finalSelectAidResp)

		// simulate handle select AID response
		time.Sleep(2 * time.Microsecond)
	}

	// return with get GPO command
	if isMastercard {
		client.LPush(upstreamChannel, composeApduCommand("80A8000002830000"))
	} else {
		client.LPush(upstreamChannel, composeApduCommand("80A80000238321F0204000000000000150000000000000084000000000000840230306001212121200"))
		//client.LPush(upstreamChannel, "80A80000238321F0204000000000000150000000000000084000000000000840230306001212121200") //"80A800000D830B000000000000000008402200")
	}

	// wait for get GPO command response, need finish fast
	getGpoResp, _ := client.BRPop(1*time.Second, downstreamChannel).Result()
	if len(getGpoResp) == 0 {
		fmt.Println("time out during get GPO command")
		client.LPush(upstreamChannel, composeErrorCommand(Terminate))
		return
	}
	finalGetGpoResp := strings.TrimSpace(getGpoResp[1])
	if len(finalGetGpoResp) == 0 {
		client.LPush(upstreamChannel, composeErrorCommand(Terminate))
		fmt.Println("no get GPO response data found")
		return
	}

	// no need parse
	fmt.Println(finalGetGpoResp)

	if isMastercard {
		// simulate handle get GPO response
		time.Sleep(2 * time.Microsecond)

		// return with generate AC command
		client.LPush(upstreamChannel, composeApduCommand("80AE50004200000010188600000000000008400000000000097823051700373EB7BD00000000000000000000001E00001725060000000000000000000000000000000000000000"))

		// wait for get GPO command response, need finish fast
		generateAcResp, _ := client.BRPop(1*time.Second, downstreamChannel).Result()
		if len(generateAcResp) == 0 {
			fmt.Println("time out during generate AC command")
			client.LPush(upstreamChannel, composeErrorCommand(Terminate))
			return
		}
		finalGenerateAcResp := strings.TrimSpace(generateAcResp[1])
		if len(finalGenerateAcResp) == 0 {
			client.LPush(upstreamChannel, composeErrorCommand(Terminate))
			fmt.Println("no generate AC response data found")
			return
		}

		// no need parse
		fmt.Println(finalGenerateAcResp)
	}

	// together with emv tags
	client.LPush(upstreamChannel, composeFinishCommand())
	// do some cleanups
}

func composeFinishCommand() string {
	tlvList := []Tlv{
		{Tag: 40730, Value: "0840"},
		{Tag: 156, Value: "00"},
	}
	tlvs := Tlvs{
		Tlvs: tlvList,
	}
	jsonBytes, _ := json.Marshal(tlvs)
	message := Message{
		ActionName: "TRANSACTION_DONE",
		ActionData: string(jsonBytes),
	}
	jsonMsgBytes, _ := json.Marshal(message)
	fmt.Println(string(jsonMsgBytes))
	return string(jsonMsgBytes)
}

func composeErrorCommand(errorCode string) string {
	kernelError:= KernelError{
		ErrorCode: errorCode,
	}
	jsonBytes, _ := json.Marshal(kernelError)
	message := Message{
		ActionName: "KERNEL_ERROR",
		ActionData: string(jsonBytes),
	}
	jsonMsgBytes, _ := json.Marshal(message)
	fmt.Println(string(jsonMsgBytes))
	return string(jsonMsgBytes)
}

func composeApduCommand(cmd string) string {
	apduList := []Apdu{
		{CommandId: commandIdStart, ApduCommand: cmd},
	}
	apdus := Apdus{
		Apdus: apduList,
	}
	jsonBytes, _ := json.Marshal(apdus)
	message := Message{
		ActionName: "EXECUTE_APDU",
		ActionData: string(jsonBytes),
	}
	jsonMsgBytes, _ := json.Marshal(message)
	fmt.Println(string(jsonMsgBytes))
	return string(jsonMsgBytes)
}
