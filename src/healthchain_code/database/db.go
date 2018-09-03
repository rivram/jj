package database

import (
"github.com/hyperledger/fabric/core/chaincode/shim"
"fmt"
"errors"
"strconv"
"encoding/json"
)

// GetTable
func CreateRecordsTable(stub shim.ChaincodeStubInterface) error {
  // Create table one

  var columnDefsRecords []*shim.ColumnDefinition
  columnOneTableOneDef := shim.ColumnDefinition{
    Name: "patientId",
    Type: shim.ColumnDefinition_STRING,
    Key: true,
  }
  columnTwoTableOneDef := shim.ColumnDefinition{
    Name: "patientName",
    Type: shim.ColumnDefinition_STRING,
    Key: false,
  }
  columnThreeTableOneDef := shim.ColumnDefinition{
    Name: "prescription",
    Type: shim.ColumnDefinition_STRING,
    Key: false,
  }
  columnFourTableOneDef := shim.ColumnDefinition{
    Name: "doctorId",
    Type: shim.ColumnDefinition_STRING,
    Key: false,
  }
  columnFiveTableOneDef := shim.ColumnDefinition{
    Name: "date",
    Type: shim.ColumnDefinition_STRING,
    Key: true,
  }
  columnDefsRecords = append(columnDefsRecords, &columnOneTableOneDef)
  columnDefsRecords = append(columnDefsRecords, &columnTwoTableOneDef)
  columnDefsRecords = append(columnDefsRecords, &columnThreeTableOneDef)
  columnDefsRecords = append(columnDefsRecords, &columnFourTableOneDef)
  columnDefsRecords = append(columnDefsRecords, &columnFiveTableOneDef)
  return stub.CreateTable("records", columnDefsRecords)
}

func CreateRequestsTable(stub shim.ChaincodeStubInterface) error {

  var columnDefsRequests []*shim.ColumnDefinition
  columnOneTableOneDef := shim.ColumnDefinition{
    Name: "patientId",
    Type: shim.ColumnDefinition_STRING,
    Key: true,
  }
  columnTwoTableOneDef := shim.ColumnDefinition{
    Name: "userId",
    Type: shim.ColumnDefinition_STRING,
    Key: true,
  }
  columnThreeTableOneDef := shim.ColumnDefinition{
    Name: "requesterName",
    Type: shim.ColumnDefinition_STRING,
    Key: false,
  }
  columnFourTableOneDef := shim.ColumnDefinition{
    Name: "approve",
    Type: shim.ColumnDefinition_BOOL,
    Key: false,
  }
  columnFiveTableOneDef := shim.ColumnDefinition{
    Name: "date",
    Type: shim.ColumnDefinition_STRING,
    Key: false,
  }
  columnDefsRequests = append(columnDefsRequests, &columnOneTableOneDef)
  columnDefsRequests = append(columnDefsRequests, &columnTwoTableOneDef)
  columnDefsRequests = append(columnDefsRequests, &columnThreeTableOneDef)
  columnDefsRequests = append(columnDefsRequests, &columnFourTableOneDef)
  columnDefsRequests = append(columnDefsRequests, &columnFiveTableOneDef)
  return stub.CreateTable("requests", columnDefsRequests)
}

// func GetRowRecords(stub shim.ChaincodeStubInterface, tableName string, key []shim.Column) (shim.Row,error) {
func GetRow(stub shim.ChaincodeStubInterface, tableName string, args []string) (shim.Row, error) {
  var row shim.Row
  if len(args) != 2 {
    return row, fmt.Errorf("Must have two args")
  }

  var columns []shim.Column
  col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
  col2 := shim.Column{Value: &shim.Column_String_{String_: args[1]}}
  columns = append(columns, col1)
  columns = append(columns, col2)

  row, err := stub.GetRow(tableName, columns)
  if err != nil {
    return row, fmt.Errorf("getRowRecords operation failed. %s", err)
  }

  rowBytes, err := json.Marshal(row)
  if err != nil {
    return row, fmt.Errorf("Marshal error: %s", err)
  }
  fmt.Printf("rowBytes: %v\n",rowBytes)
  if len(rowBytes) <= 3 {
    return row, fmt.Errorf("No such row exists %s", err)
  }
  // fmt.Printf("Juust Printing the row: %v\n",row)
  return row, nil
}

// Function to get multiple records
func GetRows(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]shim.Row, error) {
  if len(args) != 1 {
    return nil, fmt.Errorf("Must have one arg")
  }
  var columns []shim.Column
  col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
  columns = append(columns, col1)

  rowChannel, err := stub.GetRows(tableName, columns)
  if err != nil {
    return nil, fmt.Errorf("getRowRecords operation failed. %s", err)
  }
  var rows []shim.Row
  for {
    select {
    case row, ok := <-rowChannel:
      if !ok {
        rowChannel = nil
      } else {
        rows = append(rows, row)
      }
    }
    if rowChannel == nil {
      break
    }
  }
  return rows, nil
}

