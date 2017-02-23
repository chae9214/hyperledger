require 'test_helper'

class FdsbanksControllerTest < ActionDispatch::IntegrationTest
  setup do
    @fdsbank = fdsbanks(:one)
  end

  test "should get index" do
    get fdsbanks_url
    assert_response :success
  end

  test "should get new" do
    get new_fdsbank_url
    assert_response :success
  end

  test "should create fdsbank" do
    assert_difference('Fdsbank.count') do
      post fdsbanks_url, params: { fdsbank: { accountnum: @fdsbank.accountnum, cid: @fdsbank.cid, fdsauthresult: @fdsbank.fdsauthresult, fdsproducedby: @fdsbank.fdsproducedby, fdsregistreason: @fdsbank.fdsregistreason, identitynum: @fdsbank.identitynum, ipaddress: @fdsbank.ipaddress, mac: @fdsbank.mac, mobileyn: @fdsbank.mobileyn, posid: @fdsbank.posid, tid: @fdsbank.tid, transactiondate: @fdsbank.transactiondate, transactiontime: @fdsbank.transactiontime, uuid: @fdsbank.uuid } }
    end

    assert_redirected_to fdsbank_url(Fdsbank.last)
  end

  test "should show fdsbank" do
    get fdsbank_url(@fdsbank)
    assert_response :success
  end

  test "should get edit" do
    get edit_fdsbank_url(@fdsbank)
    assert_response :success
  end

  test "should update fdsbank" do
    patch fdsbank_url(@fdsbank), params: { fdsbank: { accountnum: @fdsbank.accountnum, cid: @fdsbank.cid, fdsauthresult: @fdsbank.fdsauthresult, fdsproducedby: @fdsbank.fdsproducedby, fdsregistreason: @fdsbank.fdsregistreason, identitynum: @fdsbank.identitynum, ipaddress: @fdsbank.ipaddress, mac: @fdsbank.mac, mobileyn: @fdsbank.mobileyn, posid: @fdsbank.posid, tid: @fdsbank.tid, transactiondate: @fdsbank.transactiondate, transactiontime: @fdsbank.transactiontime, uuid: @fdsbank.uuid } }
    assert_redirected_to fdsbank_url(@fdsbank)
  end

  test "should destroy fdsbank" do
    assert_difference('Fdsbank.count', -1) do
      delete fdsbank_url(@fdsbank)
    end

    assert_redirected_to fdsbanks_url
  end
end
