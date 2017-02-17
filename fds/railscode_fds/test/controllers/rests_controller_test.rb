require 'test_helper'

class RestsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @rest = rests(:one)
  end

  test "should get index" do
    get rests_url
    assert_response :success
  end

  test "should get new" do
    get new_rest_url
    assert_response :success
  end

  test "should create rest" do
    assert_difference('Rest.count') do
      post rests_url, params: { rest: { age: @rest.age, job: @rest.job, name: @rest.name } }
    end

    assert_redirected_to rest_url(Rest.last)
  end

  test "should show rest" do
    get rest_url(@rest)
    assert_response :success
  end

  test "should get edit" do
    get edit_rest_url(@rest)
    assert_response :success
  end

  test "should update rest" do
    patch rest_url(@rest), params: { rest: { age: @rest.age, job: @rest.job, name: @rest.name } }
    assert_redirected_to rest_url(@rest)
  end

  test "should destroy rest" do
    assert_difference('Rest.count', -1) do
      delete rest_url(@rest)
    end

    assert_redirected_to rests_url
  end
end
