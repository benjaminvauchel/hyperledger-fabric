package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// BaseCredential contains the common fields for all credentials
type BaseCredential struct {
	CredentialID		string `json:"CredentialID"`		// Credential unique identifier
	CredentialType		string `json:"CredentialType"`      // Type of credential: academic, professional
	FirstName         	string `json:"FirstName"`
	LastName          	string `json:"LastName"`
	Skills            	string `json:"Skills"`            	// List of skills associated with the experience/education (could be a comma-separated string or more complex structure)
	TalentID          	string `json:"TalentID"`         	// Talent identifier
	VerificationStatus 	string `json:"VerificationStatus"` 	// Status of the credential verification (e.g., "Pending", "Verified")
	VerifiedBy        	string `json:"VerifiedBy"`        	// Institution or admin that verified the credentials
}

// AcademicCredential is for academic credentials (e.g., degree, diploma)
type AcademicCredential struct {
	BaseCredential
	Education       	string `json:"Education"`      // e.g., "M.Eng. in Computer Science"
	Institution     	string `json:"Institution"`    // Institution granting the degree
}

// ProfessionalCredential is for professional credentials (e.g., work experience)
type ProfessionalCredential struct {
	BaseCredential
	Company         string `json:"Company"`        // Company where the work experience was gained
	WorkExperience  string `json:"WorkExperience"`
}

// InitLedger initializes the ledger with some sample talent credentials
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	credentials := []interface{}{
		AcademicCredential{
			BaseCredential: BaseCredential{
				CredentialID:		"credential1",
				TalentID:           "alicesmith01",
				FirstName:          "Alice",
				LastName:           "Smith",
				Skills:             "Python, Data Analysis",
				VerificationStatus: "Verified",
				VerifiedBy:         "Concordia University",
				CredentialType:		"academic",
			},
			Education:  "B.Sc. in Computer Science",
			Institution: "Concordia University",
		},
		ProfessionalCredential{
			BaseCredential: BaseCredential{
				CredentialID:		"credential2",
				TalentID:           "bobjohnson01",
				FirstName:          "Bob",
				LastName:           "Johnson",
				Skills:             "Project Management, Leadership",
				VerificationStatus: "Verified",
				VerifiedBy:         "Company ABCDEF",
				CredentialType:		"professional",
			},
			WorkExperience: "5 years as Project Manager",
			Company:        "Company ABCDEF",
		},
		AcademicCredential{
			BaseCredential: BaseCredential{
				CredentialID:		"credential3",
				TalentID:           "charliebrown02",
				FirstName:          "Charlie",
				LastName:           "Brown",
				Skills:             "Java, Software Engineering",
				VerificationStatus: "Pending",
				VerifiedBy:         "",
				CredentialType:		"academic",
			},
			Education:   "M.Sc. in Software Engineering",
			Institution: "Polytechnique Montréal",
		},
		ProfessionalCredential{
			BaseCredential: BaseCredential{
				CredentialID:		"credential4",
				TalentID:           "charliebrown02",
				FirstName:          "Charlie",
				LastName:           "Brown",
				Skills:             "C, C++, Python, Shell",
				VerificationStatus: "Pending",
				VerifiedBy:         "",
				CredentialType:		"professional",
			},
			WorkExperience: "4-Month Internship as a Software Developer",
			Company:        "Company XYZ",
		},
	}

	for _, credential := range credentials {
		credentialJSON, err := json.Marshal(credential)
		if err != nil {
			return fmt.Errorf("failed to marshal credential: %v", err)
		}

		// Store each credential in the ledger, using credentialID as the key
		credentialID := ""
		var baseCredential BaseCredential
		err = json.Unmarshal(credentialJSON, &baseCredential)
		if err != nil {
			return fmt.Errorf("failed to unmarshal credential: %v", err)
		}
		credentialID = baseCredential.CredentialID		

		err = ctx.GetStub().PutState(credentialID, credentialJSON)
		if err != nil {
			return fmt.Errorf("failed to put talent credential to world state. %v", err)
		}
	}

	return nil
}

// Issues a new academic credential
func (s *SmartContract) CreateAcademicCredential(ctx contractapi.TransactionContextInterface, credentialID string, talentID string, firstName string, lastName string, skills string, education string, institution string) error {
	exists, err := s.CredentialExists(ctx, credentialID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the credential %s already exists", credentialID)
	}

	academicCredential := AcademicCredential{
		BaseCredential: BaseCredential{
			CredentialID:		credentialID,
			TalentID:           talentID,
			FirstName:          firstName,
			LastName:           lastName,
			Skills:             skills,
			VerificationStatus: "Pending",
			CredentialType:		"academic",
			VerifiedBy:			"",
		},
		Education: education,
		Institution: institution,
	}

	academicCredentialJSON, err := json.Marshal(academicCredential)
	if err != nil {
		return fmt.Errorf("failed to marshal academic credential: %v", err)
	}

	return ctx.GetStub().PutState(credentialID, academicCredentialJSON)
}

