// Example static file server.
//
// Serves static files from the given directory.
// Exports various stats at /stats .
//https://github.com/tcnksm-sample/sarama/blob/master/LICENSE
package main

import (
	"bytes"
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"gopkg.in/Shopify/sarama.v2"
	//"io/ioutil"
	"log"
	"math"
	"math/rand"
	//"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"github.com/stretchr/stew/slice"
	//"strings"
	//"sync"
	"time"
	//"sync"
	//"os/exec"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
	//"context"
	//"github.com/prometheus/client_golang/api"
	//v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	//above three required for prometheus client
	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
)
//globals for timing and logging
var (
	videoFileName string = ""
	startTest int = 0
	stopTest int = 0
    configurationId string  = ""
    txMode string =""
	ffmpeg_args string=""
	config_name string = ""
    nNumberOfPods int = 0
	numberOfPods string = ""
	numberOfCpus string = ""
	category string = ""
	encoderToUse string = ""
	preset string = ""
	outputprotocol string = ""
	labelforvodnodes string = ""
	videoLibCount int = 0
	videonames [] string
	numUniqueVideos int = 0

	sentMsgsKeySlice [] string
	recdMsgsKeySlice [] string
	input_format string = ""
	abr_profile string = ""
    inputVidDur int = 0
	inputVidFPS int = 0
	totalFramesPerVideo int =0

	//to handle videos and their storage
	minio_ip_address string =""
	reqVideosArrstr string = ""

    testinProgress bool = false
	resultJson string = ""
	totalTestResultBuf = bytes.Buffer{}
    resultBuf = bytes.Buffer{}
    msgsStartTime = time.Now()
	msgsEndTime = time.Now()

	logBuf = bytes.Buffer{}
	runError bool = false
	runErrorDesc string = ""

	//to catch ffmpeg freezes
	lastMsgTime = time.Now()
	stopTestCh = make(chan bool)
	waitTimeSecs int = 1800
	waitTimeSecsOnePod int = 20000

	//wg sync.WaitGroup
)
//////////////////////////
//vod const
const kafkabroker string = "broker:9092"
const transcodeTopic string = "testTopicTranscodeFinal1"
//const promServerUrl string  = "http://10.233.13.12:8080" // nodeport needs to be setup for service using prometheus-service.yaml

type record struct {
	Name   string  `json:"name"`
	Addr   string  `json:"addr"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
type RequestConfig struct {
	ConfID      string `json:"conf_id"`
	GlobalID    string `json:"global_id"`
	Name        string `json:"name"`
	CodecStd    string `json:"codec_std"`
	Encoder     string `json:"encoder"`
	Preset      string `json:"preset"`
	InputFormat string `json:"input_format"`
	Abr         string `json:"abr"`
	FfmpegArgs  string `json:"ffmpeg_args"`
	Videos      []struct {
		VideoID string `json:"video_id"`
		Name    string `json:"name"`
		Fps     string `json:"fps"`
		Dur     string `json:"dur"`
	} `json:"videos"`
	Category 		 string `json:"category"`
	Mode             string `json:"mode"`
	NumVidstreams    string `json:"num_vidstreams"`
	CpusPerPod       string `json:"cpusPerPod"`
	Labelforvodnodes string `json:"labelforvodnodes"`
	Starttest        string `json:"starttest"`
	Stoptest         string `json:"stoptest"`
}

type RequestPodInfo struct {
	NumVidstreams    string `json:"num_vidstreams"`
	CpusPerPod       string `json:"cpusPerPod"`
	Labelforvodnodes string `json:"labelforvodnodes"`
	MinioIP          string `json:"minioIP"`
	Videos           []struct {
		VideoID string `json:"video_id"`
		Name    string `json:"name"`
		Fps     string `json:"fps"`
	} `json:"videos"`
}

type TranscodeResult struct {
	Key string
	Starttime string
	Endtime string
	Duration string
	Pod_id string
	Error bool
	Errordesc string
}
var (
	addr               = flag.String("addr", ":8070", "TCP address to listen to")
	addrTLS            = flag.String("addrTLS", ":8443", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS")
	byteRange          = flag.Bool("byteRange", false, "Enables byte range requests if set to true")
	certFile           = flag.String("certFile", "./cnbserver.crt", "Path to TLS certificate file")
	compress           = flag.Bool("compress", false, "Enables transparent response compression if set to true")
	dir                = flag.String("dir", "/usr/share/fileserver", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", true, "Whether to generate directory index pages")
	keyFile            = flag.String("keyFile", "./cnbserver.key", "Path to TLS key file")
	vhost              = flag.Bool("vhost", false, "Enables virtual hosting by prepending the requested path with the requested hostname")
)

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)

const (
	RANDOMTEXT = `bfeWQobe:9oY+vO.DjRV:V@JBZ)Fvj5UqWg?BrppM:'u0/y[cgi9_<,L4I8mQ>\
13fnp|@zH{'R7_VFuZ4M3Om@hMCh@suVsxA(3msi2oVS{kc)hTyc[#tRBZ"isOH
(d%nWAR48*rqFKEQH&pWIgl*DLMpbXah3C]T4[|Rq\@{6w"b<Q'<i\D,||W"ECx
!sCmnp,qeNK^46L1tE8n12"1Z,<gGh-2!R\*y,|<Mcl-01v$U\w*09qW9o>Bn+q
WWw+yW0jp7SEN8HE2Dr.*yv_yFC'P3B'$f6&)8j%G7lqR[Fn0dKw^I'-'CH\gIB
<8:zK0t<U&bzh2E)2T"|fR*:*qPjYy|!u@nq]Ch;r,=ddh?t@-,f,p&|g,@C'W1
q*G7X!pBci_NEfR)7NMqmkj=Ai8K1zgTH/IA9JgesNNXe81>^{)[Yf)[Y>+5;43
9B\':)5p3^{;6R,-vxo9dkH$I2bNI#62,'+,<OgG*)iSJ3/pv^xF2,l2tz!!cx{
:.V$SD?*ePM';x7).p0,7q>S#07{9J6pn0cP3P_-nv__Jyf;Kzz[GhXI8Ci\2?K
4h4|;>!PD'^H?F^2,QgX&k.6SgCJTM>KMnqS!hC?8s+N"O.EB'*2'$I=O|p.qmg
-.p>8DXn,CFP#w^J4?xrx${1@W)xoww9NMwXSbQ:\:^{"{vtVZ9dcV+X/S&N\61`

	genPageTitle = "General Page Title"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Parse command-line flags.
	flag.Parse()

	// Setup FS handler
	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		Compress:           *compress,
		AcceptByteRange:    *byteRange,
	}
	if *vhost {
		fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
	}
	fsHandler := fs.NewRequestHandler()

	// Create RequestHandler serving server stats on /stats and files
	// on other requested paths.
	// /stats output may be filtered using regexps. For example:
	//
	//   * /stats?r=fs will show only stats (expvars) containing 'fs'
	//     in their names.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stats":
			expvarhandler.ExpvarHandler(ctx)
		case "/vod":
			vodHandler(ctx)
		case "/vodStartTest":
			vodStartTestHandler(ctx)
		case "/vodStartTranscode":
			vodStartTranscodeHandler(ctx)
		case "/vodStatus":
			vodStatusHandler(ctx)
		case "/vodDeleteTopics":
			vodDeleteTopicHandler(ctx)
		case "/general":
			generalHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx)
		}
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		Concurrency: fasthttp.DefaultConcurrency,
	}
	// s.DisableKeepalive = true

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		go func() {
			if err := s.ListenAndServe(*addr); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}

	// Start HTTPS server.
	if len(*addrTLS) > 0 {
		log.Printf("Starting HTTPS server on %q", *addrTLS)
		go func() {
			if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, requestHandler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
	}

	log.Printf("Serving files from directory %q", *dir)
	log.Printf("See stats at http://%s/stats", *addr)

	// Wait forever.
	select {}
}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

func generalHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s\n", genPageTitle)

	round := rand.Intn(5) + 5
	for i := 0; i < round; i++ {
		fmt.Fprint(ctx, RANDOMTEXT)
	}
	fmt.Fprint(ctx, "\n")
}
func vodHandler(ctx *fasthttp.RequestCtx) {
	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(ctx, "Error response: %s\n", err.Error())
		return
	}
	fmt.Fprintf(ctx, "Hello from %v!\n\n", host)
	//resultsFolder = string(ctx.QueryArgs().Peek("resfolder"))
	//fmt.Fprintf(ctx, "Results folder: %v\n\n", resultsFolder)
	//to clean up
	var kafkaServer, kafkaTopic string
	kafkaServer =  "broker:9092"
	kafkaTopic  = "testTopicTranscodeFinal1"
	fmt.Fprintf(ctx, "Connecting to kafka broker at %v!\n\n" , kafkaServer  + ":" + kafkaTopic)
	//usingClientGo()
}
/*func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Printf(scanner.Text())
		fmt.Println(scanner.Text())
	}
}*/
func vodStatusHandler(ctx *fasthttp.RequestCtx){
	log.Printf("Test Status - step1")
	log.Printf("testinProgress5:%v \n", testinProgress)
	if testinProgress == true {
		if len(resultJson) > 0 { // this happens when one pod has returned sucessfully
			log.Printf("Test in Progress:%v\n", videoFileName)
			fmt.Fprintf(ctx, "Test in Progress:%v\n", videoFileName)
			if nNumberOfPods > 1 {
				//if there are multiple pods then the expectation is that they will return with similar durations. We wait waitTimeSecs seconds for all pods to be done
				//if all pods are not done with this extra time, then most probably the pod could be hanging. stop the test
				curTime := time.Now()
				durSinceLastSucMsg := curTime.Sub(lastMsgTime).Seconds()
				log.Printf("DURATION: %v :\n", durSinceLastSucMsg)
				//for testing setting it to 10 mins ( 600 secs)
				if (int(durSinceLastSucMsg) > waitTimeSecs) {
					testVar := lastMsgTime.Format("2006-01-02 15:04:05")
					log.Printf("vodStatusHandler lastMsgTime with more than one pod:%v\n", testVar)
					curTimeStr := curTime.Format("2006-01-02 15:04:05")
					log.Printf("vodStatusHandler curTime with more than one pod:%v\n", curTimeStr)
					log.Printf("vodStatusHandler durSinceLastSucMsg with more than one pod:%v\n", durSinceLastSucMsg)
					stopTestCh <- true
				}
			}
		} else { //resultJson is not populated so not even one pod has returned successfully
			//to catch ffmpeg and/or pod freezes or unresponsiveness
			//If only one pod was  launched then if there has been no response after a preset amount of time, stop the test
			//Wait for waitTimeSecsOnePod and then stop the test
			if nNumberOfPods == 1 {
				curTime := time.Now()
				durSinceLastSucMsg := curTime.Sub(lastMsgTime).Seconds()
				if (int(durSinceLastSucMsg) > waitTimeSecsOnePod) {
					testVar := lastMsgTime.Format("2006-01-02 15:04:05")
					log.Printf("vodStatusHandler lastMsgTime with one pod:%v\n", testVar)
					curTimeStr := curTime.Format("2006-01-02 15:04:05")
					log.Printf("vodStatusHandler curTime with one pod:%v\n", curTimeStr)
					log.Printf("vodStatusHandler durSinceLastSucMsg with one pod :%v\n", durSinceLastSucMsg)

					stopTestCh <- true
				}
				log.Printf("Test Status - step4")
				fmt.Fprintf(ctx, "Test in Progress:%s\n", videoFileName)
			}else{
				log.Printf("Test Status - step6")
				fmt.Fprintf(ctx, "Test in Progress:%s\n", videoFileName)
			}
		}
	}else{
		log.Printf("Test Status - step2")
		log.Printf("LOG totalTestResultBuf%s\n", totalTestResultBuf )
		fmt.Fprintf(ctx,"%s\n", totalTestResultBuf.String())
		fmt.Fprintf(ctx,"%s\n", resultBuf.String())
	}
	log.Printf("Test Status - step3")
	ctx.SetStatusCode(fasthttp.StatusOK)
}
func vodCreateTopic() bool{
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	brokers := []string{"broker:9092"}
	admin, err := sarama.NewClusterAdmin(brokers, config) //c.Kafka.Brokers is of type: []string
	if err != nil {
		fmt.Println("error is", err)
		return false
	}
	/*detail := sarama.TopicDetail{NumPartitions: 15, ReplicationFactor: 1, ConfigEntries: map[string]*string{
		"retention.ms": "259200000", //3 days
	}}*/
	detail := sarama.TopicDetail{NumPartitions: 50, ReplicationFactor: 1}
	err = admin.CreateTopic("testTopicTranscodeFinal1",&detail,true)
	if err != nil {
		fmt.Println("error is", err)
		return false
	}else{
		return true
	}

}
func deleteTopic() bool{
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0

	brokers := []string{"broker:9092"}
	admin, err := sarama.NewClusterAdmin(brokers, config) //c.Kafka.Brokers is of type: []string
	if err != nil {
		fmt.Println("error is", err)
		return false
	}

	err = admin.DeleteTopic("testTopicTranscodeFinal1")
	if err != nil {
		fmt.Println("error is", err)
		return false
	}

	err = admin.Close()
	if err != nil {
		fmt.Println("error is", err)
		return false
	}

	return true
}
func vodDeleteTopicHandler(ctx *fasthttp.RequestCtx){
	ret :=deleteTopic()
	if ret  {
		log.Printf("Deleted topic")
	}else{
		log.Printf("Failed to delete topic")
	}
	//delete topic folder?
}
func vodStartTestHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("In vodStartTestHandler1")
	if len(videonames) > 0 {
		videonames = nil
	}

	var podReqConf RequestPodInfo
	if ctx.IsPost() {
		logBuf.WriteString("Postreceived" + "\n")
		log.Printf("In Postreceived1")
		json.Unmarshal(ctx.PostBody(), &podReqConf)

		t := time.Now()
		testVar := t.Format("2006-01-02 15:04:05")

		logBuf.WriteString("vodStartTestHandler:" + testVar + "\n")
		log.Printf("In vodStartTestHandler")
		log.Printf("testinProgress0:%v \n", testinProgress)
		testinProgress = true
		log.Printf("testinProgress1:%v \n", testinProgress)

		ret := vodCreateTopic() //to do handle errors
		if ret {
			log.Printf("Created topic")
			time.Sleep(30 * time.Second)
		} else {
			log.Printf("Failed to create topic")
		}

		//currently ffmpeg args are pass through from vod.sh, with config name and id
		numberOfPods = string(podReqConf.NumVidstreams)
		nPods, err := strconv.Atoi(numberOfPods)
		nNumberOfPods=nPods
		if err != nil {
			fmt.Println(err)
		}
		numberOfCpus = string(podReqConf.CpusPerPod)
		labelforvodnodes = string(podReqConf.Labelforvodnodes)
		videoLibCount, _ = strconv.Atoi(string(podReqConf.NumVidstreams))
		minio_ip_address = string(podReqConf.MinioIP)

		vidTest := podReqConf.Videos

		fmt.Println(reflect.TypeOf(podReqConf.Videos))
		fmt.Println(podReqConf.Videos)


		vidcount:=0
		for _, vid := range vidTest {
			if vidcount == 0{
				reqVideosArrstr+=vid.Name
			}else{
				reqVideosArrstr = reqVideosArrstr + "," + vid.Name
			}
			vidcount++
		}
		log.Printf("reqVideosArrstr:%s", reqVideosArrstr)
		log.Printf("numberOfPods: %s\n", numberOfPods)
		log.Printf("videoLibCount: %d", videoLibCount)

		cmd := exec.Command("python3", "/root/app/changecluster.py", "create", numberOfPods, numberOfCpus, labelforvodnodes, minio_ip_address, reqVideosArrstr) //Run python script with function to run, in case of create number of pods and number of cpus per pod
		log.Printf("%s\n", "Running")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// Wait for the Python program to exit.
		errCmd := cmd.Run()
		log.Printf("Finished:%s", errCmd)
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", errCmd)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		log.Printf("outStr:%s", string(stdout.Bytes()))
		log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
		if string(stdout.Bytes()) == "error" {
			log.Printf("Status code: service Unavailable")
			ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		} else {
			time.Sleep(60 * time.Second)
			//videoFileNames[0] = "ToS_2160p_59.94fps_2min.mp4"
			//videoFileNames[1] = "ToS_2160p_59.94fps_2min_2.mp4"//default as only one video file is supported
			msgsStartTime = time.Now()
			t := time.Now()
			testVar := t.Format("2006-01-02 15:04:05")
			logBuf.WriteString("before start:" + testVar + "\n")
			log.Printf("%s\n", "Before setting status to Accepted")
			ctx.SetStatusCode(fasthttp.StatusAccepted)
		}

		ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
	}
}

