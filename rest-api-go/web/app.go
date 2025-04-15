package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// APIResponse standardizes the API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CORSMiddleware adds CORS headers to enable cross-origin requests
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // In production, specify your frontend domain instead of *
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Serve starts http web server with proper routing
func Serve(setup OrgSetup) {
	router := mux.NewRouter()

	// Credentials routes
	credentials := router.PathPrefix("/credentials").Subrouter()
	
	// Create academic credential (POST)
	credentials.HandleFunc("/academic", setup.CreateAcademicCredentialHandler).Methods("POST")
	
	// Create professional credential (POST)
	credentials.HandleFunc("/professional", setup.CreateProfessionalCredentialHandler).Methods("POST")
	
	// Approve credential (PUT)
	credentials.HandleFunc("/{id}/approve", setup.ApproveCredentialHandler).Methods("PUT")
	
	// Revoke credential (PUT)
	credentials.HandleFunc("/{id}/revoke", setup.RevokeCredentialHandler).Methods("PUT")
	
	// Delete credential (DELETE)
	credentials.HandleFunc("/{id}", setup.DeleteCredentialHandler).Methods("DELETE")

	// Update skills (PUT)
	credentials.HandleFunc("/{id}/skills", setup.UpdateSkillsHandler).Methods("PUT")

	// Update name (PUT)
	credentials.HandleFunc("/{id}/name", setup.UpdateNameHandler).Methods("PUT")

	// Get credential by type (GET) - supports ?type=academic|professional|base
	credentials.HandleFunc("/{id}", setup.GetCredentialByTypeHandler).Methods("GET")

	// Get all credentials (GET)
	credentials.HandleFunc("/all", setup.GetAllCredentialsHandler).Methods("GET")

	// Query credentials (GET)
	credentials.HandleFunc("", setup.QueryCredentialsHandler).Methods("GET")
	
	// Custom query with function and args (GET)
	credentials.HandleFunc("/query", setup.CustomQueryHandler).Methods("GET")
	
	// Apply CORS middleware to all routes
	corsRouter := CORSMiddleware(router)
	
	// Set up the server
	http.Handle("/", corsRouter)
	
	// Start the server
	fmt.Println("Listening (http://localhost:3000/)...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// HandleError sends standardized error responses
func HandleError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	resp := APIResponse{
		Success: false,
		Error:   errMsg,
	}
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleSuccess sends standardized success responses
func HandleSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}