require 'test_helper'

class BanktransfersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @banktransfer = banktransfers(:one)
  end

  test "should get index" do
    get banktransfers_url
    assert_response :success
  end

  test "should get new" do
    get new_banktransfer_url
    assert_response :success
  end

  test "should create banktransfer" do
    assert_difference('Banktransfer.count') do
      post banktransfers_url, params: { banktransfer: {  } }
    end

    assert_redirected_to banktransfer_url(Banktransfer.last)
  end

  test "should show banktransfer" do
    get banktransfer_url(@banktransfer)
    assert_response :success
  end

  test "should get edit" do
    get edit_banktransfer_url(@banktransfer)
    assert_response :success
  end

  test "should update banktransfer" do
    patch banktransfer_url(@banktransfer), params: { banktransfer: {  } }
    assert_redirected_to banktransfer_url(@banktransfer)
  end

  test "should destroy banktransfer" do
    assert_difference('Banktransfer.count', -1) do
      delete banktransfer_url(@banktransfer)
    end

    assert_redirected_to banktransfers_url
  end
end