func vodStartTranscodeHandler(ctx *fasthttp.RequestCtx){
	log.Printf("In vodStartTranscodeHandler")
	if len(videonames)> 0{
		videonames = nil
	}

	var reqConf RequestConfig
	if ctx.IsPost() {
		logBuf.WriteString("Postreceived" +  "\n")
		log.Printf("In Postreceived1")
		json.Unmarshal(ctx.PostBody(), &reqConf)

		t := time.Now()
		testVar := t.Format("2006-01-02 15:04:05")


		logBuf.WriteString("vodStartTestHandler:" + testVar + "\n")
		log.Printf("In vodStartTestHandler 1 ")
		testinProgress = true
		log.Printf("In vodStartTranscodeHandler 2 :%v \n", testinProgress)

		ret :=vodCreateTopic()//to do handle errors
		if ret  {
			log.Printf("Created topic")
			time.Sleep( 30 * time.Second)
		}else{
			log.Printf("Failed to create topic")
		}

		//currently ffmpeg args are pass through from vod.sh, with config name and id
		startTest, _ = strconv.Atoi(string(reqConf.Starttest))
		stopTest, _ = strconv.Atoi(string(reqConf.Stoptest))
		configurationId = string(reqConf.ConfID)
		txMode = string(reqConf.Mode)
		category  = string(reqConf.Category)
		//numberOfPods = string(reqConf.NumVidstreams)
		//numberOfCpus = string(reqConf.CpusPerPod)
		ffmpeg_args = string(reqConf.FfmpegArgs)
		config_name = string(reqConf.Name)
		vidTest := reqConf.Videos
		fmt.Println(reflect.TypeOf(reqConf.Videos))
		fmt.Println(reqConf.Videos)

		time.Sleep( 30 * time.Second)


		for _,vid := range vidTest {
			videonames = append(videonames, vid.Name)
			inputVidFPS,_ = strconv.Atoi(vid.Fps)
			inputVidDur, _ = strconv.Atoi(vid.Dur)
		}
		log.Printf("video names:")
		videoFileName = videonames[0]
		totalFramesPerVideo = inputVidDur * inputVidFPS

		numUniqueVideos = len(videonames)
		//to do : check config selected  by including mediaconfig.json here. Check conditionally if it is a  benchmark product or used as a tool
		//currently not using the following parameters. They are needed when checking
		encoderToUse = string(reqConf.Encoder)
		preset = string(reqConf.Preset)
		input_format = string(reqConf.InputFormat)
		abr_profile = string(reqConf.Abr)

		videoLibCount, _ = strconv.Atoi(string(reqConf.NumVidstreams))

		log.Printf("starttest: %d", startTest)
		log.Printf("stoptest: %d", stopTest)
		log.Printf("category: %s\n", category)
		log.Printf("numberOfPods: %s\n", numberOfPods)
		log.Printf("videoLibCount: %d", videoLibCount)
		log.Printf("input_formats: %s\n", input_format)
		log.Printf("abr_profile: %s\n", abr_profile)
		log.Printf("inputVidDur: %d\n", inputVidDur)
		log.Printf("inputVidFPS: %d\n", inputVidFPS)
		log.Printf("ffmpeg_args: %s\n", ffmpeg_args)


		msgsStartTime = time.Now()
		totalTestResultBuf = bytes.Buffer{}
		resultBuf = bytes.Buffer{}
		go startKafkaProducer(videoLibCount, videoFileName)
		go startKafkaConsumer(videoLibCount)


		/*if startTest == 0 {
			cmd := exec.Command("python3", "/root/app/changecluster.py", "create",numberOfPods,numberOfCpus,labelforvodnodes,minio_ip_address,reqVideosArrstr)//Run python script with function to run, in case of create number of pods and number of cpus per pod
			log.Printf("%s\n","Running")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// Wait for the Python program to exit.
			err := cmd.Run()
			log.Printf("Finished:%s", err)
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
			log.Printf("outStr:%s", string(stdout.Bytes()))
			log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
			if string(stdout.Bytes())=="error" {
				log.Printf("Status code: service Unavailable")
				ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
			}else{
				time.Sleep( 60 * time.Second)
				//videoFileNames[0] = "ToS_2160p_59.94fps_2min.mp4"
				//videoFileNames[1] = "ToS_2160p_59.94fps_2min_2.mp4"//default as only one video file is supported
				msgsStartTime = time.Now()
				t := time.Now()
				testVar := t.Format("2006-01-02 15:04:05")
				logBuf.WriteString("before start:" + testVar + "\n")
				go startKafkaProducer(videoLibCount, videoFileName)
				go startKafkaConsumer(videoLibCount)
				log.Printf("%s\n","Before setting status to Accepted")
				ctx.SetStatusCode(fasthttp.StatusAccepted)
			}
			ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
		}else{
			msgsStartTime = time.Now()
			totalTestResultBuf = bytes.Buffer{}
			resultBuf = bytes.Buffer{}
			go startKafkaProducer(videoLibCount, videoFileName)
			go startKafkaConsumer(videoLibCount)
		}*/
	}else{
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
}


/*func vodStartTranscodeHandler(ctx *fasthttp.RequestCtx){
	log.Printf("%s\n","In vodStartTranscode")
	videoNameStr := string(ctx.QueryArgs().Peek("videoname"))
	log.Printf("%s\n",videoNameStr)

	go startHarness(videoNameStr)

	ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
	return
	/*msgsStartTime = time.Now()
	t := time.Now()
	testVar := t.Format("2006-01-02 15:04:05")
	logBuf.WriteString("before start:" + testVar + "\n")
	log.Printf("%s\n","startHarness")
	for i:= 0; i< numUniqueVideos; i++ {
		wg.Add(2)
		go startKafkaProducer(videoLibCount, videoFileNames[i])
		go startKafkaConsumer(videoLibCount, i)
		vidcount++
		log.Printf("Before wait: %d",i)
		wg.Wait()
	}
	log.Printf("%s\n","Done after wait groups")
}*/
/*func startHarness(){
	msgsStartTime = time.Now()
	t := time.Now()
	testVar := t.Format("2006-01-02 15:04:05")
	logBuf.WriteString("before start:" + testVar + "\n")
	log.Printf("%s\n","startHarness")
	wrkDone := false
	//var vidName string = ""
	resultBuf.WriteString("status, key, Starttime,Endtime,Duration,Pod_id,Fps\n")
	for i:= 0; i< numUniqueVideos; i++ {
		wg.Add(2)
		go startKafkaProducer(videoLibCount, videonames[i] )
		if i==numUniqueVideos -1  {
			wrkDone = true
		}
		go startKafkaConsumer(videoLibCount,wrkDone)
		vidcount++
		log.Printf("Before wait: %d",i)
		wg.Wait()
	}
	log.Printf("%s\n","Done after wait groups")
}*/
/*func startKafkaConsumer(numMessages int, numVideosDone int){
	defer wg.Done()
	log.Printf("In startKafkaConsumer")

	var errors int
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"broker:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "TopicTranscodeDone1"
	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	//consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	quit := make(chan bool)

	// Count how many message processed
	msgCount := 0
	resultBuf.WriteString("status, key, Starttime,Endtime,Duration,Pod_id,Fps\n")
	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		log.Printf("In Go func")
		for {
			select {
			case err := <-consumer.Errors():
				errors++
				fmt.Println(err)
			case msg := <-consumer.Messages():
				t := time.Now()
				testVar := t.Format("2006-01-02 15:04:05")
				logBuf.WriteString("Receiver:" + testVar + "\n")
				msgCount++
				log.Printf("msgCount:%d", msgCount)
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))

				var tr TranscodeResult
				resultJson =  string(msg.Value)
				json.Unmarshal([]byte(resultJson),&tr)
				var dur=0.0
				if !tr.Error{
					if dur, err = strconv.ParseFloat(tr.Duration, 64); err != nil {
						fmt.Println("Error Converting dur:")
						fmt.Println(err)
					}
					fps := float64(float64(totalFramesPerVideo)/dur)
					sFps := strconv.FormatFloat(fps, 'f', 2, 64)
					resultBuf.WriteString("success," + tr.Key +  "," + tr.Starttime +  "," + tr.Endtime + "," + tr.Duration + "," + tr.Pod_id + "," + sFps+ "\n")
				}else{
					if !runError{
						runError = true
						runErrorDesc = runErrorDesc  + tr.Errordesc + ";"
					}

					resultBuf.WriteString("fail," + tr.Key +  "," + tr.Starttime +  "," + tr.Endtime + "," + tr.Duration + "," + tr.Pod_id + "," + "0" +  "\n")
				}


				recdMsgsKeySlice = append(recdMsgsKeySlice, string(tr.Key))
				log.Printf("resultBuf:%S", resultBuf)

				if msgCount == numMessages {
					testinProgress = false
					msgsEndTime = time.Now()
					timeDiff := msgsEndTime.Sub(msgsStartTime).Seconds()
					//Check if the correct messages were processed. Some time if the test fails and Kafka messages are not retired or cleaned up, old messages maybe processed.
					allMsgsPresent:= true
					log.Printf("allMsgsPresent" + strconv.FormatBool(allMsgsPresent))
					sentMsgsCount := len(sentMsgsKeySlice)
					log.Printf("sentMsgsKeySlice: %d ", sentMsgsCount)
					for _, sKey := range sentMsgsKeySlice {
						log.Printf("sentMsgsKeySlice Key: %s ", sKey)
					}


					recdMsgsCount :=len(recdMsgsKeySlice)
					log.Printf("recdMsgsCount: %d ", recdMsgsCount)
					for _, sKey := range recdMsgsKeySlice {
						log.Printf("recdMsgsKeySlice Key: %s ", sKey)
					}

					if(sentMsgsCount == recdMsgsCount) {
						//check the actual keys
						for _, sKey := range sentMsgsKeySlice {
							log.Printf("Key: %s ", sKey)
							if slice.Contains(recdMsgsKeySlice, sKey)==false{
								allMsgsPresent = false
								break
							}
						}
					}
					log.Printf("allMsgsPresent" + strconv.FormatBool(allMsgsPresent))
					if (allMsgsPresent) {
						if !runError{
							podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
							totalTestResultBuf.WriteString("Status," + "success" + "\n")
							totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
							totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
							totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
							var nNumberOfPods int = 0
							nNumberOfPods, err = strconv.Atoi(numberOfPods)
							if err != nil {
								fmt.Println(err)
							}

							aggregateFPS := float64(totalFramesPerVideo * nNumberOfPods) / timeDiff
							if numVideosDone == numUniqueVideos{
								totalTestResultBuf.WriteString("numVideosDone == numUniqueVideos :::" + strconv.FormatFloat(aggregateFPS,  'f', 2, 64) + "\n")
							}
							//totalTestResultBuf.WriteString("Aggregate FPS," + strconv.FormatFloat(aggregateFPS,  'f', 2, 64) + "\n")
						}else{
							podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
							totalTestResultBuf.WriteString("Status," + "fail," + runErrorDesc  + "\n")
							totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
							totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
							totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
						}

					}else{
						log.Printf("Error: Old kafka messages found\n")
						podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
						totalTestResultBuf.WriteString("Status," + "fail," + "Error Running the workload. Cleanup pods and Kafka. Install Kafka again to run the test," + "\n")
						totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
						totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
						totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
					}

					quit <- true

				}
			case <-quit:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Processed: %d; errors: %d\n", msgCount, errors)

}*/
//func startKafkaConsumer(numMessages int, runDone bool){
func startKafkaConsumer(numMessages int){
	log.Printf("In startKafkaConsumer")

	var errors int
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"broker:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "TopicTranscodeDone1"
	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	//consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	quit := make(chan bool)

	// Count how many message processed
	msgCount := 0
	log.Printf("In startKafkaConsumer msgCount %d\n", msgCount)
	//resultBuf.WriteString("status, key, Starttime,Endtime,Duration,Pod_id,Fps\n")
	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		log.Printf("In Go func")
		for {
			select {
			case err := <-consumer.Errors():
				errors++
				fmt.Println(err)
			case msg := <-consumer.Messages():
				//reset lastMsgTime after receiving a message from one pod when more than one pods are launched.
				//We are assuming that all pods will achive about the same amout of time as they are allocated similar resources
				//So if one pods is taking too much extra time. the probability is that ffmpeg process froze or the pod is not returning
				lastMsgTime = time.Now()
				testVar := lastMsgTime.Format("2006-01-02 15:04:05")
				logBuf.WriteString("After receiving done messages Receiver:" + testVar + "\n")
				msgCount++
				log.Printf("msgCount:%d", msgCount)
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))

				var tr TranscodeResult
				resultJson =  string(msg.Value)
				json.Unmarshal([]byte(resultJson),&tr)
				var dur=0.0
				if !tr.Error{
					if dur, err = strconv.ParseFloat(tr.Duration, 64); err != nil {
						fmt.Println("Error Converting dur:")
						fmt.Println(err)
					}
					fps := float64(float64(totalFramesPerVideo)/dur)
					sFps := strconv.FormatFloat(fps, 'f', 2, 64)
					resultBuf.WriteString("success," + tr.Key +  "," + tr.Starttime +  "," + tr.Endtime + "," + tr.Duration + "," + tr.Pod_id + "," + sFps+ "\n")
				}else{
					if !runError{
						runError = true
						runErrorDesc = runErrorDesc  + tr.Errordesc + ";"
					}
					resultBuf.WriteString("fail," + tr.Key +  "," + tr.Starttime +  "," + tr.Endtime + "," + tr.Duration + "," + tr.Pod_id + "," + "0" +  "\n")
				}

				recdMsgsKeySlice = append(recdMsgsKeySlice, string(tr.Key))

				if msgCount == numMessages {
					testinProgress = false
					log.Printf("testinProgress3:%v \n", testinProgress)

					msgsEndTime = time.Now()
					timeDiff := msgsEndTime.Sub(msgsStartTime).Seconds()
					//Check if the correct messages were processed. Some time if the test fails and Kafka messages are not retired or cleaned up, old messages maybe processed.
					allMsgsPresent:= true
					log.Printf("allMsgsPresent" + strconv.FormatBool(allMsgsPresent))
					sentMsgsCount := len(sentMsgsKeySlice)
					log.Printf("sentMsgsKeySlice: %d ", sentMsgsCount)
					for _, sKey := range sentMsgsKeySlice {
						log.Printf("sentMsgsKeySlice Key: %s ", sKey)
					}


					recdMsgsCount :=len(recdMsgsKeySlice)
					log.Printf("recdMsgsCount: %d ", recdMsgsCount)
					for _, sKey := range recdMsgsKeySlice {
						log.Printf("recdMsgsKeySlice Key: %s ", sKey)
					}

					if(sentMsgsCount == recdMsgsCount) {
						//check the actual keys
						for _, sKey := range sentMsgsKeySlice {
							log.Printf("Key: %s ", sKey)
							if slice.Contains(recdMsgsKeySlice, sKey)==false{
								allMsgsPresent = false
								break
							}
						}
					}
					log.Printf("allMsgsPresent" + strconv.FormatBool(allMsgsPresent))
					if (allMsgsPresent) {
						if !runError{
							podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
							totalTestResultBuf.WriteString("Status," + "success" + "\n")
							totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
							totalTestResultBuf.WriteString("Config_name," + config_name + "\n")
							totalTestResultBuf.WriteString("Category," + category + "\n")
							totalTestResultBuf.WriteString("Mode," + txMode + "\n")
							totalTestResultBuf.WriteString("InputVidFPS," + strconv.Itoa(inputVidFPS) + "\n")
							totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
							totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")

							aggregateFPS := float64(totalFramesPerVideo * nNumberOfPods) / timeDiff
							totalTestResultBuf.WriteString("Aggregate FPS," + strconv.FormatFloat(aggregateFPS,  'f', 2, 64) + "\n")
						}else{
							podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
							totalTestResultBuf.WriteString("Status," + "fail," + "Error from Pod:"+runErrorDesc  + "\n")
							totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
							totalTestResultBuf.WriteString("Config_name," + config_name + "\n")
							totalTestResultBuf.WriteString("Category," + category + "\n")
							totalTestResultBuf.WriteString("Mode," + txMode + "\n")
							totalTestResultBuf.WriteString("InputVidFPS," + strconv.Itoa(inputVidFPS) + "\n")
							totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
							totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
							totalTestResultBuf.WriteString("Aggregate FPS," + "0.0"+ "\n")
						}

					}else{
						log.Printf("Error: Old kafka messages found\n")
						podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
						totalTestResultBuf.WriteString("Status," + "fail," + "Error Running the workload. Cleanup pods and Kafka. Install Kafka again to run the test," + "\n")
						totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
						totalTestResultBuf.WriteString("Config_name," + config_name + "\n")
						totalTestResultBuf.WriteString("Category," + category + "\n")
						totalTestResultBuf.WriteString("Mode," + txMode + "\n")
						totalTestResultBuf.WriteString("InputVidFPS," + strconv.Itoa(inputVidFPS) + "\n")
						totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
						totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
						totalTestResultBuf.WriteString("Aggregate FPS," + "0.0"+ "\n")
					}

					runError=false
					runErrorDesc = ""

					quit <- true

				}
				case stopTest := <-stopTestCh:
					if(stopTest){
						log.Printf("Pods did not return in time\n")
						curTime := time.Now()
						timeDiff := curTime.Sub(msgsStartTime).Seconds()
						podsVideoProcDuration := strconv.FormatFloat(timeDiff, 'f', 2, 64)
						totalTestResultBuf.WriteString("Status," + "fail," + "One of the pods may not be responding. Cleanup pods and Kafka. Install Kafka again to run the test," + "\n")
						totalTestResultBuf.WriteString("Total_duration_secs," + podsVideoProcDuration + "\n")
						totalTestResultBuf.WriteString("Config_name," + config_name + "\n")
						totalTestResultBuf.WriteString("Category," + category + "\n")
						totalTestResultBuf.WriteString("Mode," + txMode + "\n")
						totalTestResultBuf.WriteString("InputVidFPS," + strconv.Itoa(inputVidFPS) + "\n")
						totalTestResultBuf.WriteString("Number_of_pods," + numberOfPods + "\n")
						totalTestResultBuf.WriteString("vCpus_per_Pod," + numberOfCpus + "\n")
						totalTestResultBuf.WriteString("Aggregate FPS," + "0.0"+ "\n")

						testinProgress = false
						log.Printf("testinProgress4:%v \n", testinProgress)
						quit <- true
					}
				case <-quit:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Processed: %d; errors: %d\n", msgCount, errors)

}
func startKafkaProducer(numMessages int,videoName string){
	//numMessages equals count of videos in the video library
	// Setup configuration
	log.Printf("In startKafkaProducer")
	config := sarama.NewConfig()
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 3
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	//config.Producer.Return.Successes = true
	//configuring partitions in kafka setup instead and not here
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Flush.Messages = 1
	brokers := []string{"broker:9092"}
	//producer, err := sarama.NewAsyncProducer(brokers, config)
	//not tested the flush frequency parameter yet
	//config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	//config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	//config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms


	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	log.Printf("%s\n","Before creating channel in producer")
	var enqueued, errors int
	enqueued = 0
	quit := make(chan bool)

	// Count how many message processed
	//msgCount := 0
	// Get signal for finish
	doneCh := make(chan struct{})

	//set start time to catch ffmpeg  freezes
	lastMsgTime = time.Now()
	testVar := lastMsgTime.Format("2006-01-02 15:04:05")
	logBuf.WriteString("Setting lastMsgTime:" + testVar + "\n")

	go func() {
		log.Printf("%s\n","startKafkaProducer Running")
		log.Printf("startKafkaProducer Produced message:%d", enqueued)
		log.Printf("startKafkaProducer numMessages:%d", numMessages)

		for {
			log.Printf(" enqueued=%d\n",enqueued)
			time.Sleep(100 * time.Millisecond)
			strTime := strconv.Itoa(int(time.Now().Unix()))
			strStreamName := strTime + "_" + strconv.Itoa(enqueued) + "/" + videoName
			value := fmt.Sprintf(`{"stream_name": "%s","configurationId": "%s",  "config_name": "%s", "ffmpeg_args":"%s"}`,strStreamName,configurationId,config_name,ffmpeg_args)

			log.Printf("Message queued to be sent%s\n",value)

			msg := &sarama.ProducerMessage{
				Topic: "testTopicTranscodeFinal1",
				Key:  sarama.ByteEncoder(strconv.Itoa(enqueued)),
				Value: sarama.ByteEncoder(value),
			}
			select {
			case producer.Input() <- msg:
				t := time.Now()
				testVar := t.Format("2006-01-02 15:04:05")
				logBuf.WriteString("Prod:" + testVar + "\n")

				enqueued++
				sentMsgsKeySlice = append(sentMsgsKeySlice, strStreamName)
				log.Printf("case Input: Produced message:%d", enqueued)
				log.Printf("case Input: numMessages:%d", numMessages)
				if enqueued == numMessages {
					quit <- true
					break
				}
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			case <-quit:
				log.Printf("In Quit")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}
/*func startKafkaProducer(numMessages int, videoName string){
	//numMessages equals count of videos in the video library
	// Setup configuration
	log.Printf("In startKafkaProducer")
	config := sarama.NewConfig()
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 3
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	//config.Producer.Return.Successes = true
	//configuring partitions in kafka setup instead and not here
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Flush.Messages = 1
	brokers := []string{"broker:9092"}
	//producer, err := sarama.NewAsyncProducer(brokers, config)
	//not tested the flush frequency parameter yet
	//config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	//config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	//config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms


	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	log.Printf("%s\n","Before creating channel in producer")
	var enqueued, errors int
	enqueued = 0
	errors = 0
	//need to channels to stop sending messsages
	//quit := make(chan bool)
	// Get signal for finish
	doneCh := make(chan struct{})

	go func() {
		log.Printf("%s\n","startKafkaProducer Running")
		log.Printf("startKafkaProducer Produced message:%d, numMessages:%d", enqueued, numMessages)


		producerloop:for {
			log.Printf(" enqueued=%d\n",enqueued)
			time.Sleep(100 * time.Millisecond)
			strTime := strconv.Itoa(int(time.Now().Unix()))
			strStreamName := strTime + "_" + strconv.Itoa(enqueued) + "/" + videoName
			value := fmt.Sprintf(`{"stream_name": "%s","configurationId": "%s",  "config_name": "%s", "ffmpeg_args":"%s"}`,strStreamName,configurationId,config_name,ffmpeg_args)

			log.Printf("Message queued to be sent%s\n",value)

			msg := &sarama.ProducerMessage{
				Topic: "testTopicTranscodeFinal1",
				Key:  sarama.ByteEncoder(strconv.Itoa(enqueued)),
				Value: sarama.ByteEncoder(value),
			}
			select {
			case producer.Input() <- msg:
				t := time.Now()
				testVar := t.Format("2006-01-02 15:04:05")
				logBuf.WriteString("Prod:" + testVar + "\n")

				enqueued++
				sentMsgsKeySlice = append(sentMsgsKeySlice, strStreamName)
				log.Printf("case Input: Produced message:%d, numMessages:%d\n", enqueued, numMessages)
				if enqueued == numMessages {
					log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
					wg.Done()
					break producerloop
				}
			case err := <-producer.Errors():
				errors++
				log.Printf("Failed to produce message:%s \n", err)
			default :
				log.Printf("Default:\n")
			}
		}
	}()
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}*/
func randString() string {
	n := rand.Intn(6) + 5
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ranDate() string {
	min := time.Date(2000, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2019, 7, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).String()
}

func randFloat(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	res = math.Floor(res*100) / 100
	return res
}

func randUser() string {
	prefix := "user_"
	id := rand.Intn(1000)

	return fmt.Sprintf("%s%d", prefix, id)
}

