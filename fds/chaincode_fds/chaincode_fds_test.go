package main

// go test chaincode_sla_test.go chaincode_sla.go

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ===========================================================
//  Utility 함수
// ===========================================================

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, key string, value string) {
	bytes := stub.State[key]
	if bytes == nil {
		fmt.Println("State", key, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Printf("State value [%v] : %v did not match the expected value %v.\n", key, string(bytes), value)
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, fnName string, args []string, value string) {
	bytes, err := stub.MockQuery(fnName, args)
	if err != nil {
		fmt.Println("Query", fnName, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", fnName, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Printf("Query value %v(%v) = %v did not match the expected value %v.\n", fnName, args, string(bytes), value)
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, fnName string, args []string) {
	_, err := stub.MockInvoke("1", fnName, args)
	if err != nil {
		fmt.Println("Invoke", fnName, "failed", err)
		t.FailNow()
	}
}

// ===========================================================
//   Test Init
// ===========================================================

func TestChaincodeFds_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	checkInit(t, stub, []string{})
	checkState(t, stub, FDS_NEXTEID_KEY, "1")
}

// ===========================================================
//   Test Query: FdsFraudEntry 조회 함수
// ===========================================================

/*
 * Get Fraud Entry with CID/MAC/UUID (cid: "cid", mac: "mac", uuid: "uuid")
 *
 *    Key            |   Value
 *  -----------------+---------------------------------------------------------------------------------------------------------------------------------
 *   [FDS_CID_cid]   |   "FDS_EID_1|FDS_EID_2"
 *   [FDS_MAC_mac]   |   "FDS_EID_1|FDS_EID_2"
 *   [FDS_UUID_uuid] |   "FDS_EID_1|FDS_EID_2"
 *   [FDS_EID_1]     |   {1 "cid" "mac" "uuid" "finaldate1" "finaltime1" "fdsproducedby1" "fdsregisteredby1" "fdsreason1" 9 "" ""} (in json string)
 *   [FDS_EID_2]     |   {2 "cid" "mac" "uuid" "finaldate2" "finaltime2" "fdsproducedby2" "fdsregisteredby2" "fdsreason2" 9 "" ""} (in json string)
 */
func TestChaincodeFds_Query_fdsGetFraudEntriesWith(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry1 := FdsFraudEntry{1, "cid", "mac", "uuid", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1", LS_BLACKLIST, "", ""}
	entry2 := FdsFraudEntry{2, "cid", "mac", "uuid", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{entry1, entry2}

	entry1InBytes, _ := json.Marshal(entry1)
	entry2InBytes, _ := json.Marshal(entry2)
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid", "mac", "uuid", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid", "mac", "uuid", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})

	checkQuery(t, stub, "fdsGetFraudEntriesWithCid", []string{"cid"}, string(entriesInBytes))
	checkQuery(t, stub, "fdsGetFraudEntriesWithMac", []string{"mac"}, string(entriesInBytes))
	checkQuery(t, stub, "fdsGetFraudEntriesWithUuid", []string{"uuid"}, string(entriesInBytes))

	checkState(t, stub, "FDS_CID_cid", "FDS_EID_1|FDS_EID_2")
	checkState(t, stub, "FDS_MAC_mac", "FDS_EID_1|FDS_EID_2")
	checkState(t, stub, "FDS_UUID_uuid", "FDS_EID_1|FDS_EID_2")
	checkState(t, stub, "FDS_EID_1", string(entry1InBytes))
	checkState(t, stub, "FDS_EID_2", string(entry2InBytes))
}

/*
 * Get All Fraud Entries
 *
 *    Key             |   Value
 *  ------------------+-----------------------------------------------------------------------------------------------------------------------------------
 *   [FDS_CID_cid1]   |   "FDS_EID_1"
 *   [FDS_CID_cid2]   |   "FDS_EID_2"
 *   [FDS_CID_cid3]   |   "FDS_EID_3"
 *   [FDS_MAC_mac1]   |   "FDS_EID_1"
 *   [FDS_MAC_mac2]   |   "FDS_EID_2"
 *   [FDS_MAC_mac3]   |   "FDS_EID_3"
 *   [FDS_UUID_uuid1] |   "FDS_EID_1"
 *   [FDS_UUID_uuid2] |   "FDS_EID_2"
 *   [FDS_UUID_uuid3] |   "FDS_EID_3"
 *   [FDS_EID_1]      |   {1 "cid1" "mac1" "uuid1" "finaldate1" "finaltime1" "fdsproducedby1" "fdsregisteredby1" "fdsreason1" 9 "" ""} (in json string)
 *   [FDS_EID_2]      |   {2 "cid2" "mac2" "uuid2" "finaldate2" "finaltime2" "fdsproducedby2" "fdsregisteredby2" "fdsreason2" 9 "" ""} (in json string)
 *   [FDS_EID_3]      |   {3 "cid3" "mac3" "uuid3" "finaldate3" "finaltime3" "fdsproducedby3" "fdsregisteredby3" "fdsreason3" 9 "" ""} (in json string)
 */
func TestChaincodeFds_Query_fdsGetAllFraudEntries(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry1 := FdsFraudEntry{1, "cid1", "mac1", "uuid1", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1", LS_BLACKLIST, "", ""}
	entry2 := FdsFraudEntry{2, "cid2", "mac2", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2", LS_BLACKLIST, "", ""}
	entry3 := FdsFraudEntry{3, "cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{entry1, entry2, entry3}

	entry1InBytes, _ := json.Marshal(entry1)
	entry2InBytes, _ := json.Marshal(entry2)
	entry3InBytes, _ := json.Marshal(entry3)
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid1", "mac1", "uuid1", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid2", "mac2", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3"})

	checkQuery(t, stub, "fdsGetAllFraudEntries", []string{}, string(entriesInBytes))

	checkState(t, stub, "FDS_EID_1", string(entry1InBytes))
	checkState(t, stub, "FDS_EID_2", string(entry2InBytes))
	checkState(t, stub, "FDS_EID_3", string(entry3InBytes))
}

// ===========================================================
//   Test Query: FdsFraudEntry 수정 함수
// ===========================================================

/*
 * Update Ledger Status of Fraud Entry ("BL" => "WL")
 *
 *   Fraud Entry              |   <Before>          |   <After>
 *  --------------------------+-------------------------------------------------------
 *   EID                      |   1                 |   1
 *   CID                      |   "cid"             |   "cid"
 *   MAC                      |   "mac"             |   "mac"
 *   FinalDate                |   "finaldate"       |   "finaldate"
 *   FinalTime                |   "finaltime"       |   "finaltime"
 *   ProducedBy               |   "fdsproducedby"   |   "fdsproducedby"
 *   RegisteredBy             |   "fdsregisteredby" |   "fdsregisteredby"
 *   Reason                   |   "fdsreason"       |   "fdsreason"
 *   LedgerStatus             |   9                 |   1
 *   LedgerStatusUpdateTime   |   ""                |   "ledgerstatusupdatetime"
 *   LedgerStatusUpdateReason |   ""                |   "ledgerstatusupdatereason"
 */
func TestChaincodeFds_Invoke_fdsUpdateLedgerStatusWithEid(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry := FdsFraudEntry{1, "cid", "mac", "uuid", "finaldate", "finaltime", "fdsproducedby", "fdsregisteredby", "fdsreason", LS_WHITELIST, "ledgerstatusupdatetime", "ledgerstatusupdatereason"}
	entryInBytes, _ := json.Marshal(entry)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid", "mac", "uuid", "finaldate", "finaltime", "fdsproducedby", "fdsregisteredby", "fdsreason"})
	checkInvoke(t, stub, "fdsUpdateLedgerStatusWithEid", []string{"1", "WL", "ledgerstatusupdatetime", "ledgerstatusupdatereason"})

	checkState(t, stub, "FDS_EID_1", string(entryInBytes))
}

// ===========================================================
//  Test Invoke: FdsFraudEntry 삭제 함수
// ===========================================================

/*
 * Delete Fraud Entries with EID (eid: 1)
 *
 *     |   <Before>               |   <After>
 *  ---+--------------------------+--------------------------
 *     |   Fraud Entry 1 (eid: 1) |   [deleted]
 *     |   Fraud Entry 2 (eid: 2) |   Fraud Entry 2
 */
func TestChaincodeFds_Invoke_fdsDeleteFraudEntryWithEid(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry2 := FdsFraudEntry{2, "cid2", "mac2", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{FdsFraudEntry{}, entry2}
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid1", "mac1", "uuid1", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid2", "mac2", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})
	checkInvoke(t, stub, "fdsDeleteFraudEntryWithEid", []string{"1"})

	checkQuery(t, stub, "fdsGetAllFraudEntries", []string{}, string(entriesInBytes))
}

/*
 * Delete Fraud Entries with CID (cid: "cid")
 *
 *     |   <Before>                    |   <After>
 *  ---+-------------------------------+---------------------
 *     |   Fraud Entry 1 (cid: "cid")  |   [deleted]
 *     |   Fraud Entry 2 (cid: "cid")  |   [deleted]
 *     |   Fraud Entry 3 (cid: "cid3") |   Fraud Entry 3
 */
func TestChaincodeFds_Invoke_DeleteFruadEntryWithCid(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry3 := FdsFraudEntry{3, "cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{FdsFraudEntry{}, FdsFraudEntry{}, entry3}
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid", "mac1", "uuid1", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid", "mac2", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3"})
	checkInvoke(t, stub, "fdsDeleteFraudEntryWithCid", []string{"cid"})

	checkQuery(t, stub, "fdsGetAllFraudEntries", []string{}, string(entriesInBytes))
}

/*
 * Delete Fraud Entries with MAC (mac: "mac")
 *
 *     |   <Before>                    |   <After>
 *  ---+-------------------------------+---------------------
 *     |   Fraud Entry 1 (mac: "mac")  |   [deleted]
 *     |   Fraud Entry 2 (mac: "mac")  |   [deleted]
 *     |   Fraud Entry 3 (mac: "mac3") |   Fraud Entry 3
 */
func TestChaincodeFds_Invoke_DeleteFruadEntryWithMac(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry3 := FdsFraudEntry{3, "cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{FdsFraudEntry{}, FdsFraudEntry{}, entry3}
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid1", "mac", "uuid1", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid2", "mac", "uuid2", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3"})
	checkInvoke(t, stub, "fdsDeleteFraudEntryWithMac", []string{"mac"})

	checkQuery(t, stub, "fdsGetAllFraudEntries", []string{}, string(entriesInBytes))
}

/*
 * Delete Fraud Entries with UUID (uuid: "uuid")
 *
 *     |   <Before>                      |   <After>
 *  ---+---------------------------------+---------------------
 *     |   Fraud Entry 1 (uuid: "uuid")  |   [deleted]
 *     |   Fraud Entry 2 (uuid: "uuid")  |   [deleted]
 *     |   Fraud Entry 3 (uuid: "uuid3") |   Fraud Entry 3
 */
func TestChaincodeFds_Invoke_DeleteFruadEntryWithUuid(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("fds_chaincode", scc)

	entry3 := FdsFraudEntry{3, "cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3", LS_BLACKLIST, "", ""}
	entries := []FdsFraudEntry{FdsFraudEntry{}, FdsFraudEntry{}, entry3}
	entriesInBytes, _ := json.Marshal(entries)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid1", "mac1", "uuid", "finaldate1", "finaltime1", "fdsproducedby1", "fdsregisteredby1", "fdsreason1"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid2", "mac2", "uuid", "finaldate2", "finaltime2", "fdsproducedby2", "fdsregisteredby2", "fdsreason2"})
	checkInvoke(t, stub, "fdsCreateFraudEntry", []string{"cid3", "mac3", "uuid3", "finaldate3", "finaltime3", "fdsproducedby3", "fdsregisteredby3", "fdsreason3"})
	checkInvoke(t, stub, "fdsDeleteFraudEntryWithUuid", []string{"uuid"})

	checkQuery(t, stub, "fdsGetAllFraudEntries", []string{}, string(entriesInBytes))
}
