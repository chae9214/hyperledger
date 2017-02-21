class OrderController < ApplicationController
	def purchase

		# Request Parameter 파싱 작업 
		params.to_json
		@orderNum 		= "#{params[:order][:orderNum]}"
		@productCode 	= "#{params[:order][:product][:productCode]}"
		@productPrice 	= "#{params[:order][:product][:productPrice]}"
		@accountID 		= "#{params[:order][:accountID]}"
		@customName 	= "#{params[:customer][:name]}"
		@customBirth 	= "#{params[:customer][:birth]}"
		@customCardnum 	= "#{params[:customer][:cardnum]}"
		@customUUID 	= "#{params[:customer][:UUID]}"
		@deviceID		= "#{params[:device][:deviceID]}"
		@deviceMAC 		= "#{params[:device][:MAC]}"
		@deviceIPAddr	= "#{params[:device][:IPAddr]}"
		@deviceType		= "#{params[:device][:Type]}"

    	@dateTime = Time.now.strftime("%Y-%d-%m %I:%M:%S") 

		#Transaction.new
		@query = "INSERT INTO transactions (	stats , 
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
								) VALUES (	'', 
											'', 
											'"+ @deviceMAC +"',
											'',
											'"+ @dateTime +"',
											'"+ @dateTime +"',
											'"+ @dateTime +"',
											'"+ @customUUID +"',
											'"+ @customCardnum +"',
											'"+ @orderNum +"',
											'"+ @accountID +"',
											'"+ @deviceType +"'
										 ) "
		
		puts @query
		ActiveRecord::Base.connection.execute(@query)

		# 처리결과 응답 
		@returnMsg = '{
			"orderNum":"N16021810001",							
			"result":"Y"
		}'

		render json: @returnMsg
	end
end
