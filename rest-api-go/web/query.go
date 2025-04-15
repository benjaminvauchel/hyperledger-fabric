// package web

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// // QueryParams represents the query parameters for credential lookups
// type QueryParams struct {
// 	Function   string   `json:"function"`
// 	Args       []string `json:"args"`
// 	ChaincodeID string   `json:"chaincodeid"`
// 	ChannelID   string   `json:"channelid"`
// }

// // QueryCredentialsHandler handles general credential query requests
// func (setup *OrgSetup) QueryCredentialsHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Received Query Credentials request")
	
// 	// Parse the request body to get the parameters
// 	var params QueryParams
// 	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
// 		HandleError(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
	
// 	// Validate chaincodeid and channelid
// 	if params.ChaincodeID == "" || params.ChannelID == "" {
// 		HandleError(w, "Missing required parameters: chaincodeid and channelid", http.StatusBadRequest)
// 		return
// 	}
	
// 	// Other optional filter parameters
// 	talentID := r.URL.Query().Get("talentId")
// 	institution := r.URL.Query().Get("institution")
// 	company := r.URL.Query().Get("company")
	
// 	// Determine which query function to use based on parameters
// 	var function string
// 	var args []string
	
// 	if talentID != "" {
// 		function = "GetCredentialsByTalent"
// 		args = []string{talentID}
// 	} else if institution != "" {
// 		function = "GetCredentialsByInstitution"
// 		args = []string{institution}
// 	} else if company != "" {
// 		function = "GetCredentialsByCompany"
// 		args = []string{company}
// 	} else {
// 		function = "GetAllCredentials"
// 		args = []string{}
// 	}
	
// 	// Log the query details for debugging
// 	log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", params.ChannelID, params.ChaincodeID, function, args)
	
// 	// Execute the query
// 	result, err := executeQuery(setup, params.ChannelID, params.ChaincodeID, function, args)
// 	if err != nil {
// 		HandleError(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	
// 	// Process the result (could be a JSON string that needs to be parsed)
// 	var responseData interface{}
// 	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
// 		// If it's not valid JSON, use the raw string
// 		responseData = result
// 	}
	
// 	// Send success response
// 	HandleSuccess(w, "Query executed successfully", responseData)
// }

// // CustomQueryHandler handles custom query requests with specific function and args
// func (setup *OrgSetup) CustomQueryHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Received Custom Query request")
	
// 	// Parse the request body to get the parameters
// 	var params QueryParams
// 	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
// 		HandleError(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
	
// 	// Validate chaincodeid, channelid, and function
// 	if params.ChaincodeID == "" || params.ChannelID == "" || params.Function == "" {
// 		HandleError(w, "Missing required parameters: chaincodeid, channelid, and function", http.StatusBadRequest)
// 		return
// 	}
	
// 	// Log the query details for debugging
// 	log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", params.ChannelID, params.ChaincodeID, params.Function, params.Args)
	
// 	// Execute the query
// 	result, err := executeQuery(setup, params.ChannelID, params.ChaincodeID, params.Function, params.Args)
// 	if err != nil {
// 		HandleError(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	
// 	// Process the result (could be a JSON string that needs to be parsed)
// 	var responseData interface{}
// 	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
// 		// If it's not valid JSON, use the raw string
// 		responseData = result
// 	}
	
// 	// Send success response
// 	HandleSuccess(w, "Query executed successfully", responseData)
// }

// // executeQuery handles the common query execution logic
// func executeQuery(setup *OrgSetup, channelID, chainCodeID, function string, args []string) (string, error) {
// 	// Get the network and contract
// 	network := setup.Gateway.GetNetwork(channelID)
// 	contract := network.GetContract(chainCodeID)
	
// 	// Evaluate transaction (query)
// 	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
// 	if err != nil {
// 		return "", err
// 	}
	
// 	return string(evaluateResponse), nil
// }

package web
import "github.com/gorilla/mux"

import (
	"encoding/json"
	"log"
	"net/http"
)

// QueryParams represents the query parameters for credential lookups
type QueryParams struct {
	Function   string   `json:"function"`
	Args       []string `json:"args"`
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
}

