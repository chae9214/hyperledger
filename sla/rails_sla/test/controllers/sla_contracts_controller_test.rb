require 'test_helper'

class SlaContractsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @sla_contract = sla_contracts(:one)
  end

  test "should get index" do
    get sla_contracts_url
    assert_response :success
  end

  test "should get new" do
    get new_sla_contract_url
    assert_response :success
  end

  test "should create sla_contract" do
    assert_difference('SlaContract.count') do
      post sla_contracts_url, params: { sla_contract: { SlaContractAssessDate: @sla_contract.SlaContractAssessDate, SlaContractAssessYn: @sla_contract.SlaContractAssessYn, SlaContractClient: @sla_contract.SlaContractClient, SlaContractClientPerson: @sla_contract.SlaContractClientPerson, SlaContractClientPersonTel: @sla_contract.SlaContractClientPersonTel, SlaContractEndDate: @sla_contract.SlaContractEndDate, SlaContractKind: @sla_contract.SlaContractKind, SlaContractName: @sla_contract.SlaContractName, SlaContractRegId: @sla_contract.SlaContractRegId, SlaContractStaDate: @sla_contract.SlaContractStaDate } }
    end

    assert_redirected_to sla_contract_url(SlaContract.last)
  end

  test "should show sla_contract" do
    get sla_contract_url(@sla_contract)
    assert_response :success
  end

  test "should get edit" do
    get edit_sla_contract_url(@sla_contract)
    assert_response :success
  end

  test "should update sla_contract" do
    patch sla_contract_url(@sla_contract), params: { sla_contract: { SlaContractAssessDate: @sla_contract.SlaContractAssessDate, SlaContractAssessYn: @sla_contract.SlaContractAssessYn, SlaContractClient: @sla_contract.SlaContractClient, SlaContractClientPerson: @sla_contract.SlaContractClientPerson, SlaContractClientPersonTel: @sla_contract.SlaContractClientPersonTel, SlaContractEndDate: @sla_contract.SlaContractEndDate, SlaContractKind: @sla_contract.SlaContractKind, SlaContractName: @sla_contract.SlaContractName, SlaContractRegId: @sla_contract.SlaContractRegId, SlaContractStaDate: @sla_contract.SlaContractStaDate } }
    assert_redirected_to sla_contract_url(@sla_contract)
  end

  test "should destroy sla_contract" do
    assert_difference('SlaContract.count', -1) do
      delete sla_contract_url(@sla_contract)
    end

    assert_redirected_to sla_contracts_url
  end
end
