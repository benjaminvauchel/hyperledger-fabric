package main

import (
	"fmt"
	"log"
	"rest-api-go/web"
)

func main() {
	// Initialize setup for Org1
	cryptoPath := "../../talent-credentials-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}

	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		log.Fatalf("Error initializing setup for Org1: %s", err)
	}

	fmt.Println("Server running at http://localhost:3000/")
	web.Serve(*orgSetup)
}

// func main() {
// 	// Initialize setup for Org2
// 	cryptoPath := "../../talent-credentials-network/organizations/peerOrganizations/org2.example.com"
// 	orgConfig := web.OrgSetup{
// 		OrgName:      "Org2",
// 		MSPID:        "Org2MSP",
// 		CertPath:     cryptoPath + "/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem",
// 		KeyPath:      cryptoPath + "/users/User1@org2.example.com/msp/keystore/",
// 		TLSCertPath:  cryptoPath + "/peers/peer0.org2.example.com/msp/tlscacerts/tlsca.org2.example.com-cert.pem",
// 		PeerEndpoint: "dns:///localhost:9051",
// 		GatewayPeer:  "peer0.org2.example.com",
// 	}

// 	orgSetup, err := web.Initialize(orgConfig)
// 	if err != nil {
// 		log.Fatalf("Error initializing setup for Org2: %s", err)
// 	}

// 	fmt.Println("Server running at http://localhost:3000/")
// 	web.Serve(*orgSetup)
// }