// Issues a new professional credential
func (s *SmartContract) CreateProfessionalCredential(ctx contractapi.TransactionContextInterface, credentialID string, talentID string, firstName string, lastName string, skills string, workExperience string, company string) error {
	exists, err := s.CredentialExists(ctx, credentialID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the credential %s already exists", credentialID)
	}

	professionalCredential := ProfessionalCredential{
		BaseCredential: BaseCredential{
			CredentialID:		credentialID,
			TalentID:           talentID,
			FirstName:          firstName,
			LastName:           lastName,
			Skills:             skills,
			VerificationStatus: "Pending",
			CredentialType:		"professional",
			VerifiedBy:			"",
		},
		WorkExperience: workExperience,
		Company:        company,
	}

	professionalCredentialJSON, err := json.Marshal(professionalCredential)
	if err != nil {
		return fmt.Errorf("failed to marshal professional credential: %v", err)
	}

	return ctx.GetStub().PutState(credentialID, professionalCredentialJSON)
}

// CredentialExists returns true when credential with given ID exists in world state
func (s *SmartContract) CredentialExists(ctx contractapi.TransactionContextInterface, credentialID string) (bool, error) {
	credentialJSON, err := ctx.GetStub().GetState(credentialID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return credentialJSON != nil, nil
}

// GetBaseCredential retrieves the talent base credential by its ID
func (s *SmartContract) GetBaseCredential(ctx contractapi.TransactionContextInterface, credentialID string) (*BaseCredential, error) {
	talentCredentialJSON, err := ctx.GetStub().GetState(credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if talentCredentialJSON == nil {
		return nil, fmt.Errorf("the talent credential %s does not exist", credentialID)
	}

	var baseCredential BaseCredential
	err = json.Unmarshal(talentCredentialJSON, &baseCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal talent credential: %v", err)
	}

	return &baseCredential, nil
}

// GetAcademicCredential retrieves the academic credential by its ID
func (s *SmartContract) GetAcademicCredential(ctx contractapi.TransactionContextInterface, credentialID string) (*AcademicCredential, error) {
	talentCredentialJSON, err := ctx.GetStub().GetState(credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if talentCredentialJSON == nil {
		return nil, fmt.Errorf("the academic credential %s does not exist", credentialID)
	}

	var academicCredential AcademicCredential
	err = json.Unmarshal(talentCredentialJSON, &academicCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal academic credential: %v", err)
	}
	if academicCredential.CredentialType != "academic" {
		return nil, fmt.Errorf("the credential is not of type academic but %v", academicCredential.CredentialType)
	}

	return &academicCredential, nil
}

// GetProfessionalCredential retrieves the professional credential by its ID
func (s *SmartContract) GetProfessionalCredential(ctx contractapi.TransactionContextInterface, credentialID string) (*ProfessionalCredential, error) {
	talentCredentialJSON, err := ctx.GetStub().GetState(credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if talentCredentialJSON == nil {
		return nil, fmt.Errorf("the professional credential %s does not exist", credentialID)
	}

	var professionalCredential ProfessionalCredential
	err = json.Unmarshal(talentCredentialJSON, &professionalCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal professional credential: %v", err)
	}
	if professionalCredential.CredentialType != "professional" {
		return nil, fmt.Errorf("the credential is not of type professional but %v", professionalCredential.CredentialType)
	}

	return &professionalCredential, nil
}


// GetTalentCredential retrieves the entire talent credential by its ID (either academic or professional)
// But we lose the structure TODO: problem
func (s *SmartContract) GetTalentCredential(ctx contractapi.TransactionContextInterface, credentialID string) (interface{}, error) {
	talentCredentialJSON, err := ctx.GetStub().GetState(credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if talentCredentialJSON == nil {
		return nil, fmt.Errorf("the talent credential %s does not exist", credentialID)
	}

	var baseCredential BaseCredential
	err = json.Unmarshal(talentCredentialJSON, &baseCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal talent credential: %v", err)
	}

	// Determine if the credential is academic or professional
	if baseCredential.CredentialType == "academic" {
		var academicCredential AcademicCredential
		err = json.Unmarshal(talentCredentialJSON, &academicCredential)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal academic credential: %v", err)
		}
		return academicCredential, nil
	} else if baseCredential.CredentialType == "professional" {
		var professionalCredential ProfessionalCredential
		err = json.Unmarshal(talentCredentialJSON, &professionalCredential)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal professional credential: %v", err)
		}
		return professionalCredential, nil
	}

	return nil, fmt.Errorf("the talent credential type %v does not exist", baseCredential.CredentialType)
}

// Updates the verification status of a talent credential
func (s *SmartContract) UpdateVerificationStatus(ctx contractapi.TransactionContextInterface, credentialID string, status string, verifiedBy string) error {
	// Get the identity of the invoker (the user calling the smart contract)
	callerMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not get MSPID: %s", err)
	}

	// Only allow org1 (institutions) to approve or revoke
	if callerMSPID != "Org1MSP" {
		return fmt.Errorf("only members of Org1 (institutions) can approve or revoke credentials")
	}
	
	talentCredential, err := s.GetTalentCredential(ctx, credentialID)
	if err != nil {
		return err
	}

	// Type assertion: Determine the type of talentCredential and update the necessary fields
	switch v := talentCredential.(type) {
	case AcademicCredential:
		v.VerificationStatus = status
		v.VerifiedBy = verifiedBy

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated academic credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	case ProfessionalCredential:
		v.VerificationStatus = status
		v.VerifiedBy = verifiedBy

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated professional credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	default:
		return fmt.Errorf("unexpected credential type: %T", v)
	}
}

// Deletes a talent credential by its ID
func (s *SmartContract) DeleteTalentCredential(ctx contractapi.TransactionContextInterface, credentialID string) error {
	talentCredential, err := s.GetTalentCredential(ctx, credentialID)
	if err != nil {
		return err
	}

	if talentCredential == nil {
		return fmt.Errorf("the talent credential %s does not exist", credentialID)
	}

	// Delete the credential from the ledger
	return ctx.GetStub().DelState(credentialID)
}

// Updates the skills of a talent credential
func (s *SmartContract) UpdateSkills(ctx contractapi.TransactionContextInterface, credentialID string, newSkills string) error {
	talentCredential, err := s.GetTalentCredential(ctx, credentialID)
	if err != nil {
		return err
	}

	// Type assertion: Determine the type of talentCredential and update the skills
	switch v := talentCredential.(type) {
	case AcademicCredential:
		v.Skills = newSkills

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated academic credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	case ProfessionalCredential:
		v.Skills = newSkills

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated professional credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	default:
		return fmt.Errorf("unexpected credential type: %T", v)
	}
}

// Updates the first and last name of a talent credential (if the talent made an error)
func (s *SmartContract) UpdateName(ctx contractapi.TransactionContextInterface, credentialID string, newFirstName string, newLastName string) error {
	talentCredential, err := s.GetTalentCredential(ctx, credentialID)
	if err != nil {
		return err
	}

	// Type assertion: Determine the type of talentCredential and update the name fields
	switch v := talentCredential.(type) {
	case AcademicCredential:
		v.FirstName = newFirstName
		v.LastName = newLastName

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated academic credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	case ProfessionalCredential:
		v.FirstName = newFirstName
		v.LastName = newLastName

		updatedCredentialJSON, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal updated professional credential: %v", err)
		}

		return ctx.GetStub().PutState(credentialID, updatedCredentialJSON)

	default:
		return fmt.Errorf("unexpected credential type: %T", v)
	}
}


// GetAllCredentials retrieves all credentials (both academic and professional) from the ledger
func (s *SmartContract) GetAllCredentials(ctx contractapi.TransactionContextInterface) ([]byte, error) { //TODO: before []interface{}, now []byte
	// Perform a range query with an empty string for startKey and endKey for an open-ended query
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var credentials []interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Check the credential type by unmarshalling into a BaseCredential first
		var baseCredential BaseCredential
		err = json.Unmarshal(queryResponse.Value, &baseCredential)
		if err != nil {
			return nil, err
		}

		// Determine whether the credential is academic or professional
		var credential interface{}
		if baseCredential.CredentialType == "academic" {
			// AcademicCredential
			var academicCredential AcademicCredential
			err = json.Unmarshal(queryResponse.Value, &academicCredential)
			if err != nil {
				return nil, err
			}
			credential = academicCredential
		} else if baseCredential.CredentialType == "professional" {
			// ProfessionalCredential
			var professionalCredential ProfessionalCredential
			err = json.Unmarshal(queryResponse.Value, &professionalCredential)
			if err != nil {
				return nil, err
			}
			credential = professionalCredential
		} else {
			// If it doesn't match either, we handle it as a base credential
			credential = baseCredential
		}

		// Append the found credential to the list
		credentials = append(credentials, credential)
	}

	// Marshal the final array to JSON bytes
	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil

	// return credentials, nil
}

