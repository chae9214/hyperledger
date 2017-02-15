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
}

type FdsFraudEntry struct {
	Cid          string `json:"cid"`
	Mac          string `json:"mac"`
	Uuid         string `json:"uuid"`
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
const PREFIX_EID, PREFIX_EID_END = "FDS_EID_", "FDS_F"
const PREFIX_CID, PREFIX_CID_END = "FDS_CID_", "FDS_D"
const PREFIX_MAC, PREFIX_MAC_END = "FDS_MAC_", "FDS_N"
const PREFIX_UUID, PREFIX_UUID_END = "FDS_UUID_", "FDS_V"

const NEXTEIDKEY = "FDS_NEXTEID"

const FIELDSEP = "|"
const ENTRYSEP = "$"

// ===========================================================
//  Helper 함수
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

func printFraudEntries(entries []FdsFraudEntry) {
	fmt.Println("[")
	for _, entry := range entries {
		//fmt.Printf("\t[%v%v] = %v\n", PREFIX_EID, i+1, entry)
		fmt.Printf("\t%v\n", entry)
	}
	fmt.Println("]")
}

// ===========================================================
//  SimpleChaincode 함수
// ===========================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Initializing requires 0 argument but given" + strconv.Itoa(len(args)))
	}
	t.fdsSetNextEid(stub, 1)
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "fdsCreateFraudEntry":
		return t.fdsCreateFraudEntry(stub, args)
	case "fdsDeleteFraudEntryWithCid":
		return t.fdsDeleteFraudEntryWithCid(stub, args)
	case "fdsDeleteFraudEntryWithMac":
		return t.fdsDeleteFraudEntryWithMac(stub, args)
	case "fdsDeleteFraudEntryWithUuid":
		return t.fdsDeleteFraudEntryWithUuid(stub, args)
	}
	return nil, errors.New("Invalid invoke function name")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "fdsGetAllFraudEntries":
		return t.fdsGetAllFraudEntries(stub, args)
	case "fdsGetFraudEntriesWithCid":
		return t.fdsGetFraudEntriesWithCid(stub, args)
	case "fdsGetFraudEntriesWithMac":
		return t.fdsGetFraudEntriesWithMac(stub, args)
	case "fdsGetFraudEntriesWithUuid":
		return t.fdsGetFraudEntriesWithUuid(stub, args)
	case "test":
		return t.testFunc(stub, args)
	}
	return nil, errors.New("Invalid query function name")
}

func main() {
	t := new(SimpleChaincode)
	err := shim.Start(t)
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ===========================================================
//   EID 조회/수정 함수
// ===========================================================

func (t *SimpleChaincode) fdsGetNextEid(stub shim.ChaincodeStubInterface) int {
	nextEidInBytes, _ := stub.GetState(NEXTEIDKEY)
	nextEidInInt, _ := strconv.Atoi(string(nextEidInBytes))
	return nextEidInInt
}

func (t *SimpleChaincode) fdsSetNextEid(stub shim.ChaincodeStubInterface, nextEidInInt int) {
	nextEidInBytes := []byte(strconv.Itoa(nextEidInInt))
	stub.PutState(NEXTEIDKEY, nextEidInBytes)
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

	nextEid := t.fdsGetNextEid(stub)
	eidKey := PREFIX_EID + strconv.Itoa(nextEid)
	cidKey := PREFIX_CID + args[IND_CID]
	macKey := PREFIX_MAC + args[IND_MAC]
	uuidKey := PREFIX_UUID + args[IND_UUID]

	entry := FdsFraudEntry{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7]}
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

	t.fdsSetNextEid(stub, nextEid+1)
	return nil, nil
}

// ===========================================================
//  FDS 삭제 함수
// ===========================================================

func (t *SimpleChaincode) fdsDeleteFraudEntryWithCid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) fdsDeleteFraudEntryWithMac(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) fdsDeleteFraudEntryWithUuid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
//   FDS 조회 함수
// ===========================================================

func (t *SimpleChaincode) fdsGetFraudEntriesWithCid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	entries := make([]FdsFraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		if entryInBytes == nil {
			continue
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FdsFraudEntry
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

func (t *SimpleChaincode) fdsGetFraudEntriesWithMac(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	entries := make([]FdsFraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		if entryInBytes == nil {
			continue
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FdsFraudEntry
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

func (t *SimpleChaincode) fdsGetFraudEntriesWithUuid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	entries := make([]FdsFraudEntry, len(eidKeys))
	for i, eidKey := range eidKeys {
		entryInBytes, err = stub.GetState(eidKey)
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		if entryInBytes == nil {
			continue
		}
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))

		var entry FdsFraudEntry
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

func (t *SimpleChaincode) fdsGetAllFraudEntries(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var entryInBytes []byte
	var entriesInBytes []byte
	var err error

	if len(args) != 0 {
		return nil, errors.New("Looking up all entries requires 0 argument but given" + strconv.Itoa(len(args)))
	}

	nextEid := t.fdsGetNextEid(stub)
	entries := make([]FdsFraudEntry, nextEid-1)
	for i := 0; i < nextEid-1; i++ {
		eidKey := PREFIX_EID + strconv.Itoa(i+1)

		entryInBytes, err = stub.GetState(eidKey)
		//fmt.Println("ENTRY looked up with", eidKey, ":", string(entryInBytes))
		if err != nil {
			return nil, errors.New("Failed to delete state for" + eidKey)
		}
		if entryInBytes == nil {
			continue
		}

		var entry FdsFraudEntry
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

// ===========================================================
//   테스트 함수
// ===========================================================

func (t *SimpleChaincode) testFunc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	iter, _ := stub.RangeQueryState(args[0], args[1])
	fmt.Println("START OF ITERATION")
	for iter.HasNext() {
		key, value, _ := iter.Next()
		fmt.Println("\t" + key + "\t:\t" + string(value))
	}
	fmt.Println("END OF ITERATION")
	return nil, nil
}
