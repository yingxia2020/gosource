package main

import (
	"fmt"
	"log"
	"net/http"
)

const DEVICE_FILE = "device.properties"
const SK_MANUF = "skManuf"
const CERT_MANUF = "certManuf"
const STATE = "state"
const PUBLIC_ID = "publicId"
const CAPACITY = "capacity"
const HASHED_USERNAME = "hashedUsername"

var propertiesMap map[string]string
var SkMfc = ""
var CertMfc = ""
var DeviceState = NONE
var PublicId = ""

type State string
const (
	NONE            State = "NONE"
	ROT_PROVISIONED State = "ROT_PROVISIONED"
	UNREGISTERED    State = "UNREGISTERED"
	REGISTERED      State = "REGISTERED"
)


func init() {
	// Perform initialization tasks here
	fmt.Println("Soft-wallet initializing the application...")
	fmt.Println("Read manufacturing keys ...")

	// Read properties from the file
	var err error
	propertiesMap, err = readProperties(DEVICE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var exist = true
	SkMfc, exist = propertiesMap[SK_MANUF]
	if !exist {
		log.Fatal("Manufacturing private key not found")
	}
	CertMfc, exist = propertiesMap[CERT_MANUF]
	if !exist {
		log.Fatal("Manufacturing certificate not found")
	}
	tmp, exist := propertiesMap[STATE]
	if !exist {
		log.Fatal("Device state not found")
	} else {
		DeviceState = State(tmp)
	}

	PublicId, exist = propertiesMap[PUBLIC_ID]
	if !exist {
		PublicId = generatePublidId()
		propertiesMap[PUBLIC_ID] = PublicId
		writeProperties(DEVICE_FILE, propertiesMap)
	} else {
		fmt.Println("publicId: ", PublicId)
	}

	switch DeviceState {
	case NONE:
		fmt.Println("Verify manufacturing certificate")
		ok := verifyManufacturingCertificate(CertMfc)
		if !ok {
			log.Fatal("Verify manufacturing certificate failed")
		}
		propertiesMap[STATE] = string(ROT_PROVISIONED)
		writeProperties(DEVICE_FILE, propertiesMap)
		fmt.Println("Init provisioning")
		ok = initProvisioning(propertiesMap)
		if !ok {
			log.Fatal("Init provisioning failed")
		}
		propertiesMap[STATE] = string(UNREGISTERED)
		propertiesMap[CAPACITY] = "0"
		writeProperties(DEVICE_FILE, propertiesMap)
		fmt.Println("Heart beat is omitted now")
		break
	case ROT_PROVISIONED:
		fmt.Println("Init provisioning")
		ok := initProvisioning(propertiesMap)
		if !ok {
			log.Fatal("Init provisioning failed")
		}
		propertiesMap[STATE] = string(UNREGISTERED)
		propertiesMap[CAPACITY] = "0"
		writeProperties(DEVICE_FILE, propertiesMap)
		fmt.Println("Heart beat is omitted now")
		break
	case UNREGISTERED:
		fmt.Println("Provision is done, please go for online registration")
		break
	default:
		// device is provisioned and registered, do nothing
		fmt.Println("Welcome to use Intel Soft Wallet")
		break
	}
}


func main() {
	// Your application logic goes here
	fmt.Println("Soft-wallet service started.")

	http.HandleFunc("/v1/command", handleCommand)
	http.HandleFunc("/v1/public_command", handlePublicCommand)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if it's a preflight request
		if r.Method == "OPTIONS" {
			return
		}
		fmt.Fprint(w, "Hello world!")
	})

	// Read the PKCS12 (.pfx) file
	//pfxData, err := ioutil.ReadFile("keystore.pfx")
	//if err != nil {
	//	fmt.Println("Failed to read PKCS12 file:", err)
	//	return
	//}
	//
	//// Decrypt the PKCS12 data with the provided password
	//password := "intel123"
	//privateKey, cert, err := pkcs12.Decode(pfxData, password)
	//if err != nil {
	//	fmt.Println("Failed to decode PKCS12 data:", err)
	//	return
	//}
	//
	//// Encode and save the private key as PEM
	//keyBytes := x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))
	//keyPEM := pem.EncodeToMemory(&pem.Block{
	//	Type:  "RSA PRIVATE KEY",
	//	Bytes: keyBytes,
	//})
	//err = ioutil.WriteFile("key.pem", keyPEM, 0644)
	//if err != nil {
	//	fmt.Println("Failed to save private key PEM file:", err)
	//	return
	//}
	//
	//// Encode and save the certificate as PEM
	//certPEM := pem.EncodeToMemory(&pem.Block{
	//	Type:  "CERTIFICATE",
	//	Bytes: cert.Raw,
	//})
	//
	//err = ioutil.WriteFile("cert.pem", certPEM, 0644)
	//if err != nil {
	//	fmt.Println("Failed to save certificate PEM file:", err)
	//	return
	//}

	// Load the private key and certificate as a tls.Certificate
	// certificate, err := tls.X509KeyPair(certPEM, keyPEM)
	//certificate, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Create the TLS configuration
	//config := &tls.Config{
	//	Certificates: []tls.Certificate{certificate},
	//	InsecureSkipVerify: true,
	//}
	//
	//// Create the HTTPS server
	//server := &http.Server{
	//	Addr:      ":20600",
	//	TLSConfig: config,
	//}
	//
	//log.Fatal(server.ListenAndServeTLS("", ""))

	log.Fatal(http.ListenAndServe(":20600", nil))
}
