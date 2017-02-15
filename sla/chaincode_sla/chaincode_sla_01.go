/*copyright IBM Corp. 2016 All Rights Reserved.

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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"strings"
)

// ===========================================================
//  Struct �� Constant ����
// ===========================================================
type SimpleChaincode struct {
}

// Sla Contract ����ü�� �����մϴ�.
type SlaContract struct {
	RegId           string           `json:  "RegId"`           // SLA����Ϲ�ȣ
	Name            string           `json:  "Name"`            // SLA����
	Kind            string           `json:  "Kind"`            // SLA�������
	StaDate         string           `json:  "StaDate"`         // SLA��������
	EndDate         string           `json:  "EndDate"`         // SLA���������
	Client          string           `json:  "Client"`          // �����
	ClientPerson    string           `json:  "ClientPerson"`    // ������ڸ�
	ClientPersonTel string           `json:  "ClientPersonTel"` // ���������ȭ��ȣ
	AssessDate      string           `json:  "AssessDate"`      // �򰡿�����
	Progression     string           `json:  "Progression"`     // ����ܰ�
	AssessYn        string           `json:  "AssessYn"`        // SLA�� ��󿩺�
	Approvals       []SlaApproval    `json:  "Approvals"`       // SLA���缱����
	ServiceItems    []SlaServiceItem `json:  "ServiceItems"`    // SLA���׸�
}

// Sla Approval ����ü�� �����մϴ�.
type SlaApproval struct {
	ApprovalUserId     string `json:  "ApprovalUserId"`     // ��������ID
	ApprovalCompany    string `json:  "ApprovalCompany"`    // ����ȸ���
	ApprovalDepartment string `json:  "ApprovalDepartment"` // ����μ���
	ApprovalName       string `json:  "ApprovalName"`       // �����ڸ�
	ApprovalState      string `json:  "ApprovalState"`      // �������
	ApprovalDate       string `json:  "ApprovalDate"`       // ��������
	ApprovalComment    string `json:  "ApprovalComment"`    // �ǰ߳���
	ApprovalAlram      string `json:  "ApprovalAlram"`      // �˶�����  TODO Alram --> Alarm
}

// Sla ServiceItem ����ü�� �����մϴ�.
type SlaServiceItem struct {
	ServiceItem     string `json:  "ServiceItem"`     // �����׸�
	ScoreItem       string `json:  "ScoreItem"`       // ���׸�
	MeasurementItem string `json:  "MeasurementItem"` // ��������
	ExplainItem     string `json:  "ExplainItem"`     // ����
	DivideScore     string `json:  "DivideScore"`     // SLA�������
}

// Sla EvaluationRoot����ü�� �����մϴ�.
type SlaEvalutionRoot struct {
	RegId        string           `json:  "RegId"`        // SLA����Ϲ�ȣ
	ContractId   string           `json:  "ContractId"`   // SLA����
	Status       string           `json:  "Status"`       // SLA����
	Evaluations  []SlaEvaluation  `json:  "Evaluations"`  // SLA�򰡵�Ϲ�ȣ
	ServiceItems []SlaServiceItem `json:  "ServiceItems"` // SLA���׸�
}

// Sla Evaluation ����ü�� �����մϴ�.
type SlaEvaluation struct {
	RegId                 string        `json:  "SlaContractRegId"` // SLA����Ϲ�ȣ
	EvaluationRootId      string        `json:  "SlaContractName"`  // SLA����
	ScoresForServiceItems string        `json:  "SlaContractName"`  // SLA�������׸�
	Approvals             []SlaApproval `json:  "Approvals"`        // SLA���缱����
}

// key-value store �� Ű ������
const FIELDSEP = "|"
const ENTRYSEP = ","
const SLA_ALL_DATA = "SLA_ALL_DATA"

// ===========================================================
//  Initialization �Լ�
// ===========================================================

// �ʱ�ȭ�� ó���մϴ�.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var A, B string
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ��� �̺�Ʈ�� ȣ���մϴ�.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	// ��� ���� ����
	case "slaCreateContract": // ��û�ڰ� ��� ���� (���� ���� �� �ӽ�����)
		return t.slaCreateContract(stub, args)

	case "slaUpdateContract": // ��û�ڰ� ��� ���� (�ӽ����� �� / ���ΰ��� ��)
		return t.slaUpdateContract(stub, args)

	case "slaAbandonContract": // ��û�ڰ� ��� ��� (�ӽ����� �� / ���ΰ��� ��)
		return t.slaUpdateContract(stub, args)

	// ���� ��û / ���� / ����
	case "slaSubmitContract": // ��û�� --> ������  (���� ���� �� ���� ��û / ����  �� ���� ��û )
		return t.slaSubmitContract(stub, args)

	case "slaApproveContract": // ������ --> ���� ������
		return t.slaApproveContract(stub, args)

	case "slaRejectContract": // ������ --> ��û��
		return t.slaRejectContract(stub, args)

	// ���� ����
	case "slaCloseContract": // ���� ������ ����
		return t.slaCloseContract(stub, args)

	// ��ü �� ����
	case "slaCreateEvaluationTemplateFromContract": // ���� �� ���� (����� ���� ���� ��),
		_, err1 := t.slaCreateEvaluationRootFromContract(stub, args)
		_, err2 := t.slaCreateEvaluationsFromContract(stub, args)

		if err1 == nil {
			err1 = err2
		}
		return nil, err1

	// ���� �� ����
	case "slaInitEvaluationValues": // ���� ���� ������ �Է�
		return t.slaInitEvaluationValues(stub, args)

	case "slaUpdateEvaluationValues": // ������ ����
		return t.slaUpdateEvaluationValues(stub, args)

	// ���� ��û / ���� / ����
	case "slaSubmitEvaluation": // ��û�� --> ������  (���� ���� �� ���� ��û / ����  �� ���� ��û )
		return t.slaSubmitEvaluation(stub, args)

	case "slaApproveEvaluation": // ������ --> ���� ������
		return t.slaApproveEvaluation(stub, args)

	case "slaRejectEvaluation": // ������ --> ��û��
		return t.slaRejectEvaluation(stub, args)

	// ���� ��û / ���� / ����
	case "slaSubmitPayment": // ��û�� --> ������
		return t.slaSubmitPayment(stub, args)

	case "slaClosePayment": // ������
		return t.slaClosePayment(stub, args)

	// ���� �� ������
	case "slaCloseEvaluation":
		return t.slaCloseEvaluation(stub, args)

	// ��ü �� ������
	case "slaCloseEvaluationRoot": // ������ ���� �򰡰� �������� ���, �ڵ� ȣ��
		return t.slaCloseEvaluationRoot(stub, args)

	}
	return nil, errors.New("Invalid invoke function name. Expecting \"slaCreateContract\" \"slaUpdateContract\" \"slaApproveContract\" \"slaRejectContract\"")
}

// ���� �̺�Ʈ�� ó���մϴ�.
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "slaGetAllContracts":
		return t.slaGetAllContracts(stub, args)

	case "slaGetContractWithId":
		return t.slaGetContractWithId(stub, args)

	case "slaGetContractsWithName":
		return t.slaGetContractsWithName(stub, args)

	case "slaGetContractsWithClient":
		return t.slaGetContractsWithClient(stub, args)

	case "slaGetAllEvaluations":
		return t.slaGetAllContracts(stub, args)

	case "slaGetEvaluationWithId":
		return t.slaGetContractWithId(stub, args)

	case "slaGetEvaluationsWithName":
		return t.slaGetContractsWithName(stub, args)

	case "slaGetEvaluationsWithClient":
		return t.slaGetContractsWithClient(stub, args)

	}
	return nil, errors.New("Invalid Query function name. Expecting \"slaGetAllContracts\" \"slaGetContractWithId\" \"slaGetContractsWithName\" \"slaGetContractsWithClient\"")
}

// �����Լ��� ó���մϴ�.
func main() {





	// ���ü�� �̺�Ʈ�� ȣ���մϴ�.
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Er4ror starting Simple chaincode: %s", err)
	}
}

// ===========================================================
//  SLAChaincodeStub ��� �Լ�
// ===========================================================

// ����� ����մϴ�.
func (t *SimpleChaincode) slaCreateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var data SlaContract

	content := args[0]

	fmt.Printf("slaCreateContract Input Args:%s\n", args[0])

	// JSON �����͸� ���ڵ�(Unmarshal)�մϴ�.
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	}

	// JSON �����͸� �����Ͽ� ���ڵ�(Unmarshal)�մϴ�.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	// }

	// fmt.Printf("= jsonData ====================================================\n")
	// fmt.Println(data.RegId)
	// fmt.Println(data.Name)
	// fmt.Println(data.Client)
	// fmt.Println(SLA_ALL_DATA)
	// fmt.Println("")
	// fmt.Println(string(jsonData))
	// fmt.Printf("===============================================================\n")

	contractID := data.RegId
	contractName := data.Name
	contractClient := data.Client

	// A01. ���ID ����մϴ�.
	err = stub.PutState(data.RegId, []byte(content))

	if err != nil {
		return nil, errors.New("Failed to put state with" + content)

	} else {
		fmt.Println("SlaContractRegId : ok")
	}

	// A02. ���� ����մϴ�.
	{
		var err error

		// �������� ���������� ��ȸ�մϴ�.
		contractIDsInBytes, err := stub.GetState(contractName) // ���ϰ� ([]byte, error)

		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractIDsInBytes))
		}

		contractIDsInString := string(contractIDsInBytes)

		//���������� ������� "����"-"���ID���" ���
		if contractIDsInString == "" {
			err = stub.PutState(contractName, []byte(contractID))
			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(contractName, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A03. ����� ����մϴ�.
	{
		var err error

		//�������� ����������ȸ
		contractIDsInBytes, _ := stub.GetState(contractClient) // ���ϰ� ([]byte, error)
		contractIDsInString := string(contractIDsInBytes)

		//���������� ������� "�����"-"���ID���" ���
		if contractIDsInString == "" {
			err = stub.PutState(contractClient, []byte(contractID))
			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(contractClient, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A04. ��ü ��ȸ ����մϴ�.
	{
		var err error

		// �����͸� ��ü ��ȸ�մϴ�.
		contractALLIDsInBytes, err := stub.GetState(SLA_ALL_DATA) // ���ϰ� ([]byte, error)

		fmt.Println("= A04 == 01:" + string(contractALLIDsInBytes))

		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractALLIDsInBytes))
		}

		fmt.Println("= A04 == 02:" + string(contractALLIDsInBytes))

		contractALLIDsInString := string(contractALLIDsInBytes)

		fmt.Println("= A04 == 03:" + contractALLIDsInString)

		// ���������� ������� "����"-"���ID���" ���
		if contractALLIDsInString == "" {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractID))
			fmt.Println("= A04 == 04:")

			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractALLIDsInString+FIELDSEP+contractID))

			fmt.Println("= A04 == 05:")

			if err != nil {
				return nil, err
			}
		}
	}

	/*
		err = stub.PutState(data[0].SlaContractName, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SlaContractName : ok")
		}

		err = stub.PutState(data[0].SlaContractClient, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SlaContractClient : ok")
		}

		err = stub.PutState(SLA_ALL_DATA, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SLA_ALL_DATA : ok")
		}

	*/
	return nil, nil
}

