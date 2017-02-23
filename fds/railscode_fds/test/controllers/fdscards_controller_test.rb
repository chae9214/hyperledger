require 'test_helper'

class FdscardsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @fdscard = fdscards(:one)
  end

  test "should get index" do
    get fdscards_url
    assert_response :success
  end

  test "should get new" do
    get new_fdscard_url
    assert_response :success
  end

  test "should create fdscard" do
    assert_difference('Fdscard.count') do
      post fdscards_url, params: { fdscard: { cardnum: @fdscard.cardnum, cid: @fdscard.cid, correspondid: @fdscard.correspondid, customername: @fdscard.customername, fdsproducedby: @fdscard.fdsproducedby, fdsregistreason: @fdscard.fdsregistreason, fdsstatus: @fdscard.fdsstatus, fraudproductcode: @fdscard.fraudproductcode, ipaddr: @fdscard.ipaddr, mac: @fdscard.mac, mobileyn: @fdscard.mobileyn, ordernum: @fdscard.ordernum, tid: @fdscard.tid, transactiondate: @fdscard.transactiondate, transactiontime: @fdscard.transactiontime, uuid: @fdscard.uuid } }
    end

    assert_redirected_to fdscard_url(Fdscard.last)
  end

  test "should show fdscard" do
    get fdscard_url(@fdscard)
    assert_response :success
  end

  test "should get edit" do
    get edit_fdscard_url(@fdscard)
    assert_response :success
  end

  test "should update fdscard" do
    patch fdscard_url(@fdscard), params: { fdscard: { cardnum: @fdscard.cardnum, cid: @fdscard.cid, correspondid: @fdscard.correspondid, customername: @fdscard.customername, fdsproducedby: @fdscard.fdsproducedby, fdsregistreason: @fdscard.fdsregistreason, fdsstatus: @fdscard.fdsstatus, fraudproductcode: @fdscard.fraudproductcode, ipaddr: @fdscard.ipaddr, mac: @fdscard.mac, mobileyn: @fdscard.mobileyn, ordernum: @fdscard.ordernum, tid: @fdscard.tid, transactiondate: @fdscard.transactiondate, transactiontime: @fdscard.transactiontime, uuid: @fdscard.uuid } }
    assert_redirected_to fdscard_url(@fdscard)
  end

  test "should destroy fdscard" do
    assert_difference('Fdscard.count', -1) do
      delete fdscard_url(@fdscard)
    end

    assert_redirected_to fdscards_url
  end
end
