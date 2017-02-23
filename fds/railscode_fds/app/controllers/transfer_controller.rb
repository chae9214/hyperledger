class TransferController < ApplicationController
	def empty
	  render :nothing => true
	end

	def query_from_hyperledger(key,value)
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

	def makeQuery (stats)
			@query = "INSERT INTO transactions (stats, 
											cid, 
											mac, 
											accountnum, 
											txtime, 
											created_at, 
											updated_at, 
											uuid, 
											cardnum, 
											ordernum, 
											correspondentid, 
											posid
								) VALUES ('"+stats +"',
										  '"+params[:customer][:CID]+"',
										  '"+params[:device][:MAC]+"',
										  '"+params[:customer][:account]+"',
										  '"+@dateTime+"',
										  '"+@dateTime+"',
										  '"+@dateTime+"',
										  '"+@uuid+"',
										  '"+@cardnum+"',
										  '"+params[:ordernum]+"',
										  '"+params[:device][:deviceID]+"',
										  '"+params[:device][:deviceID]+"'
										 ) "
	end

	def trans
=begin
		puts "#{params}"
		puts "#{params[:ordernum]}"
		puts "#{params[:customer][:CID]}"  #출금 고객 번호
		puts "#{params[:customer][:name]}" #출금 고객명
		puts "#{params[:customer][:bankcode]}" #출금 은행
		puts "#{params[:customer][:account]}" #출금 계좌 
		puts "#{params[:depositAccount][:bankCode]}" #입금 은행
		puts "#{params[:depositAccount][:customerName]}" #입금 고객명
		puts "#{params[:depositAccount][:account]}" #입금 계쫘
		puts "#{params[:depositAccount][:money]}" #입금 금액
		puts "#{params[:device][:deviceID]}" #단말기 ID
		puts "#{params[:device][:MAC]}" #MAC
		puts "#{params[:device][:IPAddr]}" #IP 주소
		puts "#{params[:device][:Type]}" #모바일 여부
=end

    	@dateTime = Time.now.strftime("%Y-%m-%d %H:%M:%S") 
    	#Random Uuid 생성
    	require 'securerandom'
    	@uuid = p SecureRandom.uuid
    	#Random Card Number 생성
    	@randomNum1 = p SecureRandom.random_number(9999)
    	@randomNum2 = p SecureRandom.random_number(9999)
    	@randomNum3 = p SecureRandom.random_number(9999)
    	@randomNum4 = p SecureRandom.random_number(9999)
    	@cardnum = "#{@randomNum1}-#{@randomNum2}-#{@randomNum3}-#{@randomNum4}"
	    
		#이상 거래 판단: 모바일에서 넘어온 CID로 Hyperledger를 호출하여 결과값으로 이상거래 구분
		@hyperledger_response = JSON.parse(query_from_hyperledger("fdsGetFraudEntriesWithCid",params[:customer][:CID]))
	    if @hyperledger_response["result"]["message"].length > 2
	    	params[:result] = "N"
			@returnMsg = params
			@fdsStats = "이상"
			makeQuery(@fdsStats)
			ActiveRecord::Base.connection.execute(@query)
	    else
			params[:result] = "Y"
			@returnMsg = params
			@fdsStats = "정상"
			makeQuery(@fdsStats)
			ActiveRecord::Base.connection.execute(@query)
	    end

		render json: @returnMsg
	end
end