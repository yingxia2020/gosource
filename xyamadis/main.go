package main

import (
	"encoding/json"
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
const DownStream string = ":message_for_emv_kernel"
const Terminate string = "TERMINATE"

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

type Event struct {
	EventName string   `json:"eventName"`
	EventData string   `json:"eventData"`
}

type ExecuteTransaction struct {
	PaymentEntryPointParameter string `json:"paymentEntryPointParameter"`
	Configuration              string `json:"configuration"`
}

type UpdateTransaction struct {
	AuthorizedAmount           string `json:"authorizedAmount"`
}

type KernelError struct {
	ErrorCode string   `json:"errorCode"`
}

type Tlvs struct {
	Tlvs []Tlv      `json:"tlvs"`
}

type Tlv struct {
	Tag uint32      `json:"tag"`
	Value string    `json:"value"`
}

func main() {

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
		fmt.Println("waiting for transaction...")
		select {
		case <-sigs:
			fmt.Println("Got signal, terminating loop...")
			client.Close()
			return
		default:
			initData, _ := client.BRPop(10*time.Second, TaskPool).Result()
			// actually one BRPop only returns one message
			if len(initData) > 0 {
				finalData := strings.TrimSpace(initData[1])
				if len(finalData) > 0 {
					fmt.Println("Receive a task: ", finalData)
					doTransaction(client, finalData)
					time.Sleep(200 * time.Microsecond)
				}
			}
		}
	}
}

func handleMessage(client *redis.Client, messageData string, upstreamChannel string) bool {
	var event Event
	json.Unmarshal([]byte(messageData), &event)

	if event.EventName == "TRANSACTION_UPDATED" {
		fmt.Println("Received TRANSACTION_UPDATED message: " + event.EventData)
		var updateTrans UpdateTransaction
		json.Unmarshal([]byte(event.EventData), &updateTrans)
		client.LPush(upstreamChannel, composeExecuteTransactionCommand(true, updateTrans.AuthorizedAmount))
		return false
	} else if event.EventName == "TRANSACTION_EXECUTED" {
		fmt.Println(event.EventData)
		if strings.Contains(event.EventData, "TRANSACTION_DONE") {
			fmt.Println("Transaction is finished successfully")
			client.LPush(upstreamChannel, composeFinishCommand())
		} else { // TRANSACTION_CANCELED
			fmt.Println("Something goes wrong...")
			client.LPush(upstreamChannel, composeErrorCommand("USER_CANCEL_OBSOLETE"))
		}
		return true
	} else if event.EventName == "TRANSACTION_EXECUTED" {
		fmt.Println(event.EventData)
		fmt.Println("Transaction is canceled by user")
		return true
	}

	// error condition
	fmt.Println("Received TRANSACTION_CANCELED message: " + event.EventData)
	return true
}

func doTransaction(client *redis.Client, finalInitData string) {
	// JSON format need parse
	var transaction Transaction
	json.Unmarshal([]byte(finalInitData), &transaction)
	fmt.Printf("Transaction Id is: %d\n", transaction.TransactionId)

	taskId = strconv.FormatUint(transaction.TransactionId, 10)

	upstreamChannel := taskId + UpStream
	downstreamChannel := taskId + DownStream

	// send back first payment context message
	client.LPush(upstreamChannel, composeExecuteTransactionCommand(false, ""))

	//simulate handle select PPSE response
	//time.Sleep(12 * time.Microsecond)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Loop forever until a signal is received
	outerloop:
	for {
		fmt.Println("waiting for messages...")
		select {
		case <-sigs:
			fmt.Println("Got signal, terminating loop...")
			client.Close()
			return
		default:
			msgData, _ := client.BRPop(10*time.Second, downstreamChannel).Result()
			// actually one BRPop only returns one message
			if len(msgData) > 0 {
				finalData := strings.TrimSpace(msgData[1])
				if len(finalData) > 0 {
					fmt.Println("Receive a message: ", finalData)
					if handleMessage(client, finalData, upstreamChannel) {
						break outerloop
					}
					time.Sleep(200 * time.Microsecond)
				}
			}
		}
	}
}

func composeExecuteTransactionCommand(update bool, authorizedAmount string) string {
	var executeTrans ExecuteTransaction
	if update {
		executeTrans = ExecuteTransaction {
			Configuration: "",
			PaymentEntryPointParameter: "TRANSACTION_CONTINUE",
		}
	} else {
		executeTrans = ExecuteTransaction{
			Configuration: "",
			PaymentEntryPointParameter: "TRANSACTION_START",
		}
	}

	jsonBytes, _ := json.Marshal(executeTrans)
	message := Message{
		ActionName: "EXECUTE_TRANSACTION",
		ActionData: string(jsonBytes),
	}
	jsonMsgBytes, _ := json.Marshal(message)
	fmt.Println(string(jsonMsgBytes))
	return string(jsonMsgBytes)
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