// // ===========================================================
// //  SLAChaincodeStub ������Ʈ �Լ�
// // ===========================================================

// // ����� ������Ʈ�մϴ�. (�⺻)
// func (t *SimpleChaincode) updateContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

// 	var dataInBytes string
// 	var err error

// 	fmt.Printf("updateContractId Input Args:%s\n", args[0])

// 	if len(args) != 1 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithClient")
// 	}

// 	dataInBytes = args[0]
// 	contractID := args[0]

// 	// �������� ��ȸ
// 	contractIDsInBytes, err := stub.GetState(contractID)
// 	if err != nil {
// 		return nil, errors.New("Failed to get state with " + string(contractIDsInBytes))
// 	}

// 	// UPDATDE ó��
// 	stub.PutState(contractID, []byte(args[1]))

// 	// ���泻�� ��ȸ
// 	update_value, err := stub.GetState(contractID)
// 	if err != nil {
// 		return nil, errors.New("Failed to get state with " + dataInBytes)
// 	}

// 	fmt.Printf("slaGetContractsWithClient Response:%s\n", update_value)

// 	return []byte(update_value), nil
// }

// 1.����� ������Ʈ�մϴ�.
func (t *SimpleChaincode) slaUpdateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// create �� ����
	return nil, nil
}

