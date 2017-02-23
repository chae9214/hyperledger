require 'test_helper'

class MalltransfersControllerTest < ActionDispatch::IntegrationTest
  setup do
    @malltransfer = malltransfers(:one)
  end

  test "should get index" do
    get malltransfers_url
    assert_response :success
  end

  test "should get new" do
    get new_malltransfer_url
    assert_response :success
  end

  test "should create malltransfer" do
    assert_difference('Malltransfer.count') do
      post malltransfers_url, params: { malltransfer: {  } }
    end

    assert_redirected_to malltransfer_url(Malltransfer.last)
  end

  test "should show malltransfer" do
    get malltransfer_url(@malltransfer)
    assert_response :success
  end

  test "should get edit" do
    get edit_malltransfer_url(@malltransfer)
    assert_response :success
  end

  test "should update malltransfer" do
    patch malltransfer_url(@malltransfer), params: { malltransfer: {  } }
    assert_redirected_to malltransfer_url(@malltransfer)
  end

  test "should destroy malltransfer" do
    assert_difference('Malltransfer.count', -1) do
      delete malltransfer_url(@malltransfer)
    end

    assert_redirected_to malltransfers_url
  end
end
