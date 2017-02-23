class FdscardsController < ApplicationController
  before_action :set_fdscard, only: [:show, :edit, :update, :destroy]

  # GET /fdscards
  # GET /fdscards.json
  def index
    @fdscards = Fdscard.all

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

  # GET /fdscards/1
  # GET /fdscards/1.json
  def show
  end

  # GET /fdscards/new
  def new
    @fdscard = Fdscard.new
  end

  # GET /fdscards/1/edit
  def edit
  end

  # POST /fdscards
  # POST /fdscards.json
  def create
    @fdscard = Fdscard.new(fdscard_params)

    respond_to do |format|
      if @fdscard.save
        format.html { redirect_to @fdscard, notice: 'Fdscard was successfully created.' }
        format.json { render :show, status: :created, location: @fdscard }
      else
        format.html { render :new }
        format.json { render json: @fdscard.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /fdscards/1
  # PATCH/PUT /fdscards/1.json
  def update
    respond_to do |format|
      if @fdscard.update(fdscard_params)
        format.html { redirect_to @fdscard, notice: 'Fdscard was successfully updated.' }
        format.json { render :show, status: :ok, location: @fdscard }
      else
        format.html { render :edit }
        format.json { render json: @fdscard.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /fdscards/1
  # DELETE /fdscards/1.json
  def destroy
    @fdscard.destroy
    respond_to do |format|
      format.html { redirect_to fdscards_url, notice: 'Fdscard was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_fdscard
      @fdscard = Fdscard.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def fdscard_params
      params.require(:fdscard).permit(:tid, :cid, :mac, :uuid, :customername, :transactiondate, :transactiontime, :fdsproducedby, :fdsregistreason, :ordernum, :fraudproductcode, :fdsstatus, :correspondid, :ipaddr, :mobileyn, :cardnum)
    end

    def reload_currentPage 
      respond_to do |format|
        format.js {render inline: "location.reload();" }
      end
    end

    def query_from_hyperledger(key,value)
      uri = URI('http://10.243.224.161:7050/chaincode')
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
