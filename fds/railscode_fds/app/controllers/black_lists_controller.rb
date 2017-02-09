class BlackListsController < ApplicationController
	before_action :set_black_list, only: [:show, :edit, :update, :destroy]

	# GET /black_lists
	# GET /black_lists.json
	def index
		@black_lists = BlackList.all
		@hyperledger_initialized_list = JSON.parse(call_hyperledger)

		logger.debug "hyperledger_initialized_list $$$$$$$$$$$$  : #{@hyperledger_initialized_list}"
	end

	# GET /black_lists/1
	# GET /black_lists/1.json
	def show
	end

	# GET /black_lists/new
	def new
		@black_list = BlackList.new
	end

	# GET /black_lists/1/edit
	def edit
	end

	# POST /black_lists
	# POST /black_lists.json
	def create
		@black_list = BlackList.new(black_list_params) 
		#HyperLedger 데이터를 출력하기 위한 객체 생성
		@hyperledger_list = BlackList.new()

		#HyperLedger 응답값 저장
		@hyperledger_data = put_hyperledger

		logger.debug "hyperledger_response $$$$$$$$$$$$ : #{@hyperledger_data}"
		logger.debug black_list_params

		#HyperLedger 응답값 저장 객체
		@hyperledger_response = JSON.parse(put_hyperledger)

		logger.debug "parsed_hyperledger_response $$$$$$$$$$$$$$$ : #{@hyperledger_response}"

=begin		
		#추후에는 HyperLedger 실 데이터로 변경 처리
		@dummy_data = "123456|AA-BB-CC-DD|3fdd6dfg98ac|2017-07-07|신한은행"
		responseArr = @dummy_data.split("|")
		logger.debug "split data $$$$$$$$$$$$$ #{responseArr}" 
=end

		

		@black_list.cid = responseArr[0]
		@black_list.mac = responseArr[1]
		@black_list.uuid = responseArr[2]
		@black_list.transcationtime = responseArr[3]    
		@black_list.registeredby = responseArr[4]       


		@hyperledger_list.cid = responseArr[0]
		@hyperledger_list.mac = responseArr[1]
		@hyperledger_list.uuid = responseArr[2]
		@hyperledger_list.transcationtime = responseArr[3]    
		@hyperledger_list.registeredby = responseArr[4]       

		logger.debug "hyperledger_list.cid $$$$$$$$$$$$$$$$$ : #{@hyperledger_list.cid}"
		logger.debug "hyperledger_list.mac $$$$$$$$$$$$$$$$$ : #{@hyperledger_list.mac}"
		logger.debug "hyperledger_list.uuid $$$$$$$$$$$$$$$$$ : #{@hyperledger_list.uuid}"
		logger.debug "hyperledger_list.transcationtime $$$$$$$$$$$$$$$$$ : #{@hyperledger_list.transcationtime}"
		logger.debug "hyperledger_list.registeredby $$$$$$$$$$$$$$$$$ : #{@hyperledger_list.registeredby}"

		respond_to do |format|
				if @black_list.save
				format.html { redirect_to @black_list, notice: 'Black list was successfully created.' }
				format.json { render :show, status: :created, location: @black_list }
			else
				format.html { render :new }
				format.json { render json: @black_list.errors, status: :unprocessable_entity }
			end
		end
	end

	# PATCH/PUT /black_lists/1
	# PATCH/PUT /black_lists/1.json
	def update
		respond_to do |format|
			if @black_list.update(black_list_params)
				format.html { redirect_to @black_list, notice: 'Black list was successfully updated.' }
				format.json { render :show, status: :ok, location: @black_list }
			else
				format.html { render :edit }
				format.json { render json: @black_list.errors, status: :unprocessable_entity }
			end
		end
	end

	# DELETE /black_lists/1
	# DELETE /black_lists/1.json
	def destroy
		@black_list.destroy
		respond_to do |format|
			format.html { redirect_to black_lists_url, notice: 'Black list was successfully destroyed.' }
			format.json { head :no_content }
		end
	end

	private
		# Use callbacks to share common setup or constraints between actions.
		def set_black_list
			@black_list = BlackList.find(params[:id])
		end

		# Never trust parameters from the scary internet, only allow the white list through.
		def black_list_params
			 params.require(:black_list).permit(:cid, :mac, :uuid, :transcationtime, :registeredby)
		end

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
			json['params']['ctorMsg']['args'] = [ "register", @black_list.cid , @black_list.mac, @black_list.uuid, @black_list.transcationtime, @black_list.registeredby, "6", "7", "8" ]
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

		def query_from_hyperledger
				uri = URI('http://192.168.150.129:7050/chaincode')
			req = Net::HTTP::Post.new(uri)

			json = Hash.new()
			json['jsonrpc'] = "2.0"
			json['method'] = "query"

			json['params'] = Hash.new()
			json['params']['type'] = 1

			json['params']['chaincodeID'] = Hash.new()
			json['params']['chaincodeID']['name'] = "mycc"

			json['params']['ctorMsg'] = Hash.new()
			json['params']['ctorMsg']['args'] = [ "lookupwithcid" , "a"]
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
