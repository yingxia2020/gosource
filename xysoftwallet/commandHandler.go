package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WriteDataRequest struct {
	Token string `json:"token"`
}

type DevicePortDiscoveryData struct {
	SsId string `json:"ssId"`
	EncryptedData string `json:"encryptedData"`
}

type BrowserCommandRequest struct {
	Cmd string `json:"cmd"`
	Token string `json:"token"`
}

type BrowserCommandResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	RegStatus string `json:"regStatus"`
	Data string `json:"data"`
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	// Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Check if it's a preflight request
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only POST requests are allowed")
		return
	}

	// Parse the request body
	var commandRequest BrowserCommandRequest
	err := json.NewDecoder(r.Body).Decode(&commandRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}

	// Validate the required fields
	if commandRequest.Cmd == "" || commandRequest.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cmd and token fields are required")
		return
	}

	// Process the command and token
	// You can perform any desired logic here
	fmt.Printf("Received command: %s\n", commandRequest.Cmd)
	//fmt.Printf("Received token: %s\n", commandRequest.Token)
	fmt.Println("---------------------------------------------------------")

	var response []byte
	if commandRequest.Cmd == "REGISTER_DEVICE" {
		writeDataRequest := WriteDataRequest{
			Token: commandRequest.Token,
		}
		data, _ := json.Marshal(writeDataRequest)
		request := SocRequest{
			PublicID: PublicId,
			Api: "REGISTER_SOFTWALLET_DEVICE",
			Data: string(data),
		}
		result, err := sendRegisterDeviceRequest(request)
		if err != nil || result != true {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to send register device to cloud: %v", err)
			return
		}
		// update device state to registered
		propertiesMap[STATE] = string(REGISTERED)
		// update capacity value
		propertiesMap[CAPACITY] = "2"
		// write sha256 hash of username to device
		propertiesMap[HASHED_USERNAME] = getHashedUsernameFromJwtToken(commandRequest.Token)
		writeProperties(DEVICE_FILE, propertiesMap)

		response = composeBrowserResponse("OK", true)
	} else if commandRequest.Cmd == "GET_NEXT_JOB" {
		response = composeBrowserResponse("", true)
	} else if commandRequest.Cmd == "QUERY_STATUS" {
		if propertiesMap[STATE] != string(REGISTERED) {
			response = composeBrowserResponseReg("NOT_REGISTERED", true)
		} else {
			response = composeBrowserResponseReg("REGISTERED", true)
		}
	} else if commandRequest.Cmd == "VERIFY_DEVICE_OWNER" {
		if propertiesMap[STATE] != string(REGISTERED) {
			response = composeBrowserResponse("unauthorized", false)
		} else {
			response = composeBrowserResponse("SUCCESS", true)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Received not able to handle command: %v", commandRequest.Cmd)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func handlePublicCommand(w http.ResponseWriter, r *http.Request) {		// Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Check if it's a preflight request
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only POST requests are allowed")
		return
	}

	// Parse the request body
	var commandRequest BrowserCommandRequest
	err := json.NewDecoder(r.Body).Decode(&commandRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}

	// Validate the required fields
	if commandRequest.Cmd == ""  {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cmd field is required")
		return
	}

	// Process the command
	// You can perform any desired logic here
	fmt.Printf("Received public command: %s\n", commandRequest.Cmd)
	fmt.Println("---------------------------------------------------------")

	var response []byte
	if commandRequest.Cmd == "PORT_DISCOVERY" {
		portDiscoverData := DevicePortDiscoveryData {
			SsId: "640279da-eac6-4ab6-a0c4-36af7fed2400",
			EncryptedData: "Gfz9/cGfGnF6PWviwnlH6W90l9t2sUTmk/EMzVhoyLFjMugZkEuMSjmH2PWVfTfhfPJeFL2+4u3y9zR+HhdTdsuoOXhkQA1sr4O5RUqj2EY=\\",
		}
		discoverData, _ := json.Marshal(portDiscoverData)
		response = composeBrowserResponseData(string(discoverData), true)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Received not able to handle command: %v", commandRequest.Cmd)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}


func composeBrowserResponse(message string, status bool) []byte {
	var response BrowserCommandResponse
	if len(message) == 0 {
		response = BrowserCommandResponse{
			Status: status,
		}
	} else {
		response = BrowserCommandResponse{
			Message: message,
			Status:  status,
		}
	}
	jsonData, _ := json.Marshal(response)
	return jsonData
}


func composeBrowserResponseReg(regStatus string, status bool) []byte {
	response := BrowserCommandResponse{
		Status:  status,
		RegStatus: regStatus,
	}
	jsonData, _ := json.Marshal(response)
	return jsonData
}

func composeBrowserResponseData(data string, status bool) []byte {
	response := BrowserCommandResponse{
		Status:  status,
		Data: data,
	}
	jsonData, _ := json.Marshal(response)
	return jsonData
}