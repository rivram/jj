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

package main

import (
	"errors"
	"fmt"
	// "strconv"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/approveTx"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/createRecord"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/requestRecords"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/getRecord"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/getRequest"
	// "github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/hasher"
	"github.com/hyperledger/fabric/examples/chaincode/go/healthchain_code/database"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	err = database.DeleteTable(stub,"requests")
	if err != nil {
		return nil, fmt.Errorf("Error deleting table during init. %s", err)
	}
	err = database.CreateRequestsTable(stub)
	if err != nil {
		return nil, fmt.Errorf("Error creating table during init. %s", err)
	}

	err = database.DeleteTable(stub,"records")
	if err != nil {
		return nil, fmt.Errorf("Error deleting table during init. %s", err)
	}
	err = database.CreateRecordsTable(stub)
	if err != nil {
		return nil, fmt.Errorf("Error creating table during init. %s", err)
	}

	/************
			// Write the state to the ledger
			err = stub.PutState(A, []byte(strconv.Itoa(Aval))
			if err != nil {
				return nil, err
			}

			stub.PutState(B, []byte(strconv.Itoa(Bval))
			err = stub.PutState(B, []byte(strconv.Itoa(Bval))
			if err != nil {
				return nil, err
			}
	************/
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	switch function {
		// Create a new Record
		case "createRecord":
			return createRecord.CreateRecord(stub, args)
		case "requestRecords":
			return requestRecords.RequestRecords(stub,args)
		case "approveTx":
			return approveTx.ApproveTx(stub,args)
		case "insertRow":
			return nil, database.InsertRow(stub,"requests",args)
		case "updateRow":
			return database.UpdateRow(stub,"requests",args)
		case "deleteRow":
			return database.DeleteRow(stub,"requests",args)
		default:
			return nil,errors.New("Expecting a Function Name, Invoke doesn't know what to do")
	}
	return nil, err
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "getRequest":
		return getRequest.GetPatientRequest(stub,args)
	case "getRecordClinic":
		return getRecord.GetRecordClinic(stub,args)
	case "getRecordPatient":
		return getRecord.GetRecordPatient(stub,args)
	case "getRow":
		row,err := database.GetRow(stub,"requests",args)
		if err != nil {
			fmt.Printf("error: %s\n",err)
		}
		fmt.Printf("Printing the row MAIN: %v\n",row)
		return nil,nil
	case "getRows":
		rows,err := database.GetRows(stub,"requests",args)
		if err != nil {
			fmt.Printf("error: %s\n",err)
		}
		fmt.Printf("Printing all rows MAIN: %v\n",rows)
		// fmt.Printf("Printing first row MAIN: %v\n",rows[0])
		// fmt.Printf("Printing first row approve MAIN: %v\n",rows[0].Columns[3].GetBool())
		return nil,nil
	default:
		return nil,errors.New("Expecting a Function Name, Invoke doesn't know what to do")
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
