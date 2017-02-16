package main

// go test chaincode_sla_test.go chaincode_sla.go

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Printf("State value %v [%v] was not the expected value[%v]\n", name, string(bytes), value)
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
		fmt.Printf("State value %v [%v] was not the expected value[%v]\n", fnName, string(bytes), value)
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

func TestChaincodeSla_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	checkInit(t, stub, []string{})

	checkState(t, stub, SLA_CONTRACT_ID_COUNT_KEY, "1")
	checkState(t, stub, SLA_EVALUATION_ID_COUNT_KEY, "1")
	checkState(t, stub, CURRENT_YEAR_KEY, "2017")
}

func TestChaincodeSla_Query_slaGetContractId(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	checkInit(t, stub, []string{})
	checkQuery(t, stub, "slaGetContractId", []string{}, "SLA_CONT_2017_00001")
}

func TestChaincodeSla_Invoke_slaCreateContract(t *testing.T) {
	contractContentInJson := `{
		  "RegId": "SLA_CONT_2017_00005",
		  "Name": "홍길동",
		  "Kind": "보통계약",
		  "StaDate": "2017-02-01",
		  "EndDate": "2017-12-01",
		  "Client": "신한은행",
		  "ClientPerson": "개인",
		  "ClientPersonTel": "010-1111-2222",
		  "AssessDate": "2017-12-31",
		  "Progression": "작성",
		  "AssessYn": "포함",
		  "Approvals": [
		    {
		      "ApprovalUserId": "test",
		      "ApprovalCompany": "test",
		      "ApprovalDepartment": "test",
		      "ApprovalName": "test",
		      "ApprovalState": "test",
		      "ApprovalDate": "test",
		      "ApprovalComment": "test",
		      "ApprovalAlram": "test"
		    },
		    {
		      "ApprovalUserId": "test2",
		      "ApprovalCompany": "test2",
		      "ApprovalDepartment": "test2",
		      "ApprovalName": "test2",
		      "ApprovalState": "test2",
		      "ApprovalDate": "test2",
		      "ApprovalComment": "test2",
		      "ApprovalAlram": "test2"
		    }
		  ],
		  "ServiceItems": [
		    {
		      "ServiceItem": "test",
		      "ScoreItem": "test",
		      "MeasurementItem": "test",
		      "ExplainItem": "test",
		      "DivideScore": "test"
		    }
		  ]
		}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "slaCreateContract", []string{contractContentInJson})

	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00005"}, contractContentInJson)
}
