package getRecord

import (
"errors"
"fmt"
"encoding/json"
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/hasher"
"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/database"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type RequestStruct struct {
	PatientId string
	UserId string
	RequesterName string
	Approve bool
	Date string
}

type Record struct {
	PatientName  string `json:patientName,omitempty`
	LastName  string `json:lastName`
	FirstName  string `json:firstName`
	DOB  string `json:dob`
	PatientID string `json:patientID`
	Name string `json:name`
	Prescription string `json:prescription`
	DoctorID string `json:DoctorID`
	Date string `json:,omitempty`
}

type Response struct {
	PatientName  string `json:patientName,omitempty`
	PatientID string `json:patientID`
	Prescription string `json:prescription`
	DoctorID string `json:DoctorID`
	Date string `json:,omitempty`
}

type Patient struct {
	PatientName string `json:,omitempty`
	FirstName string
	LastName string
	DOB string
}

func GetRecordClinic(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])
	var request Record
	err := json.Unmarshal(jsonBlob, &request)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	fmt.Printf("%+v", request)

	myPatient := Patient{
		FirstName: request.FirstName,
		LastName: request.LastName,
		DOB: request.DOB,
	}

	a, err := json.Marshal(myPatient)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	patientHash := hasher.Hash(a)

	fmt.Printf("\nHash: %s\n\n", patientHash)


	userId := request.DoctorID

	row,err := database.GetRow(stub,"requests",[]string{patientHash,userId})
	if err != nil {
		return nil, errors.New("Access Denied")
	}

	var requestStruct RequestStruct
	fmt.Printf("From getRecord: %v\n",row)
	requestStruct.PatientId = row.Columns[0].GetString_()
	requestStruct.UserId = row.Columns[1].GetString_()
	requestStruct.RequesterName = row.Columns[2].GetString_()
	requestStruct.Approve = row.Columns[3].GetBool()
	requestStruct.Date = row.Columns[4].GetString_()
	fmt.Printf("Struct From getRecord: %v\n",requestStruct)

	if requestStruct.Approve {
		return GetRecord(stub, patientHash)
	} else {
		fmt.Println("Denied :(")
		return nil, errors.New("Access Denied")
	}
	return nil, nil
}

func GetRecordPatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])
	var request Record
	err := json.Unmarshal(jsonBlob, &request)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	fmt.Printf("%+v", request)

	myPatient := Patient{
		FirstName: request.FirstName,
		LastName: request.LastName,
		DOB: request.DOB,
	}

	a, err := json.Marshal(myPatient)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	patientHash := hasher.Hash(a)

	fmt.Printf("\nHash: %s\n\n", patientHash)
	return GetRecord(stub, patientHash)
}

func GetRecord(stub shim.ChaincodeStubInterface, patientHash string) ([]byte, error) {
	rows,err := database.GetRows(stub,"records",[]string{patientHash})
	if err != nil {
		return nil, err
	}
	var records []Response
	var record Response
	for _, row := range rows {
		fmt.Printf("From getRequest: %v\n",row)
		record.PatientID = row.Columns[0].GetString_()
		record.PatientName = row.Columns[1].GetString_()
		record.Prescription = row.Columns[2].GetString_()
		record.DoctorID = row.Columns[3].GetString_()
		record.Date = row.Columns[4].GetString_()
		records = append(records,record)
	}
	response, err := json.Marshal(records)
	if err == nil {
		fmt.Printf("Marshal From getRequest: %v\n",string(response))
	}
	responseStr := "{\"Record\":"+string(response)+"}"
	response = []byte(responseStr)
	fmt.Printf("Marshal From getRequest: %v\n",string(response))

	fmt.Printf("Query Response:%s\n", response)
	return response, nil
}