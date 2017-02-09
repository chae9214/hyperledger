class TransactionsController < ApplicationController
  before_action :set_transaction, only: [:show, :edit, :update, :destroy]

  # GET /transactions
  # GET /transactions.json
  def index
    @transactions = Transaction.all
    #초기 화면에서 hyperledger 데이터 호출
    @hyperledger_response = JSON.parse(query_from_hyperledger)
    @hyperledger_result_list = @hyperledger_response["result"]["message"]
    @hyperledger_response_array = @hyperledger_result_list.split("$")
   
    #테스트 데이터
    #@dummy_data = "123456|AA-BB-CC-DD|3fdd6dfg98ac|2017-07-07|신한은행$234567|AA-BB-CC-DD|3fdd6dfg98ac|2017-07-07|신한카드"
    
    #Entry별로 배열 생성
    @cnt = 0
    @hyperledger_row_array = [@hyperledger_response_array.length]
    begin
      @hyperledger_row_array[@cnt] = @hyperledger_response_array[@cnt].split("|")
      @cnt +=1
    end while @cnt < @hyperledger_response_array.length

    logger.debug "hyperledger row data $$$$$$$$$$$$$ #{@hyperledger_row_array}" 
        
    respond_to do |format|
      format.html { render :index }
      format.json { render json: @hyperledger_row_array}
    end

  end

  # GET /transactions/1
  # GET /transactions/1.json
  def show
  end

  # GET /transactions/new
  def new
    @transaction = Transaction.new
  end

  # GET /transactions/1/edit
  def edit
  end

  # POST /transactions
  # POST /transactions.json
  def create
    @transaction = Transaction.new(transaction_params)

=begin
    #HyperLedger 데이터를 출력하기 위한 객체 생성
    @hyperledger_list = Transaction.new()

    #HyperLedger 응답값 저장
    @hyperledger_data = put_hyperledger
    logger.debug "hyperledger_response $$$$$$$$$$$$ : #{@hyperledger_data}"
    
    #HyperLedger 응답값 저장 객체
    @hyperledger_response = JSON.parse(put_hyperledger)
    logger.debug "parsed_hyperledger_response $$$$$$$$$$$$$$$ : #{@hyperledger_response}"

    
    #추후에는 HyperLedger 실 데이터로 변경 처리
    @dummy_data = "123456|AA-BB-CC-DD|3fdd6dfg98ac|2017-07-07|신한은행"
    responseArr = @dummy_data.split("|")
    logger.debug "split data $$$$$$$$$$$$$ #{responseArr}" 
=end

    respond_to do |format|
      if @transaction.save
        format.html { redirect_to @transaction, notice: 'Transaction was successfully created.' }
        format.json { render :show, status: :created, location: @transaction }
      else
        format.html { render :new }
        format.json { render json: @transaction.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /transactions/1
  # PATCH/PUT /transactions/1.json
  def update
    respond_to do |format|
      if @transaction.update(transaction_params)
        format.html { redirect_to @transaction, notice: 'Transaction was successfully updated.' }
        format.json { render :show, status: :ok, location: @transaction }
      else
        format.html { render :edit }
        format.json { render json: @transaction.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /transactions/1
  # DELETE /transactions/1.json
  def destroy
    @transaction.destroy
    respond_to do |format|
      format.html { redirect_to transactions_url, notice: 'Transaction was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_transaction
      @transaction = Transaction.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def transaction_params
      params.require(:transaction).permit(:seq, :stats, :cid, :mac, :accountnum, :txtime)
    end

=begin
    #등록 테스트를 위한 function
    def register_to_hyperledger
      uri = URI('http://192.168.150.129:7050/chaincode')
      req = Net::HTTP::Post.new(uri)

      json = Hash.new()
      json['jsonrpc'] = "2.0"
      json['method'] = "invoke"

      json['params'] = Hash.new()
      json['params']['type'] = 1

      json['params']['chaincodeID'] = Hash.new()
      json['params']['chaincodeID']['name'] = "mycc"

      json['params']['ctorMsg'] = Hash.new()
      json['params']['ctorMsg']['args'] = [ "register", "12345678" , "AA-BB-CC-DD", "uuid-uuid-uuid", "2017-02-08 14:00:00", "신한 카드", "6", "7", "8" ]
      json['params']['secureContext'] = "bob"
      json['id'] = 1

      req.body = json.to_json

      req.content_type = 'application/json'

      res = Net::HTTP.start(uri.hostname, uri.port) do |http|
        puts req.body
        http.request(req)
      end

      return res.body
    end
=end

    def query_from_hyperledger
      uri = URI('http://192.168.99.101:7050/chaincode')
      req = Net::HTTP::Post.new(uri)

      json = Hash.new()
      json['jsonrpc'] = "2.0"
      json['method'] = "query"

      json['params'] = Hash.new()
      json['params']['type'] = 1

      json['params']['chaincodeID'] = Hash.new()
      json['params']['chaincodeID']['name'] = "mycc3"

      json['params']['ctorMsg'] = Hash.new()
      json['params']['ctorMsg']['args'] = [ "lookupwithcid" , "cid1"]
      json['params']['secureContext'] = "admin"
      json['id'] = 3

      req.body = json.to_json

      req.content_type = 'application/json'

      res = Net::HTTP.start(uri.hostname, uri.port) do |http|
        puts req.body
        http.request(req)
      end

      return res.body
    end


    def invoke_to_hyperledger
        uri = URI('http://192.168.99.101:7050/chaincode')
      req = Net::HTTP::Post.new(uri)

      json = Hash.new()
      json['jsonrpc'] = "2.0"
      json['method'] = "invoke"

      json['params'] = Hash.new()
      json['params']['type'] = 1

      json['params']['chaincodeID'] = Hash.new()
      json['params']['chaincodeID']['name'] = "mycc"

      json['params']['ctorMsg'] = Hash.new()
      json['params']['ctorMsg']['args'] = [ "register","cidValue","macValue","uuidValue","2017-01-01","e","f","g","h"]
      json['params']['secureContext'] = "bob"
      json['id'] = 1

      req.body = json.to_json

      req.content_type = 'application/json'

      res = Net::HTTP.start(uri.hostname, uri.port) do |http|
        puts req.body
        http.request(req)
      end

      return res.body
    end
end
