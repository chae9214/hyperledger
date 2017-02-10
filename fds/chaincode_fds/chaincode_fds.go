package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ===========================================================
//  Struct 및 Constant 정의
// ===========================================================

type SimpleChaincode struct {
	nextEID int
}

// fraud entry 의 필드갯수
const NUM_FIELDS = 8

// fraud entry 의 각 필드별 인덱스
const IND_CID = 0
const IND_MAC = 1
const IND_UUID = 2
const IND_FINALDATE = 3
const IND_FINALTIME = 4
const IND_FDSPRODUCEDBY = 5
const IND_FDSREGISTEREDBY = 6
const IND_FDSREASON = 7

// fraud entry 의 각 필드별 json 명칭
// var fieldnames = make([]string, NUM_FIELDS)
// fieldnames[IND_CID] = "cid"
// fieldnames[IND_MAC] = "mac"
// fieldnames[IND_UUID] = "uuid"
// fieldnames[IND_FINALDATE] = "finalDate"
// fieldnames[IND_FINALTIME] = "finalTime"
// fieldnames[IND_FDSPRODUCEDBY] = "producedBy"
// fieldnames[IND_FDSREGISTEREDBY] = "registeredBy"
// fieldnames[IND_FDSREASON] = "reason"

// key-value store 의 키 구분자
const PREFIX_EID = "eid_"
const PREFIX_CID = "cid_"
const PREFIX_MAC = "mac_"
const PREFIX_UUID = "uuid_"

const FIELDSEP = "|"
const ENTRYSEP = "$"

/*
// ===========================================================
//  Initialization 함수
// ===========================================================

func CreateStub() ChaincodeStubInterface {
	var stub ChaincodeStubInterface
	stub.kvs = make(map[string][]byte)
	return stub
}

func CreateFDSChaincodeStub() FDSChaincodeStub {
	kvs := make(map[string][]byte)
	nextEID := 1
	stub := FDSChaincodeStub{ChaincodeStubInterface{kvs}, nextEID}
	return stub
}

func CreateSLAChaincodeStub() SLAChaincodeStub {
	kvs := make(map[string][]byte)
	stub := SLAChaincodeStub{ChaincodeStubInterface{kvs}}
	return stub
}

func (stub *ChaincodeStubInterface) String() string {
	var s string

	s += "<KVS>\n"
	for k, v := range stub.kvs {
		s += "\t" + k + ": " + string(v) + "\n"
	}
	return s
}
*/

// ===========================================================
//  Serialize / Deserialize 함수
// ===========================================================

func stringArrayToByteArray(slist []string) []byte {
	return []byte(strings.Join(slist, FIELDSEP))
}

func byteArrayToStringArray(b []byte) []string {
	if string(b) == "" {
		return []string{}
	} else {
		return strings.Split(string(b), FIELDSEP)
	}
}

func appendToEIDList(b []byte, eid string) []byte {
	eidKeys := byteArrayToStringArray(b)
	return stringArrayToByteArray(append(eidKeys, eid))
}

func entryStringsToJsonString(entries []string) string {
	var fieldnames = make([]string, NUM_FIELDS)
	fieldnames[IND_CID] = "cid"
	fieldnames[IND_MAC] = "mac"
	fieldnames[IND_UUID] = "uuid"
	fieldnames[IND_FINALDATE] = "finalDate"
	fieldnames[IND_FINALTIME] = "finalTime"
	fieldnames[IND_FDSPRODUCEDBY] = "producedBy"
	fieldnames[IND_FDSREGISTEREDBY] = "registeredBy"
	fieldnames[IND_FDSREASON] = "reason"

	jsonStr := "[\n"
	for i, entry := range entries {
		jsonStr += "\t{"
		fields := strings.Split(entry, FIELDSEP)
		for j, fieldname := range fieldnames {
			jsonStr += "\"" + fieldname + "\" : \"" + fields[j] + "\""
			if j < len(fields)-1 {
				jsonStr += ", "
			}
		}
		jsonStr += "}"
		if i < len(entries)-1 {
			jsonStr += ",\n"
		}
	}
	jsonStr += "\n]"

	return jsonStr
}

/*
// ===========================================================
//  ChaincodeStubInterface 함수
// ===========================================================

func (stub *ChaincodeStubInterface) PutState(key string, value []byte) error {
	if value == nil {
		return errors.New("entry cannot be empty")
	} else {
		stub.kvs[key] = value
		return nil
	}
}

func (stub *ChaincodeStubInterface) GetState(key string) ([]byte, error) {
	value := stub.kvs[key]
	return value, nil
}

func (stub *ChaincodeStubInterface) DelState(key string) error {
	delete(stub.kvs, key)
	return nil
}

func (stub *ChaincodeStubInterface) GetKVSLength() int {
	return len(stub.kvs)
}
*/

// ===========================================================
//  SimpleChaincode 함수
// ===========================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Initializing requires 0 argument but given" + strconv.Itoa(len(args)))
	}
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "register":
		return t.RegisterFraudEntry(stub, args)
	case "removewithcid":
		return t.RemoveWithCID(stub, args)
	case "removewithmac":
		return t.RemoveWithMAC(stub, args)
	case "removewithuuid":
		return t.RemoveWithUUID(stub, args)
	}
	return nil, errors.New("Invalid invoke function name. Expecting \"register\" \"removewith~\"")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "lookupall":
		return t.LookupAll(stub, args)
	case "lookupwithcid":
		return t.LookupWithCID(stub, args)
	case "lookupwithmac":
		return t.LookupWithMAC(stub, args)
	case "lookupwithuuid":
		return t.LookupWithUUID(stub, args)
	}
	return nil, errors.New("Invalid query function name. Expecting \"lookupwith~\"")
}

