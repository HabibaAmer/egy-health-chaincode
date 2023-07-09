package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract provides functions for managing a Patient
type SmartContract struct {
	// contractapi.Contract
}

// PatientData describes basic details of what makes up a simple Patient

type CounterNO struct {
	Counter int `json:"counter"`
}

type User struct {
	Username string `json:"Username"`
	UserID   string `json:"UserID"`
	Email    string `json:"Email"`
	UserRole string `json:"UserRole"`
	Password string `json:"Password"`
}

type PatientData struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Age       int             `json:"age"`
	Gender    string          `json:"gender"`
	BloodType string          `json:"bloodType"`
	Allergies string          `json:"allergies"`
	Access    map[string]bool `json:"Doctor Access"`
	Record    MedicalRecord   `json:"record"`
}
type MedicalRecord struct {
	Diagnose           string   `json:"diagnose"`
	Medications        string   `json:"medications"`
	DiagnosesHistory   []string `json:"diagnoseshistory"`
	MedicationsHistory []string `json:"medicationhistory"`
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initLedger" {
		//init ledger
		return t.InitLedger(stub)
	} else if function == "signIn" {
		//login user
		return t.signIn(stub, args)
	} else if function == "createUser" {
		//create a new user
		return t.createUser(stub, args)

	} else if function == "createPatient" {
		//create a new patient
		age, _ := strconv.Atoi(args[2])

		access := make(map[string]bool)
		json.Unmarshal([]byte(args[6]), &access)
		fmt.Printf("access %v", access)

		diagnosehistory := []string{}
		json.Unmarshal([]byte(args[9]), &diagnosehistory)
		fmt.Printf("diagnosehistory %v", diagnosehistory)

		medicationhistory := []string{}
		json.Unmarshal([]byte(args[10]), &medicationhistory)
		fmt.Printf("medicationhistory %v", medicationhistory)

		return t.CreatePatient(stub, args[0], args[1], age, args[3], args[4], args[5], access, args[7], args[8], diagnosehistory, medicationhistory)
	} else if function == "updatePatientMedicalRecords" {
		// update patient medcial records
		return t.UpdatePatientMedicalRecords(stub, args[0], args[1], args[2], args[3])
	} else if function == "grantAccess" {
		// grant access to a provider
		return t.GrantAccess(stub, args[0], args[1], args[2])
	} else if function == "revokeAccess" {
		// revoke access from a provider
		return t.RevokeAccess(stub, args[0], args[1], args[2])
	} else if function == "readPatientMedicalInfo" {
		// read the patient medical records
		return t.ReadPatientMedicalInfo(stub, args[0], args[1])
	} else if function == "readPatientAllInfo" {
		// read all the patient's info
		return t.ReadPatientAllInfo(stub, args[0], args[1])
	} else if function == "deletePatient" {
		// delete the patient EHR
		return t.DeletePatient(stub, args[0])
	}
	// else if function == "patientExists" {
	// 	// check if the patient exists
	// 	return t.PatientExists(stub, args[0])
	// }
	//  else if function == "hasPermission" {
	// 	// check if the provider has permission to access the EHR
	// 	return t.hasPermission(stub, args[0], args[1])
	// }

	fmt.Println("invoke did not find func: " + function)
	//error
	return shim.Error("Received unknown function invocation")
}

func (t *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	// Initializing User Counter
	UserCounterBytes, _ := APIstub.GetState("UserCounterNO")
	if UserCounterBytes == nil {
		var UserCounter = CounterNO{Counter: 0}
		UserCounterBytes, _ := json.Marshal(UserCounter)
		err := APIstub.PutState("UserCounterNO", UserCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Intitate User Counter"))
		}
	}

	return shim.Success([]byte("Initiated successfully."))
}

//getCounter to the latest value of the counter based on the Asset Type provided as input parameter
func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Sprintf("Counter Current Value %d of Asset Type %s", counterAsset.Counter, AssetType)

	return counterAsset.Counter
}

//incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Sprintf("Failed to Increment Counter")

	}

	fmt.Println("Success in incrementing counter  %v", counterAsset)

	return counterAsset.Counter
}

// InitLedger adds a base set of Patients to the ledger --> The init function is called when the smart contract is first deployed to the network
func (s *SmartContract) InitLedger(APIstub shim.ChaincodeStubInterface) pb.Response {

	Patients := []PatientData{
		{ID: "Patient1", Name: "test1", Age: 5, Gender: "male", BloodType: "B+", Allergies: "xx", Record: MedicalRecord{"Diagnose1", "Medications1", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient2", Name: "test2", Age: 5, Gender: "male", BloodType: "A+", Allergies: "yy", Record: MedicalRecord{"Diagnose2", "Medications2", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
		{ID: "Patient3", Name: "test3", Age: 10, Gender: "female", BloodType: "AB", Allergies: "cc", Record: MedicalRecord{"Diagnose3", "Medications3", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient4", Name: "test4", Age: 10, Gender: "female", BloodType: "O+", Allergies: "nn", Record: MedicalRecord{"Diagnose4", "Medications4", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
		{ID: "Patient5", Name: "test5", Age: 15, Gender: "female", BloodType: "O-", Allergies: "mm", Record: MedicalRecord{"Diagnose5", "Medications5", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient6", Name: "test6", Age: 15, Gender: "female", BloodType: "B-", Allergies: "jj", Record: MedicalRecord{"Diagnose6", "Medications6", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
	}

	for _, Patient := range Patients {
		PatientJSON, err := json.Marshal(Patient) // take each patient and convert it to json format then store this format in PatientJSON file and check for error
		if err != nil {                           //if error field is not empty
			return shim.Error("Error parsing patient to json") //then return the error and return to main and exit the init function
		}

		err = APIstub.PutState(Patient.ID, PatientJSON) // used to access the BC ledger , PutState: generate the key-value pair==(Patient.ID-PatientJSON) , GetStub: provide APIs access to the world state
		// = is used for assigment while := is used for decleration
		if err != nil {
			fmt.Errorf("failed to put to world state. %v", err)
			return shim.Error("failed to put to world state")
		}
	}

	return shim.Success([]byte("initiated ledger successfully."))
}

func (t *SmartContract) signIn(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expected 2 argument")
	}

	if len(args[0]) == 0 {
		return shim.Error("User ID must be provided")
	}

	if len(args[1]) == 0 {
		return shim.Error("Password must be provided")
	}

	entityUserBytes, _ := APIstub.GetState("user_" + args[0])
	if entityUserBytes == nil {
		return shim.Error("Cannot Find User")
	}
	entityUser := User{}
	// unmarsahlling the entity data
	json.Unmarshal(entityUserBytes, &entityUser)

	// check if password matched
	if entityUser.Password != args[1] {
		return shim.Error("Either id or password is wrong")
	}

	return shim.Success(entityUserBytes)
}

//create user
func (t *SmartContract) createUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, Required 5 arguments")
	}

	if len(args[0]) == 0 {
		return shim.Error("Name must be provided to register user")
	}

	if len(args[1]) == 0 {
		return shim.Error("Email is mandatory")
	}

	if len(args[2]) == 0 {
		return shim.Error("User role must be specified")
	}

	if len(args[3]) == 0 {
		return shim.Error("Password must be non-empty ")
	}

	userCounter := getCounter(APIstub, "UserCounterNO")
	userCounter++

	var createdUser = User{Username: args[0], UserID: "User" + strconv.Itoa(userCounter), Email: args[1], UserRole: args[2], Password: args[3]}

	createdUserAsBytes, errMarshal := json.Marshal(createdUser)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Product: %s", errMarshal))
	}

	errPut := APIstub.PutState("user_"+createdUser.UserID, createdUserAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to register user: %s", createdUser.UserID))
	}

	//TO Increment the User Counter
	incrementCounter(APIstub, "UserCounterNO")

	fmt.Println("User register successfully %v", createdUser)

	return shim.Success(createdUserAsBytes)

}

func (s *SmartContract) CreatePatient(APIstub shim.ChaincodeStubInterface, patientID string, name string, age int, gender string, bloodType string, allergies string, access map[string]bool, diagnose string, medication string, diagnosehistory []string, medicationhistory []string) pb.Response {
	exists, err := s.PatientExists(APIstub, patientID)

	if err != nil {
		fmt.Errorf("%s", err)
		return shim.Error("")
	}
	if exists {
		fmt.Errorf("the Patient %s already exists", patientID)
		return shim.Error("The patient already exists")
	}

	NewPatient := PatientData{
		ID:        patientID,
		Name:      name,
		Age:       age,
		Gender:    gender,
		BloodType: bloodType,
		Allergies: allergies,
		Access:    access,
		Record:    MedicalRecord{Diagnose: diagnose, Medications: medication, DiagnosesHistory: diagnosehistory, MedicationsHistory: medicationhistory},
	}
	PatientJSON, err := json.Marshal(NewPatient)
	if err != nil { //if error field is not empty
		fmt.Errorf("%s", err)
		return shim.Error("Failed to add patient") //then return the error and return to main and exit the init function
	}

	err = APIstub.PutState(patientID, PatientJSON) // used to access the BC ledger , PutState: generate the key-value pair==(Patient.ID-PatientJSON) , GetStub: provide APIs access to the world state
	// = is used for assigment while := is used for decleration
	if err != nil {
		fmt.Errorf("failed to put to world state. %v", err)
		return shim.Error("Failed to put to world state")
	}
	return shim.Success(PatientJSON)
}

func (s *SmartContract) PatientExists(APIstub shim.ChaincodeStubInterface, id string) (bool, error) {
	PatientJSON, err := APIstub.GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return PatientJSON != nil, nil
}

func (s *SmartContract) UpdatePatientMedicalRecords(APIstub shim.ChaincodeStubInterface, providerID string, patientID string, diagnose string, medication string) pb.Response {

	patientRecord, err := APIstub.GetState(patientID) // get state of ID from the ledger and store it in PatientRecords , if an error occur store it in err
	if err != nil {
		fmt.Errorf("failed to read from world state: %v", err)
		return shim.Error("failed to read from world state")
	}
	if patientRecord == nil {
		fmt.Errorf("the patient with ID %s does not exist", patientID)
		return shim.Error("the patient does not exist")
	}

	access, err := s.hasPermission(APIstub, providerID, patientID)

	if access == false && err != nil {
		fmt.Errorf("Provider does not have permission to the update data for patient %s", patientID)
		return shim.Error("Povider doesn't have permission to update data for the patient")
	}

	var patientData PatientData //declare updated data as a variable
	err = json.Unmarshal(patientRecord, &patientData)
	if err != nil {
		fmt.Errorf("failed to unmarshal patient record: %v", err)
		return shim.Error("failed to unmarshal patient record")
	}

	//add the new diagnose and medication to the history
	patientData.Record.DiagnosesHistory = append(patientData.Record.DiagnosesHistory, patientData.Record.Diagnose)
	patientData.Record.MedicationsHistory = append(patientData.Record.MedicationsHistory, patientData.Record.Medications)
	//overwrite the diagnose and medication fields
	patientData.Record.Diagnose = diagnose
	patientData.Record.Medications = medication

	updatedpatientData, err := json.Marshal(patientData)
	if err != nil {
		fmt.Errorf("failed to marshal updated patient record: %v", err)
		return shim.Error("failed to marshal updated patient record")
	}

	err = APIstub.PutState(patientID, updatedpatientData)
	if err != nil {
		fmt.Errorf("failed to update patient record: %v", err)
		return shim.Error("failed to update patient record")
	}

	return shim.Success(updatedpatientData)
}

// ReadPatient returns the Medical info only stored in the world state with given id.
func (s *SmartContract) ReadPatientMedicalInfo(APIstub shim.ChaincodeStubInterface, providerID string, patientID string) pb.Response {

	PatientRecordJSON, err := APIstub.GetState(patientID)
	if err != nil {
		fmt.Errorf("failed to read from world state: %v", err)
		return shim.Error("failed to read from world state")

	}
	if PatientRecordJSON == nil {
		fmt.Errorf("the patient %s does not exist", patientID)
		return shim.Error("the patient does not exist")

	}
	access, err := s.hasPermission(APIstub, providerID, patientID)

	if access == false && err != nil {
		fmt.Errorf("Provider does not have permission to read data for patient %s", patientID)
		return shim.Error("Provider does not have permission to read data for patient")
	}

	var Patient PatientData
	err = json.Unmarshal(PatientRecordJSON, &Patient)
	if err != nil {
		fmt.Errorf("failed to unmarshal patient data %s", err)
		return shim.Error("failed to unmarshal patient data")
	}

	fmt.Printf("Diagnosis: %s  , Medication: %s \n", Patient.Record.Diagnose, Patient.Record.Medications)
	fmt.Println("Diagnoses history:")
	for _, diagnosis := range Patient.Record.DiagnosesHistory {
		fmt.Printf("- %s\n", diagnosis)
	}
	fmt.Println("Medication history:")
	for _, history := range Patient.Record.MedicationsHistory {
		fmt.Printf("- %s\n", history)
	}

	return shim.Success(PatientRecordJSON)
}

// ReadPatient returns the all patient info stored in the world state with given id.
func (s *SmartContract) ReadPatientAllInfo(APIstub shim.ChaincodeStubInterface, providerID string, patientID string) pb.Response {

	PatientJSON, err := APIstub.GetState(patientID)
	if err != nil {
		fmt.Errorf("failed to read from world state: %v", err)
		return shim.Error("failed to read from world state")
	}
	if PatientJSON == nil {
		fmt.Errorf("the patient %s does not exist", patientID)
		return shim.Error("the patient does not exist")
	}

	access, err := s.hasPermission(APIstub, providerID, patientID)

	if access == false && err != nil {
		fmt.Errorf("Provider does not have permission to read data for patient %s", patientID)
		return shim.Error("Provider does not have permission to read data for patient")
	}
	var Patient PatientData
	err = json.Unmarshal(PatientJSON, &Patient)
	if err != nil {
		fmt.Errorf("failed to unmarshal patient %s", err)
		return shim.Error("failed to unmarshal patient")
	}

	return shim.Success(PatientJSON)
}

// DeletePatient deletes an given patient from the world state.
func (s *SmartContract) DeletePatient(APIstub shim.ChaincodeStubInterface, id string) pb.Response {
	exists, err := s.PatientExists(APIstub, id)
	if err != nil {
		fmt.Errorf("failed to check patient exist: %v", err)
		return shim.Error("failed to check patient exist")
	}
	if !exists {
		fmt.Errorf("the Patient %s does not exist", id)
		return shim.Error("the Patient does not exist")
	}

	APIstub.DelState(id)
	return shim.Success([]byte("Patient deleted successfully."))
}

// Define the grantAccess function

func (p *SmartContract) GrantAccess(APIstub shim.ChaincodeStubInterface, userID string, patientID string, providerID string) pb.Response {
	// Check if the caller is authorized to grant access

	// Retrieve patient data from the ledger
	patientDataJSON, err := APIstub.GetState(patientID)
	if err != nil {
		fmt.Errorf("failed to read patient data from world state: %v", err)
		return shim.Error("failed to read patient data from world")

	}
	if patientDataJSON == nil {
		fmt.Errorf("patient data with ID %s does not exist", patientID)
		return shim.Error("patient data does not exist")
	}

	access, err := p.hasPermission(APIstub, userID, patientID)

	if access == false && err != nil {
		fmt.Errorf("user does not have permission to change the AccessList of the Patient %s", patientID)
		return shim.Error("user does not have permission to change the AccessList of the Patient")
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
		return shim.Error("failed to unmarshal patient data JSON")
	}

	access, ok := patientData.Access[providerID]
	if access == false || ok == false {
		patientData.Access = map[string]bool{providerID: true}
	}

	updatedpatientData, err := json.Marshal(patientData)
	if err != nil {
		fmt.Errorf("failed to marshal updated patient record: %v", err)
		return shim.Error("failed to marshal updated patient record")
	}

	err = APIstub.PutState(patientID, updatedpatientData)
	if err != nil {
		fmt.Errorf("failed to update patient record: %v", err)
		return shim.Error("failed to update patient record")
	}

	return shim.Success(updatedpatientData)
}

// define Revoke Access function
func (s *SmartContract) RevokeAccess(APIstub shim.ChaincodeStubInterface, userID string, patientID string, DoctorID string) pb.Response {

	// Retrieve patient data from the ledger
	patientDataJSON, err := APIstub.GetState(patientID)
	if err != nil {
		fmt.Errorf("failed to read patient data from world state: %v", err)
		return shim.Error("failed to read patient data from world state")
	}
	if patientDataJSON == nil {
		fmt.Errorf("patient data with ID %s does not exist", patientID)
		return shim.Error("patient data does not exist")
	}

	access, err := s.hasPermission(APIstub, userID, patientID)

	if access == false && err != nil {
		fmt.Errorf("user does not have permission to change the AccessList of the Patient %s", patientID)
		return shim.Error("user does not have permission to change the AccessList of the Patient")
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
		return shim.Error("failed to unmarshal patient data JSON")
	}

	// Check if doctor has access to patient data
	if _, ok := patientData.Access[DoctorID]; !ok {
		fmt.Errorf("Doctor with ID %s does not have access to patient data with ID %s", DoctorID, patientID)
		return shim.Error("Doctor with does not have access to patient data")
	}

	// Revoke access for specific hospital
	//patientData.Access[DoctorID] = false
	patientData.Access = map[string]bool{DoctorID: false}

	// Marshal updated patient data into JSON format
	patientDataJSON, err = json.Marshal(patientData)
	if err != nil {
		fmt.Errorf("failed to marshal updated patient data into JSON format: %v", err)
		return shim.Error("failed to marshal updated patient data into JSON format")
	}

	// Update patient data in the ledger
	err = APIstub.PutState(patientID, patientDataJSON)
	if err != nil {
		fmt.Errorf("failed to update patient data in the ledger: %v", err)
		return shim.Error("failed to update patient data in the ledger")
	}

	return shim.Success(patientDataJSON)
}

// define share data function

func (s *SmartContract) hasPermission(APIstub shim.ChaincodeStubInterface, userID string, patientID string) (bool, error) {
	// Check if the user has permission to access the patient data
	// In this example implementation, we assume that only the patient and healthcare providers with the patient's permission have access to the data
	if userID == patientID {
		return true, nil
	}

	// Retrieve patient data from the ledger
	patientDataJSON, err := APIstub.GetState(patientID)
	if err != nil {
		return false, fmt.Errorf("failed to read patient data from world state: %v", err)
	}
	if patientDataJSON == nil {
		return false, fmt.Errorf("patient data with ID %s does not exist", patientID)
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
	}

	access, ok := patientData.Access[userID]
	if !ok {
		return false, fmt.Errorf("%s doesn't exist in %s access list", userID, patientID)
	}
	if access == false {
		return false, fmt.Errorf("%s doesn't have access to %s", userID, patientID)
	}

	return true, nil
}
