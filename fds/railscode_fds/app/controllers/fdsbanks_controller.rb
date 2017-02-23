class FdsbanksController < ApplicationController
  before_action :set_fdsbank, only: [:show, :edit, :update, :destroy]

  # GET /fdsbanks
  # GET /fdsbanks.json
  def index
    @fdsbanks = Fdsbank.all
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

  # GET /fdsbanks/1
  # GET /fdsbanks/1.json
  def show
  end

  # GET /fdsbanks/new
  def new
    @fdsbank = Fdsbank.new
  end

  # GET /fdsbanks/1/edit
  def edit
  end

  # POST /fdsbanks
  # POST /fdsbanks.json
  def create
    @fdsbank = Fdsbank.new(fdsbank_params)

    respond_to do |format|
      if @fdsbank.save
        format.html { redirect_to @fdsbank, notice: 'Fdsbank was successfully created.' }
        format.json { render :show, status: :created, location: @fdsbank }
      else
        format.html { render :new }
        format.json { render json: @fdsbank.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /fdsbanks/1
  # PATCH/PUT /fdsbanks/1.json
  def update
    respond_to do |format|
      if @fdsbank.update(fdsbank_params)
        format.html { redirect_to @fdsbank, notice: 'Fdsbank was successfully updated.' }
        format.json { render :show, status: :ok, location: @fdsbank }
      else
        format.html { render :edit }
        format.json { render json: @fdsbank.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /fdsbanks/1
  # DELETE /fdsbanks/1.json
  def destroy
    @fdsbank.destroy
    respond_to do |format|
      format.html { redirect_to fdsbanks_url, notice: 'Fdsbank was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_fdsbank
      @fdsbank = Fdsbank.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def fdsbank_params
      params.require(:fdsbank).permit(:tid, :cid, :mac, :uuid, :fdsstatus, :transactiondate, :transactiontime, :fdsauthresult, :fdsproducedby, :fdsregistreason, :ipaddress, :posid, :mobileyn, :accountnum)
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