func main() {
	t := new(SimpleChaincode)
	t.nextEID = 1
	err := shim.Start(t)
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ===========================================================
//  FDS 등록/수정 함수
// ===========================================================

func (t *SimpleChaincode) RegisterFraudEntry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var err error

	if len(args) != NUM_FIELDS {
		return nil, errors.New("Register requires" + strconv.Itoa(NUM_FIELDS) + "arguements but given" + strconv.Itoa(len(args)))
	}

	eidKey := PREFIX_EID + strconv.Itoa(t.nextEID)
	cidKey := PREFIX_CID + args[IND_CID]
	macKey := PREFIX_MAC + args[IND_MAC]
	uuidKey := PREFIX_UUID + args[IND_UUID]

	err = stub.PutState(eidKey, stringArrayToByteArray(args))
	if err != nil {
		return nil, err
	}

	eidsInBytes, err = stub.GetState(cidKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + cidKey)
	}
	err = stub.PutState(cidKey, appendToEIDList(eidsInBytes, eidKey))
	if err != nil {
		return nil, err
	}

	eidsInBytes, err = stub.GetState(macKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + macKey)
	}
	err = stub.PutState(macKey, appendToEIDList(eidsInBytes, eidKey))
	if err != nil {
		return nil, err
	}

	eidsInBytes, err = stub.GetState(uuidKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + uuidKey)
	}
	err = stub.PutState(uuidKey, appendToEIDList(eidsInBytes, eidKey))
	if err != nil {
		return nil, err
	}

	t.nextEID++
	return nil, nil
}

// ===========================================================
//  FDS 삭제 함수
// ===========================================================

func (t *SimpleChaincode) RemoveWithCID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Removing with CID requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	cidKey := PREFIX_CID + args[0]

	eidsInBytes, err = stub.GetState(cidKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + cidKey)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	for _, eidKey := range eidKeys {
		err = stub.DelState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
	}
	err = stub.DelState(cidKey)
	if err != nil {
		return nil, errors.New("Failed to delete state for" + cidKey)
	}
	return nil, nil
}

func (t *SimpleChaincode) RemoveWithMAC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Removing with MAC requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	macKey := PREFIX_MAC + args[0]

	eidsInBytes, err = stub.GetState(macKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + macKey)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	for _, eidKey := range eidKeys {
		err = stub.DelState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
	}
	err = stub.DelState(macKey)
	if err != nil {
		return nil, errors.New("Failed to delete state for" + macKey)
	}
	return nil, nil
}

func (t *SimpleChaincode) RemoveWithUUID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Removing with UUID requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	uuidKey := PREFIX_UUID + args[0]

	eidsInBytes, err = stub.GetState(uuidKey)
	if err != nil {
		return nil, errors.New("Failed to get state for" + uuidKey)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	for _, eidKey := range eidKeys {
		err = stub.DelState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
	}
	err = stub.DelState(uuidKey)
	if err != nil {
		return nil, errors.New("Failed to delete state for" + uuidKey)
	}
	return nil, nil
}

// ===========================================================
//  FDS 조회 함수
// ===========================================================

//func (t *SimpleChaincode) LookupWithCID(cid string) (entries [][]string, result bool) {
func (t *SimpleChaincode) LookupWithCID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Looking up with CID requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	cidKey := PREFIX_CID + args[0]

	eidsInBytes, err = stub.GetState(cidKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + cidKey + "\"}"
		return nil, errors.New(jsonResp)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]string, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		entries[i] = string(entryInBytes)
	}

	jsonResp := entryStringsToJsonString(entries)
	fmt.Println("Query response:%s\n", jsonResp)
	return []byte(strings.Join(entries, ENTRYSEP)), nil
}

//func (t *SimpleChaincode) LookupWithMAC(mac string) (entries [][]string, result bool) {
func (t *SimpleChaincode) LookupWithMAC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Looking up with MAC requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	macKey := PREFIX_MAC + args[0]

	eidsInBytes, err = stub.GetState(macKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + macKey + "\"}"
		return nil, errors.New(jsonResp)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]string, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		entries[i] = string(entryInBytes)
	}

	jsonResp := entryStringsToJsonString(entries)
	fmt.Println("Query response:%s\n", jsonResp)
	return []byte(strings.Join(entries, ENTRYSEP)), nil
}

//func (t *SimpleChaincode) LookupWithUUID(uuid string) (entries [][]string, result bool) {
func (t *SimpleChaincode) LookupWithUUID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Looking up with UUID requires 1 argument but given" + strconv.Itoa(len(args)))
	}

	uuidKey := PREFIX_UUID + args[0]

	eidsInBytes, err = stub.GetState(uuidKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + uuidKey + "\"}"
		return nil, errors.New(jsonResp)
	}

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]string, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		entries[i] = string(entryInBytes)
	}

	jsonResp := entryStringsToJsonString(entries)
	fmt.Println("Query response:%s\n", jsonResp)
	return []byte(strings.Join(entries, ENTRYSEP)), nil
}

func (t *SimpleChaincode) LookupAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entryInBytes []byte
	var err error

	if len(args) != 0 {
		return nil, errors.New("Looking up all entries requires 0 argument but given" + strconv.Itoa(len(args)))
	}

	entries := make([]string, t.nextEID-1)
	for i := 0; i < t.nextEID-1; i++ {
		eidKey := PREFIX_EID + strconv.Itoa(i+1)

		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, err
		}
		if entryInBytes == nil {
			continue
		}
		entries[i] = string(entryInBytes)
	}

	jsonResp := entryStringsToJsonString(entries)
	fmt.Println("Query response:\n", jsonResp)
	return []byte(strings.Join(entries, ENTRYSEP)), nil
}
