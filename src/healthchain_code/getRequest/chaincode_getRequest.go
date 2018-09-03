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

package getRequest

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

type Request struct{
	ReqID int
	ReqName string
	RecordID int
	Date string
}

type Patient struct {
	PatientName string `json:,omitempty`
	FirstName string
	LastName string
	DOB string
}

type RequestStruct struct {
  PatientId string
  UserId string
  RequesterName string
  Approve bool
  Date string
}

func GetPatientRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])

	var patient Patient
	err := json.Unmarshal(jsonBlob, &patient)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%+v\n", patient)
	a, err := json.Marshal(patient)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	patientHash := hasher.Hash(a)

	fmt.Printf("\nHash: %s\n\n", patientHash)
	rows,err := database.GetRows(stub,"requests",[]string{patientHash})
	if err != nil {
		return nil, err
	}

	var requests []RequestStruct
	var requestStruct RequestStruct
	for _, row := range rows {
		requestStruct.PatientId = row.Columns[0].GetString_()
		requestStruct.UserId = row.Columns[1].GetString_()
		requestStruct.RequesterName = row.Columns[2].GetString_()
		requestStruct.Approve = row.Columns[3].GetBool()
		requestStruct.Date = row.Columns[4].GetString_()
		if !requestStruct.Approve{
			requests = append(requests,requestStruct)
		}
	}
	response, err := json.Marshal(requests)
	if err != nil {
		fmt.Println("error:", err)
	}
	responseStr := "{\"Request\":"+string(response)+"}"
	response = []byte(responseStr)
	return response, nil
}