// QueryCredentialsHandler handles general credential query requests
func (setup *OrgSetup) QueryCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Query Credentials request")

	// Extract query parameters
	query := r.URL.Query()
	chainCodeID := query.Get("chaincodeid")
	channelID := query.Get("channelid")
	credentialID := query.Get("credentialid")
	credentialType := query.Get("credentialtype")
	talentID := query.Get("talentid")
	institution := query.Get("institution")
	company := query.Get("company")

	// Validate chaincodeid and channelid
	if chainCodeID == "" || channelID == "" {
		HandleError(w, "Missing required parameters: chaincodeid and channelid", http.StatusBadRequest)
		return
	}

	// Decide which function to call
	var function string
	var args []string

	switch {
	case credentialID != "":
		switch credentialType {
		case "academic":
			function = "GetAcademicCredential"
		case "professional":
			function = "GetProfessionalCredential"
		default:
			function = "GetBaseCredential"
		}
		args = []string{credentialID}

	case talentID != "":
		function = "GetCredentialsByTalent"
		args = []string{talentID}

	case institution != "":
		function = "GetCredentialsByInstitution"
		args = []string{institution}

	case company != "":
		function = "GetCredentialsByCompany"
		args = []string{company}

	default:
		function = "GetAllCredentials"
		args = []string{}
	}

	log.Printf("Executing query: channel=%s, chaincode=%s, function=%s, args=%v",
		channelID, chainCodeID, function, args)

	// Execute query
	result, err := executeQuery(setup, channelID, chainCodeID, function, args)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		HandleError(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Try to decode JSON result
	var responseData interface{}
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		responseData = result // fallback to raw string
	}
	log.Printf("Raw query result: %s\n", result)

	HandleSuccess(w, "Query executed successfully", responseData)
}


// CustomQueryHandler handles custom query requests with specific function and args
func (setup *OrgSetup) CustomQueryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Custom Query request")
	
	// Retrieve query parameters directly from the URL query string
	chainCodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")
	function := r.URL.Query().Get("function")
	args := r.URL.Query()["args"]

	// Validate required parameters
	if chainCodeID == "" || channelID == "" || function == "" {
		HandleError(w, "Missing required parameters: chaincodeid, channelid, and function", http.StatusBadRequest)
		return
	}
	
	// Log the query details for debugging
	log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", channelID, chainCodeID, function, args)
	
	// Execute the query
	result, err := executeQuery(setup, channelID, chainCodeID, function, args)
	if err != nil {
		HandleError(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Process the result (could be a JSON string that needs to be parsed)
	var responseData interface{}
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		// If it's not valid JSON, use the raw string
		responseData = result
	}
	
	// Send success response
	HandleSuccess(w, "Query executed successfully", responseData)
}

// executeQuery handles the common query execution logic
func executeQuery(setup *OrgSetup, channelID, chainCodeID, function string, args []string) (string, error) {
	// Get the network and contract
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeID)
	
	// Evaluate transaction (query)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		return "", err
	}
	
	return string(evaluateResponse), nil
}

// maybe useless
func (setup *OrgSetup) GetCredentialByTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	credentialID := vars["id"]
	credType := r.URL.Query().Get("type") // "academic", "professional", or "base"
	chaincodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")

	if credentialID == "" || chaincodeID == "" || channelID == "" {
		HandleError(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	var function string
	switch credType {
	case "academic":
		function = "GetAcademicCredential"
	case "professional":
		function = "GetProfessionalCredential"
	case "base":
		function = "GetBaseCredential"
	default:
		function = "GetTalentCredential"
	}

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chaincodeID)

	evaluateResult, err := contract.EvaluateTransaction(function, credentialID)
	if err != nil {
		HandleError(w, "Failed to evaluate transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(evaluateResult)
}

// maybe useless
func (setup *OrgSetup) GetAllCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	chaincodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")

	if chaincodeID == "" || channelID == "" {
		HandleError(w, "Missing chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chaincodeID)

	result, err := contract.EvaluateTransaction("GetAllCredentials")
	if err != nil {
		HandleError(w, "Failed to evaluate transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
