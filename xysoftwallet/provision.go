package main

import "fmt"

func initProvisioning(propertiesMap map[string]string) bool {
	fmt.Println("Create CSR for LAK(Soc)")
	privkeySoc, pubkeySoc := generate384KeyPair()
	if len(privkeySoc) == 0 {
		return false
	}
	fmt.Println("Create CSR for LDEVID(Tee)")
	privkeyTee, pubkeyTee := generate384KeyPair()
	if len(privkeyTee) == 0 {
		return false
	}

	fmt.Println("Send both to issue certificates")
	request := SocRequest{
		PublicID: PublicId,
		Api: "INIT_SOFTWALLET_PROVISIONING",
		CsrSoc: pubkeySoc,
		CsrTee: pubkeyTee,
		CertMfg: CertMfc,
		Data: SkMfc,  // need send private key to generate csr
	}
	certSoc, certTee, err := sendInitProvisioningRequest(request)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println("Cert for LAK(Soc): ", certSoc)
	fmt.Println("Cert for LDEVID(Tee): ", certTee)
	propertiesMap["skLak"] = privkeySoc
	propertiesMap["certLak"] = certSoc
	propertiesMap["skLDevId"] = privkeyTee
	propertiesMap["certLDevId"] = certTee
	writeProperties(DEVICE_FILE, propertiesMap)
	return true
}

func verifyManufacturingCertificate(certManu string) bool {
	publicId := PublicId
	if len(publicId) == 0 {
		return false
	} else {
		fmt.Println(publicId)
	}

	if len(certManu) == 0 {
		return false
	}

	// check certManu with cloud
	fmt.Println(certManu)
	request := SocRequest{
		PublicID: publicId,
		Api: "VERIFY_MANUFACTURE_CERTIFICATE",
		Data: certManu,
	}
	result, err := sendVerifyRotRequest(request)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return result
}