func InsertRow(stub shim.ChaincodeStubInterface, tableName string, args []string) error  {
  switch tableName {
  case "requests":
    if len(args) < 5 {
      return errors.New("insertRowRecords failed. Must include 5 column values")
    }
    fmt.Println(args)

    patientId := args[0]
    userId := args[1]
    requesterName := args[2]
    approve, err := strconv.ParseBool(args[3])
    if err != nil {
      return errors.New("insertRowRequests failed. arg[4] must be convertable to bool")
    }
    date := args[4]

    fmt.Println(patientId,userId,requesterName,approve,date)

    var columns []*shim.Column
    col1 := shim.Column{Value: &shim.Column_String_{String_: patientId}}
    col2 := shim.Column{Value: &shim.Column_String_{String_: userId}}
    col3 := shim.Column{Value: &shim.Column_String_{String_: requesterName}}
    col4 := shim.Column{Value: &shim.Column_Bool{Bool: approve}}
    col5 := shim.Column{Value: &shim.Column_String_{String_: date}}
    columns = append(columns, &col1)
    columns = append(columns, &col2)
    columns = append(columns, &col3)
    columns = append(columns, &col4)
    columns = append(columns, &col5)
    fmt.Println(columns)

    row := shim.Row{Columns: columns}
    fmt.Println(row)
    ok, err := stub.InsertRow(tableName, row)
    if err != nil {
      return fmt.Errorf("insertRowRequests operation failed. %s", err)
    }
    if !ok {
      return errors.New("insertRowRequests operation failed. Row with given key already exists")
    }
  case "records":
    if len(args) < 5 {
      return errors.New("insertRowRecords failed. Must include 5 column values")
    }
    fmt.Println(args)
    patientId := args[0]
    patientName := args[1]
    prescription := args[2]
    doctorId := args[3]
    date := args[4]

    var columns []*shim.Column
    col1 := shim.Column{Value: &shim.Column_String_{String_: patientId}}
    col2 := shim.Column{Value: &shim.Column_String_{String_: patientName}}
    col3 := shim.Column{Value: &shim.Column_String_{String_: prescription}}
    col4 := shim.Column{Value: &shim.Column_String_{String_: doctorId}}
    col5 := shim.Column{Value: &shim.Column_String_{String_: date}}
    columns = append(columns, &col1)
    columns = append(columns, &col2)
    columns = append(columns, &col3)
    columns = append(columns, &col4)
    columns = append(columns, &col5)
    fmt.Println(columns)
    row := shim.Row{Columns: columns}
    fmt.Println(row)
    ok, err := stub.InsertRow(tableName, row)
    if err != nil {
      return fmt.Errorf("insertRowRequests operation failed. %s", err)
    }
    if !ok {
      return errors.New("insertRowRequests operation failed. Row with given key already exists")
    }
  default :
    return errors.New("insertRowRequests failed. No such table.")
  }
  return nil
}

func UpdateRow(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]byte, error)  {
  switch tableName {
  case "requests":
    if len(args) < 5 {
      return nil, errors.New("insertRowRecords failed. Must include 5 column values")
    }
    fmt.Println(args)

    patientId := args[0]
    userId := args[1]
    requesterName := args[2]
    approve, err := strconv.ParseBool(args[3])
    if err != nil {
      return nil, errors.New("insertRowRequests failed. arg[3] must be convertable to bool")
    }
    date := args[4]

    fmt.Println(patientId,userId,requesterName,approve,date)

    var columns []*shim.Column
    col1 := shim.Column{Value: &shim.Column_String_{String_: patientId}}
    col2 := shim.Column{Value: &shim.Column_String_{String_: userId}}
    col3 := shim.Column{Value: &shim.Column_String_{String_: requesterName}}
    col4 := shim.Column{Value: &shim.Column_Bool{Bool: approve}}
    col5 := shim.Column{Value: &shim.Column_String_{String_: date}}
    columns = append(columns, &col1)
    columns = append(columns, &col2)
    columns = append(columns, &col3)
    columns = append(columns, &col4)
    columns = append(columns, &col5)
    fmt.Println(columns)

    row := shim.Row{Columns: columns}
    fmt.Println(row)
    ok, err := stub.ReplaceRow(tableName, row)
    if err != nil {
      return nil, fmt.Errorf("updateRowRequests operation failed. %s", err)
    }
    if !ok {
      return nil, errors.New("updateRowRequests operation failed. Row with given key does not exists")
    }
  default :
    return nil, errors.New("updateRowRequests failed. No such table.")
  }
  return nil, nil
}

func GetTable(stub shim.ChaincodeStubInterface, tableName string) ([]byte, error) {
  var table *shim.Table
  table, err := stub.GetTable(tableName)
  if err != nil{
    return nil, fmt.Errorf("Something went wrong while getting table: %s",tableName)
  }
  fmt.Printf("Table %v",table)
  return nil,nil
}

func DeleteTable(stub shim.ChaincodeStubInterface, tableName string) error {
  err := stub.DeleteTable(tableName)
  if err != nil {
    return fmt.Errorf("Couldn't Delete %s, beacuse: %s \n",tableName,err)
  }
  return nil
}

// Delete row from a table
func DeleteRow(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]byte, error)  {
  patientId := args[0]
  userId := args[1]

  var columns []shim.Column
  col1 := shim.Column{Value: &shim.Column_String_{String_: patientId}}
  col2 := shim.Column{Value: &shim.Column_String_{String_: userId}}
  columns = append(columns, col1)
  columns = append(columns, col2)

  err := stub.DeleteRow(tableName, columns)
  if err != nil {
    return nil, fmt.Errorf("The operation delete row was not successful. Cause: %s\n", err)
  }

  return nil,nil
}
