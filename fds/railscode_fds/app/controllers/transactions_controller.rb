class TransactionsController < ApplicationController
  before_action :set_transaction, only: [:show, :edit, :update, :destroy]

  # GET /transactions
  # GET /transactions.json
  def index
    @transactions = Transaction.all
    #초기 화면에서 hyperledger 데이터 호출
    @hyperledger_response = JSON.parse(query_from_hyperledger("fdsGetAllFraudEntries",nil))
    logger.debug "hyperledger_response$$$$$$$$$$$$$ #{@hyperledger_response}"
    @hyperledger_result_list = @hyperledger_response["result"]["message"]
    logger.debug "hyperledger_response_message$$$$$$$$$$$$$ #{@hyperledger_result_list}"
    @parsed_hyperledger_result_list = JSON.parse(@hyperledger_result_list.tr("\\", ""))
    logger.debug "parsed_hyperledger_response_message$$$$$$$$$$$$$ #{@parsed_hyperledger_result_list}"

    respond_to do |format|
      format.html { render :index }
      format.json { render json: @parsed_hyperledger_result_list}
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

    respond_to do |format|
      if @transaction.save
        format.html { redirect_to @transaction, notice: '거래 내역이 정상 등록 되었습니다.' }
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
        format.html { redirect_to @transaction, notice: '거래 내역이 수정 되었습니다.' }
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
      format.html { redirect_to transactions_url, notice: '거래 내역이 삭제 되었습니다.' }
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
      params.require(:transaction).permit(:uuid, :stats, :cid, :mac, :accountnum, :txtime, :cardnum, :ordernum, :posid, :correspondentid)
    end

    def reload_currentPage 
      respond_to do |format|
        format.js {render inline: "location.reload();" }
      end
    end

    def query_from_hyperledger(key,value)
      uri = URI('http://192.168.99.100:7050/chaincode')
      req = Net::HTTP::Post.new(uri)

      json = Hash.new()
      json['jsonrpc'] = "2.0"
      json['method'] = "query"

      json['params'] = Hash.new()
      json['params']['type'] = 1

      json['params']['chaincodeID'] = Hash.new()
      json['params']['chaincodeID']['name'] = "mycc"

      json['params']['ctorMsg'] = Hash.new()

      if value == nil 
       json['params']['ctorMsg']['args'] = [key]  
      else 
        json['params']['ctorMsg']['args'] = [key, value]
      end

      json['params']['secureContext'] = "admin"
      json['id'] = 1    

      req.body = json.to_json

      req.content_type = 'application/json'
      res = Net::HTTP.start(uri.hostname, uri.port) do |http|
        puts req.body
        http.request(req)
      end

      case res
      when Net::HTTPSuccess, Net::HTTPRedirection
        return res.body
      else
        res.value
      end
    end
    helper_method :query_from_hyperledger
end
