package main

import (
	"encoding/json"
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

type FraudEntry struct {
	CID          string `json:"cid"`
	MAC          string `json:"mac"`
	UUID         string `json:"uuid"`
	FinalDate    string `json:"finalDate"`
	FinalTime    string `json:"finalTime"`
	ProducedBy   string `json:"producedBy"`
	RegisteredBy string `json:"registeredBy"`
	Reason       string `json:"reason"`
}

// fraud entry 의 필드갯수
const NUM_FIELDS = 8

// fraud entry 의 각 필드별 인덱스
const IND_CID = 0
const IND_MAC = 1
const IND_UUID = 2
const IND_FINALDATE = 3
const IND_FINALTIME = 4
const IND_PRODUCEDBY = 5
const IND_REGISTEREDBY = 6
const IND_REASON = 7

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
	fieldnames[IND_PRODUCEDBY] = "producedBy"
	fieldnames[IND_REGISTEREDBY] = "registeredBy"
	fieldnames[IND_REASON] = "reason"

	jsonStr := "[\n"
	for i, entry := range entries {
		jsonStr += "\t{ "
		if entry != "" {
			fields := strings.Split(entry, FIELDSEP)
			for j, fieldname := range fieldnames {
				jsonStr += "\"" + fieldname + "\" : \"" + fields[j] + "\""
				if j < len(fields)-1 {
					jsonStr += ", "
				}
			}
		}
		jsonStr += " }"
		if i < len(entries)-1 {
			jsonStr += ",\n"
		}
	}
	jsonStr += "\n]"

	return jsonStr
}

func printFraudEntries(entries []FraudEntry) {
	fmt.Println("[\n")
	for _, entry := range entries {
		//fmt.Printf("\t[%v%v] = %v\n", PREFIX_EID, i+1, entry)
		fmt.Printf("\t%v", entry)
	}
	fmt.Println("\n]")
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
		return t.fdsCreateFraudEntry(stub, args)
	case "fdsDeleteWithCid":
		return t.fdsDeleteWithCid(stub, args)
	case "fdsDeleteWithMac":
		return t.fdsDeleteWithMac(stub, args)
	case "fdsDeleteWithUuid":
		return t.fdsDeleteWithUuid(stub, args)
	}
	return nil, errors.New("Invalid invoke function name. Expecting \"register\" \"removewith~\"")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "fdsGetAll":
		return t.fdsGetAll(stub, args)
	case "fdsGetWithCid":
		return t.fdsGetWithCid(stub, args)
	case "fdsGetWithMac":
		return t.fdsGetWithMac(stub, args)
	case "fdsGetWithUuid":
		return t.fdsGetWithUuid(stub, args)
	case "getnexteid":
		return []byte(strconv.Itoa(t.nextEID)), nil
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

func (t *SimpleChaincode) fdsCreateFraudEntry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entryInBytes []byte
	var eidsInBytes []byte
	var err error

	if len(args) != NUM_FIELDS {
		return nil, errors.New("Register requires" + strconv.Itoa(NUM_FIELDS) + "arguements but given" + strconv.Itoa(len(args)))
	}

	eidKey := PREFIX_EID + strconv.Itoa(t.nextEID)
	cidKey := PREFIX_CID + args[IND_CID]
	macKey := PREFIX_MAC + args[IND_MAC]
	uuidKey := PREFIX_UUID + args[IND_UUID]

	entry := FraudEntry{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7]}
	entryInBytes, err = json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(eidKey, entryInBytes)
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

func (t *SimpleChaincode) fdsDeleteWithCid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) fdsDeleteWithMac(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) fdsDeleteWithUuid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
//   조회 함수
// ===========================================================

//func (t *SimpleChaincode) fdsGetWithCid(cid string) (entries [][]string, result bool) {
func (t *SimpleChaincode) fdsGetWithCid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var entriesInBytes []byte
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
	//fmt.Println("EIDS looked up with", cidKey, ":", string(eidsInBytes))

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]FraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FraudEntry
		err = json.Unmarshal(entryInBytes, &entry)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	entriesInBytes, err = json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	fmt.Println("Query response:")
	printFraudEntries(entries)
	return entriesInBytes, nil
}

//func (t *SimpleChaincode) fdsGetWithMac(mac string) (entries [][]string, result bool) {
func (t *SimpleChaincode) fdsGetWithMac(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var entriesInBytes []byte
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
	//fmt.Println("EIDS looked up with", macKey, ":", string(eidsInBytes))

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]FraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FraudEntry
		err = json.Unmarshal(entryInBytes, &entry)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	entriesInBytes, err = json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	fmt.Println("Query response:")
	printFraudEntries(entries)
	return entriesInBytes, nil
}

//func (t *SimpleChaincode) fdsGetWithUuid(uuid string) (entries [][]string, result bool) {
func (t *SimpleChaincode) fdsGetWithUuid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var eidsInBytes []byte
	var entryInBytes []byte
	var entriesInBytes []byte
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
	//fmt.Println("EIDS looked up with", uuidKey, ":", string(eidsInBytes))

	eidKeys := byteArrayToStringArray(eidsInBytes)
	entries := make([]FraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FraudEntry
		err = json.Unmarshal(entryInBytes, &entry)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	entriesInBytes, err = json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	fmt.Println("Query response:")
	printFraudEntries(entries)
	return entriesInBytes, nil
}

func (t *SimpleChaincode) fdsGetAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entryInBytes []byte
	var entriesInBytes []byte
	var err error

	if len(args) != 0 {
		return nil, errors.New("Looking up all entries requires 0 argument but given" + strconv.Itoa(len(args)))
	}

	entries := make([]FraudEntry, t.nextEID-1)
	for i := 0; i < t.nextEID-1; i++ {
		eidKey := PREFIX_EID + strconv.Itoa(i+1)

		entryInBytes, err = stub.GetState(eidKey)
		fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))
		if err != nil {
			return nil, err
		}
		if entryInBytes == nil {
			continue
		}

		var entry FraudEntry
		err = json.Unmarshal(entryInBytes, &entry)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	entriesInBytes, err = json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	fmt.Println("Query response:")
	printFraudEntries(entries)
	return entriesInBytes, nil
}
