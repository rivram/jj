/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package requestRecords

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/hasher"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/database"
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

func RequestRecords(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])

	var request RequestStruct
	err := json.Unmarshal(jsonBlob, &request)
	if err != nil {
		fmt.Println("error:", err)
	}

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

	approveStatus := "false"
	userId := request.UserId
	requesterName := request.RequesterName
	t := time.Now()
	date := t.Format("2006-01-02T15:04")
	patientHash := hasher.Hash(a) //Hashed object of patient

	data := [6]string{patientHash, userId, requesterName, approveStatus, date}

	
	return nil,database.InsertRow(stub,"requests",data[0:5])
}
