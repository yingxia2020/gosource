package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NOTE: SOC <-> LDEVID   TEE <-> LAK
type SocResponse struct {
	// Define the structure of the response JSON
	PublicID string `json:"publicId"`
	Api string `json:"api"`
	CertSoc string `json:"certSoc"`
	CertTee string `json:"certTee"`
	ErrorCode string `json:"errorCode"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SocRequest struct {
	// Define the structure of the request JSON
	PublicID string `json:"publicId"`
	Api string `json:"api"`
	Data string `json:"data"`
	CsrSoc string `json:"csrSoc"`
	CsrTee string `json:"csrTee"`
	CertMfg string `json:"certMfg"`
}

const URL = "http://localhost:9007/v1/msg"

//func main() {
//	requestBody := SocRequest{
//		Data: "Hello, Server!",
//	}
//
//	resp, err := sendHttpRequest(requestBody)
//
//}

func sendInitProvisioningRequest(requestBody SocRequest) (string, string, error) {
	response, err := sendHttpRequest(requestBody)

	if err != nil {
		return "", "", err
	} else {
		if response.ErrorCode != "OK" {
			return "", "", fmt.Errorf("error in http response")
		}
		return response.CertSoc, response.CertTee, nil
	}
}

func sendVerifyRotRequest(requestBody SocRequest) (bool, error) {
	response, err := sendHttpRequest(requestBody)

	if err != nil {
		return false, err
	}

	if response.ErrorCode != "OK" {
		return false, nil
	}

	return true, nil
}

func sendRegisterDeviceRequest(requestBody SocRequest) (bool, error) {
	response, err := sendHttpRequest(requestBody)

	if err != nil {
		return false, err
	}

	if response.ErrorCode != "OK" {
		return false, nil
	}

	return true, nil
}

func sendHttpRequest(requestBody SocRequest) (SocResponse, error) {
	var response SocResponse
	// Convert the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return response, err
	}

	method := "POST"
	payload := strings.NewReader(string(jsonData))

	req, err := http.NewRequest(method, URL, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return response, err
	}

	req.Header.Set("Content-Type", "application/json") // Set the desired Content-Type header value
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	// Send the HTTP POST request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return response, err
	}
	defer resp.Body.Close()

	// Deserialize the response into the struct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return response, err
	}

	return response, nil
}