// 2.����� �����մϴ�.
func (t *SimpleChaincode) slaApproveContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// ����ȣ > ��� ������ >  ���� ����
	
	return nil, nil
}

// 3.����� �ݷ��մϴ�.
func (t *SimpleChaincode) slaRejectContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// ����ȣ > ��� ������ >  ���� ����
	
	return nil, nil
}

func (t *SimpleChaincode) slaAbandonContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// ������Ʈ�� abandon state  ����
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// ������Ʈ�� submit state  ����
	return nil, nil
}
func (t *SimpleChaincode) slaCloseContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// ������Ʈ�� submit state  ����
	return nil, nil
}

func (t *SimpleChaincode) slaCreateEvaluationRootFromContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCreateEvaluationsFromContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) slaInitEvaluationValues(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaUpdateEvaluationValues(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaApproveEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaRejectEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaClosePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCloseEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCloseEvaluationRoot(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ===========================================================
//  SLAChaincodeStub �˻� �Լ�
// ===========================================================

// SLA ������ ��ü�� ��ȸ�մϴ�.  (abandon ����)
func (t *SimpleChaincode) slaGetAllContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error

	fmt.Printf("slaGetAllContracts Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithName")
	}

	dataInBytes = args[0]

	// �������� ���ID��� ��ȸ
	contractIDsInBytes, err := stub.GetState(SLA_ALL_DATA)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// ���ID����� ���¸� ��Ʈ������ �迭�� ��ȯ
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// ���ϰ� �ʱ�ȭ
	contractList := make([]string, len(contractIDs))

	// ��� ��ü ID��� ��ȸ
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}

	contractListBytes := strings.Join(contractList, ENTRYSEP)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

// ID���� ��ȸ�մϴ�.
func (t *SimpleChaincode) slaGetContractWithId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error

	fmt.Printf("slaGetContractWithId Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Value to slaGetContractWithId")
	}

	dataInBytes = args[0]
	Valuebytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Failed to get state with" + dataInBytes)
	}

	fmt.Printf("searchbyid Response:%s\n", Valuebytes)

	return Valuebytes, nil
}

