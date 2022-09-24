package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gopkg.in/restruct.v1"
	"io/ioutil"
	"os"
	"strings"
)

type Header struct {
	Version            uint16
	AttestationKeyType uint16
	TeeType            uint32
	QeSvn              uint16
	PceSvn             uint16
	QeVendorId         [16]byte
	UserData           [20]byte
}

type SgxReportBodyT struct {
	CpuSvn          [16]uint8                                    /* (0) Security Version of the CPU */
	MiscSelect      uint32                                       /* (16) Which fields defined in SSA.MISC */
	Reserved1       [12]uint8                                    /* (20) */
	SgxIsvextProdId [16]uint8                                    /* (32) ISV assigned Extended Product ID */
	SgxAttributes   struct {
		Flags uint64
		Xfrm  uint64
	} /* (48) Any special Capabilities the Enclave possess */
	MrEnclave      [32]byte                                     /* (64) The value of the enclave's ENCLAVE measurement */
	Reserved2      [32]uint8                                    /* (96) */
	MrSigner       [32]byte                                     /* (128) The value of the enclave's SIGNER measurement */
	Reserved3      [32]uint8                                    /* (160) */
	ConfigId       [64]uint8                                    /* (192) CONFIGID 11*/
	SgxIsvProdId   uint16                                       /* (256) Product ID of the Enclave */
	SgxIsvSvn      uint16                                       /* (258) Security Version of the Enclave */
	SgxConfigSvn   uint16                                       /* (260) CONFIGSVN 11*/
	Reserved4      [42]uint8                                    /* (262) */
	SgxIsvFamilyId [16]uint8                                    /* (304) ISV assigned Family ID */
	SgxReportData  [64]byte                                     /* (320) Data provided by the user */
}

// Ecdsa Structure sequence - 1
type SgxQuote struct {
	Header
	SgxReportBodyT
	SignatureLen uint32 /* 432 */
	Signature    []byte
}

type AttestationTokenClaim struct {
	MrEnclave    string          `json:"mrenclave,omitempty"`
	MrSigner     string          `json:"mrsigner,omitempty"`
	ReportData   string          `json:"reportdata,omitempty"`
	IsvProductId *uint16         `json:"isvprodid,omitempty"`
	IsvSvn       *uint16         `json:"isvsvn,omitempty"`
}

const (
	BINARY_FILE = ".dat"
	BASE64_FILE = ".txt"
	CONVERT_OP = "convert"
	BASE64_FMT = "base64"
)

func getContent() ([]byte, error) {
	var quote []byte
	var err error
	var fileNameNew string
	var convert = false
	fileType := os.Args[1]
	fileName := os.Args[2]

	if len(os.Args) > 3 && strings.ToLower(os.Args[3]) == CONVERT_OP {
		convert = true
	}

	if strings.HasSuffix(strings.ToLower(fileName), BINARY_FILE) {
		fileNameNew = fileName + BASE64_FILE
	} else {
		fileNameNew = fileName + BINARY_FILE
	}
	content, err1 := ioutil.ReadFile(fileName)
	if err1 != nil {
		fmt.Println(fileName, " not found")
		os.Exit(1)
	}
	if strings.HasPrefix(strings.ToLower(fileType), BASE64_FMT) {
		quote, err = base64.StdEncoding.DecodeString(string(content))
		if err != nil {
			fmt.Println("Input file is not base64 format")
			os.Exit(1)
		}
		if convert {
			ioutil.WriteFile(fileNameNew, quote, 0664)
		}
		return quote, err
	}
	if convert {
		tmp := base64.StdEncoding.EncodeToString(content)
		ioutil.WriteFile(fileNameNew, []byte(tmp), 0664)
	}
	return content, err1
}

func main() {
	var parsedSGXQuote SgxQuote

	content, err := getContent()
	if err != nil {
		fmt.Println("Failed to get quote contents from file ", os.Args[2])
	}
	err = restruct.Unpack(content, binary.LittleEndian, &parsedSGXQuote)
	if err != nil {
		fmt.Println(err, "Failed to parse SGX quote")
	}

	ac := AttestationTokenClaim{
		MrEnclave:    fmt.Sprintf("%02x", parsedSGXQuote.MrEnclave),
		MrSigner:     fmt.Sprintf("%02x", parsedSGXQuote.MrSigner),
		ReportData:   fmt.Sprintf("%02x", getSHA256Hash(parsedSGXQuote.SgxReportData)),
		IsvProductId: &parsedSGXQuote.SgxIsvProdId,
		IsvSvn:       &parsedSGXQuote.SgxIsvSvn,
	}
	result, _ := json.Marshal(ac)
	fmt.Println(string(result))

}

func getSHA256Hash(reportData [64]byte) []byte {
	hashValue := make([]byte, sha256.Size)
	for i := 0; i < sha256.Size; i++ {
		hashValue[i] = reportData[i]
	}
	return hashValue
}
