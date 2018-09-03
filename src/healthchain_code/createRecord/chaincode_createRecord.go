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

package createRecord

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/database"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/hasher"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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

type Patient struct {
	PatientName string `json:,omitempty`
	FirstName string
	LastName string
	DOB string
}

func CreateRecord(stub shim.ChaincodeStubInterface,args []string) ([]byte,error)  {
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting more arguments")
	}

	var jsonBlob = []byte(args[0])
	var patient Record
	err := json.Unmarshal(jsonBlob, &patient)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	myPatient := Patient{
		FirstName: patient.FirstName,
		LastName: patient.LastName,
		DOB: patient.DOB,
	}

	a, err := json.Marshal(myPatient)
	if err != nil {
		fmt.Println("error:", err)
		return nil,err
	}

	patientHash := hasher.Hash(a)
	patientId := patientHash
	patientName := patient.FirstName + " " +patient.LastName
	prescription := patient.Prescription
	doctorId := patient.DoctorID
	t := time.Now()
	date := t.Format("2006-01-02T15:04")

	data := []string{patientId, patientName, prescription, doctorId, date}

	err = database.InsertRow(stub,"records",data[0:5])
	if err != nil {
		return nil, err
	}

	return nil, nil
}