// �������� ��ȸ�մϴ�.
func (t *SimpleChaincode) slaGetContractsWithName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error
	var data SlaContract

	fmt.Printf("slaGetContractsWithName Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithName")
	}

	dataInBytes = args[0]
	contractName := args[0]

	// �������� ���ID��� ��ȸ
	contractIDsInBytes, err := stub.GetState(contractName)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// ���ID����� ���¸� ��Ʈ������ �迭�� ��ȯ�մϴ�.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// ���ϰ� �ʱ�ȭ
	contractList := make([]SlaContract, len(contractIDs))

	// ���ID������� ��೻���� �����Ͽ� ����� �ۼ�
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])

		err = json.Unmarshal(contractInBytes, &data)
		contractList[i] = data
	}

	contractListBytes, _ := json.Marshal(contractList)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

// ��������� ��ȸ�մϴ�.
func (t *SimpleChaincode) slaGetContractsWithClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error
	var data SlaContract

	fmt.Printf("slaGetContractsWithClient Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithClient")
	}

	dataInBytes = args[0]
	contractClient := args[0]

	// �������� ���ID��� ��ȸ
	contractIDsInBytes, err := stub.GetState(contractClient)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// ���ID����� ���¸� ��Ʈ������ �迭�� ��ȯ�մϴ�.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// ���ϰ� �ʱ�ȭ
	contractList := make([]SlaContract, len(contractIDs))

	// ���ID������� ��೻���� �����Ͽ� ����� �ۼ�
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])

		err = json.Unmarshal(contractInBytes, &data)
		contractList[i] = data
	}

	contractListBytes, _ := json.Marshal(contractList)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

func (t *SimpleChaincode) slaGetAllEvaluations(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationWithId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationsWithName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationsWithClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
