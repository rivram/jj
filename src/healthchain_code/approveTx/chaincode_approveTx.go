package approveTx

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/database"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/hasher"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type RequestStruct struct {
	UserId string `json:,omitempty`
	PatientId string `json:,omitempty`
	LastName  string `json:lastName`
	FirstName  string `json:firstName`
	DOB  string `json:dob`
	RequesterName string `json:,omitempty`
	Approve bool
	Date string `json:,omitempty`
}

type Patient struct {
	PatientName string `json:,omitempty`
	FirstName string
	LastName string
	DOB string
}

func ApproveTx(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])
	var request RequestStruct
	err := json.Unmarshal(jsonBlob, &request)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("From approveTx%+v", request)

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
	approveStatus := request.Approve
	requesterName := request.RequesterName
	userId := request.UserId
	date := request.Date

	data := [5]string{patientHash, userId, requesterName, strconv.FormatBool(approveStatus), date}

	if approveStatus {
		return database.UpdateRow(stub,"requests",data[0:5])
	} else {
		return database.DeleteRow(stub,"requests",data[0:5])
	}
}