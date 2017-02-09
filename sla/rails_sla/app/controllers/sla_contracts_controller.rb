class SlaContractsController < ApplicationController
  before_action :set_sla_contract, only: [:show, :edit, :update, :destroy]

  # GET /sla_contracts
  # GET /sla_contracts.json
  def index
    @sla_contracts = SlaContract.all
      respond_to do |format|
          format.html { render :index }
          format.json { render :json => @sla_contracts }
      end
  end
  
  
  # GET /sla_contracts/1
  # GET /sla_contracts/1.json
  def show
    @sla_contract = SlaContract.find(params[:id])
  end

  # GET /sla_contracts/new
  # encoding: UTF-8
  def new
    @sla_contract = SlaContract.new
  end

  # GET /sla_contracts/1/edit
  def edit
  end

  # POST /sla_contracts
  # POST /sla_contracts.json
  def create
    #puts call_hyperlegder        
    @sla_contract = SlaContract.new(sla_contract_params)

    respond_to do |format|
      if @sla_contract.save
        @sla_contracts = SlaContract.all
        format.html { render :index}
        format.json { render json:@sla_contracts}
        #format.html { redirect_to @sla_contract, notice: 'Sla contract was successfully created.' }
        #format.json { render :index, status: :created, location: @sla_contract }
      else
        format.html { render :new }
        format.json { render json: @sla_contract.errors, status: :unprocessable_entity }
      end
    end
    
  end



  # PATCH/PUT /sla_contracts/1
  # PATCH/PUT /sla_contracts/1.json
  def update
    respond_to do |format|
      if @sla_contract.update(sla_contract_params)
        format.html { redirect_to @sla_contract, notice: 'Sla contract was successfully updated.' }
        format.json { render :show, status: :ok, location: @sla_contract }
      else
        format.html { render :edit }
        format.json { render json: @sla_contract.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /sla_contracts/1
  # DELETE /sla_contracts/1.json
  def destroy
    @sla_contract.destroy
    respond_to do |format|
      format.html { redirect_to sla_contracts_url, notice: 'Sla contract was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_sla_contract
      @sla_contract = SlaContract.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def sla_contract_params
      params.require(:sla_contract).permit(:SlaContractRegId, :SlaContractName, :SlaContractKind, :SlaContractStaDate, :SlaContractEndDate, :SlaContractClient, :SlaContractClientPerson, :SlaContractClientPersonTel, :SlaContractAssessDate, :SlaContractAssessYn)
    end
    
    
    
    def call_hyperlegder
      uri = URI('http://httpbin.org/post')
      req = Net::HTTP::Post.new(uri)

      json = Hash.new()
      json['jsonrpc'] = "2.0"
      json['method'] = "deploy"

      json['params'] = Hash.new()
      json['params']['type'] = 1

      json['params']['chaincodeID'] = Hash.new()
      json['params']['chaincodeID']['path'] = "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"

      json['params']['ctorMsg'] = Hash.new()

      json['params']['ctorMsg']['args'] = sla_contract_params.values_at(:SlaContractRegId) , sla_contract_params.values_at(:SlaContractRegId, :SlaContractName, :SlaContractKind, :SlaContractStaDate, :SlaContractEndDate, :SlaContractClient, :SlaContractClientPerson, :SlaContractClientPersonTel, :SlaContractAssessDate, :SlaContractAssessYn)
      json['params']['secureContext'] = "bob"
      json['id'] = 1

      req.body = json.to_json

      req.content_type = 'application/json'

      res = Net::HTTP.start(uri.hostname, uri.port) do |http|
        http.request(req)
      end

      return res.body
    end
end
