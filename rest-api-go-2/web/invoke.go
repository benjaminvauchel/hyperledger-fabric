package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Credential models the credential data for requests
type Credential struct {
	CredentialID   string `json:"credentialId"`
	TalentID       string `json:"talentId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Skills         string `json:"skills"`
	Education      string `json:"education,omitempty"`
	WorkExperience string `json:"workExperience,omitempty"`
	Institution    string `json:"institution,omitempty"`
	Company        string `json:"company,omitempty"`
	VerifiedBy	   string `json:"verifiedBy,omitempty"`
}

// CredentialRequest models the data for requests
type CredentialRequest struct {
    ChainCodeID string    `json:"chaincodeid"`
    ChannelID   string    `json:"channelid"`
    Credential  Credential `json:"credential"`
}

// // VerificationRequest models the verification data for requests
// type VerificationRequest struct {
// 	VerifiedBy string `json:"verifiedBy"`
// }

// TransactionResult holds the result of a blockchain transaction
type TransactionResult struct {
	TxID     string `json:"transactionId"`
	Response string `json:"response"`
}

// ValidateCredential validates credential fields
func ValidateCredential(cred *Credential, credType string) error {
	if cred.CredentialID == "" {
		return errors.New("credentialId is required")
	}
	if cred.TalentID == "" {
		return errors.New("talentId is required")
	}
	if cred.FirstName == "" {
		return errors.New("firstName is required")
	}
	if cred.LastName == "" {
		return errors.New("lastName is required")
	}
	
	// Validate specific credential type fields
	if credType == "academic" {
		if cred.Education == "" {
			return errors.New("education is required for academic credentials")
		}
		if cred.Institution == "" {
			return errors.New("institution is required for academic credentials")
		}
	} else if credType == "professional" {
		if cred.WorkExperience == "" {
			return errors.New("workExperience is required for professional credentials")
		}
		if cred.Company == "" {
			return errors.New("company is required for professional credentials")
		}
	}
	
	return nil
}

// CreateAcademicCredentialHandler handles requests to create academic credentials
func (setup *OrgSetup) CreateAcademicCredentialHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received Create Academic Credential request")
    
    // Parse the request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        HandleError(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Parse the request body into CredentialRequest struct
    var request CredentialRequest
    if err := json.Unmarshal(body, &request); err != nil {
        HandleError(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // Validate the credential data
    if err := ValidateCredential(&request.Credential, "academic"); err != nil {
        HandleError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get the network and contract from the gateway
    network := setup.Gateway.GetNetwork(request.ChannelID)
    contract := network.GetContract(request.ChainCodeID)

    // Define function name and arguments
    function := "CreateAcademicCredential"
    args := []string{
        request.Credential.CredentialID,
        request.Credential.TalentID,
        request.Credential.FirstName,
        request.Credential.LastName,
        request.Credential.Skills,
        request.Credential.Education,
        request.Credential.Institution,
    }

    // Log the transaction details for debugging
    log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", request.ChannelID, request.ChainCodeID, function, args)

    // Execute the transaction
    result, err := executeTransaction(contract, function, args)
    if err != nil {
        HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Send success response
    HandleSuccess(w, "Academic credential created successfully", result)
}


// CreateProfessionalCredentialHandler handles requests to create professional credentials
func (setup *OrgSetup) CreateProfessionalCredentialHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received Create Professional Credential request")
    
    // Parse the request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        HandleError(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Parse the request data (includes chaincodeid, channelid, and credential)
    var request CredentialRequest
    if err := json.Unmarshal(body, &request); err != nil {
        HandleError(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate the credential data
    if err := ValidateCredential(&request.Credential, "professional"); err != nil {
        HandleError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get the network and contract from the gateway
    network := setup.Gateway.GetNetwork(request.ChannelID)
    contract := network.GetContract(request.ChainCodeID)

    // Define function name and arguments
    function := "CreateProfessionalCredential"
    args := []string{
        request.Credential.CredentialID,
        request.Credential.TalentID,
        request.Credential.FirstName,
        request.Credential.LastName,
        request.Credential.Skills,
        request.Credential.WorkExperience,
        request.Credential.Company,
    }

    // Log the transaction details for debugging
    log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", request.ChannelID, request.ChainCodeID, function, args)

    // Execute the transaction
    result, err := executeTransaction(contract, function, args)
    if err != nil {
        HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Send success response
    HandleSuccess(w, "Professional credential created successfully", result)
}

func (setup *OrgSetup) ApproveCredentialHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Approve Credential request")

	// Check if the client is from Org1MSP
	if setup.MSPID != "Org1MSP" {
		HandleError(w, "Permission denied: only users from Org1MSP can approve credentials", http.StatusForbidden)
		return
	}

	// Get the credential ID from the URL
	vars := mux.Vars(r)
	credentialID := vars["id"]
	if credentialID == "" {
		HandleError(w, "Credential ID is required", http.StatusBadRequest)
		return
	}

	// Get chaincodeid and channelid from query parameters
	chaincodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")
	if chaincodeID == "" || channelID == "" {
		HandleError(w, "Missing chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	// For simplicity, we hardcode "VerifiedBy" as Org1's name
	verifiedBy := "Org1"

	// Get the network and contract from the gateway
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chaincodeID)

	// Define function name and arguments
	function := "UpdateVerificationStatus"
	args := []string{
		credentialID,   // credentialID
		"Verified",     // verification status
		verifiedBy,     // verifiedBy
	}

	log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", channelID, chaincodeID, function, args)

	// Execute the transaction
	result, err := executeTransaction(contract, function, args)
	if err != nil {
		HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	HandleSuccess(w, "Credential approved successfully", result)
}

func (setup *OrgSetup) RevokeCredentialHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Revoke Credential request")

	// Ensure only Org1 can revoke credentials
	if setup.MSPID != "Org1MSP" {
		HandleError(w, "Permission denied: only users from Org1MSP can revoke credentials", http.StatusForbidden)
		return
	}

	// Get the credential ID from the URL
	vars := mux.Vars(r)
	credentialID := vars["id"]
	if credentialID == "" {
		HandleError(w, "Credential ID is required", http.StatusBadRequest)
		return
	}

	// Get chaincodeid and channelid from query parameters
	chaincodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")
	if chaincodeID == "" || channelID == "" {
		HandleError(w, "Missing chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	// Hardcoded verifier for now
	verifiedBy := "Org1"

	// Get the network and contract
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chaincodeID)

	// Prepare arguments
	function := "UpdateVerificationStatus"
	args := []string{
		credentialID,
		"Revoked",
		verifiedBy,
	}

	log.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", channelID, chaincodeID, function, args)

	// Execute transaction
	result, err := executeTransaction(contract, function, args)
	if err != nil {
		HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Success
	HandleSuccess(w, "Credential revoked successfully", result)
}

func (setup *OrgSetup) DeleteCredentialHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Delete Credential request")

	// // Only Org1 can delete credentials (nope)
	// if setup.MSPID != "Org1MSP" {
	//	HandleError(w, "Permission denied: only Org1MSP can delete credentials", http.StatusForbidden)
	//	return
	//}

	vars := mux.Vars(r)
	credentialID := vars["id"]
	if credentialID == "" {
		HandleError(w, "Credential ID is required", http.StatusBadRequest)
		return
	}

	chaincodeID := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")
	if chaincodeID == "" || channelID == "" {
		HandleError(w, "Missing chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chaincodeID)

	result, err := executeTransaction(contract, "DeleteTalentCredential", []string{credentialID})
	if err != nil {
		HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	HandleSuccess(w, "Credential deleted successfully", result)
}

type UpdateSkillsRequest struct {
	NewSkills   string `json:"newSkills"`
	ChainCodeID string `json:"chaincodeid"`
	ChannelID   string `json:"channelid"`
}

func (setup *OrgSetup) UpdateSkillsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Update Skills request")

	vars := mux.Vars(r)
	credentialID := vars["id"]

	var req UpdateSkillsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	network := setup.Gateway.GetNetwork(req.ChannelID)
	contract := network.GetContract(req.ChainCodeID)

	args := []string{credentialID, req.NewSkills}

	result, err := executeTransaction(contract, "UpdateSkills", args)
	if err != nil {
		HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	HandleSuccess(w, "Skills updated successfully", result)
}

type UpdateNameRequest struct {
	NewFirstName string `json:"newFirstName"`
	NewLastName  string `json:"newLastName"`
	ChainCodeID  string `json:"chaincodeid"`
	ChannelID    string `json:"channelid"`
}

func (setup *OrgSetup) UpdateNameHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Update Name request")

	vars := mux.Vars(r)
	credentialID := vars["id"]

	var req UpdateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	network := setup.Gateway.GetNetwork(req.ChannelID)
	contract := network.GetContract(req.ChainCodeID)

	args := []string{credentialID, req.NewFirstName, req.NewLastName}

	result, err := executeTransaction(contract, "UpdateName", args)
	if err != nil {
		HandleError(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	HandleSuccess(w, "Name updated successfully", result)
}

// executeTransaction handles the common transaction execution logic
func executeTransaction(contract *client.Contract, function string, args []string) (*TransactionResult, error) {
	// Create the transaction proposal
	txnProposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		return nil, fmt.Errorf("error creating txn proposal: %s", err)
	}
	
	// Endorse the transaction
	txnEndorsed, err := txnProposal.Endorse()
	if err != nil {
		return nil, fmt.Errorf("error endorsing txn: %s", err)
	}
	
	// Submit the transaction
	txnCommitted, err := txnEndorsed.Submit()
	if err != nil {
		return nil, fmt.Errorf("error submitting transaction: %s", err)
	}
	
	// Return the result
	return &TransactionResult{
		TxID:     txnCommitted.TransactionID(),
		Response: string(txnEndorsed.Result()),
	}, nil
}
