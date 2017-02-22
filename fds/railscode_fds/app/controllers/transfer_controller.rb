class TransferController < ApplicationController
def empty
  render :nothing => true
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

    	@dateTime = Time.now.strftime("%Y-%m-%d %H:%M:%S") 
    	@uuid = "1121 44356 1123"
    	@cardnum = "1234-1234-1234"

		#이상 거래 판단: 신한은행 -> 정상 
		puts "#{params[:depositAccount][:bankCode]}"
		if params[:depositAccount][:bankCode] == "신한은행"
			params[:result] = "Y"
			@returnMsg = params
			@fdsStats = "정상"
			makeQuery(@fdsStats)
			ActiveRecord::Base.connection.execute(@query)
		else
			params[:result] = "N"
			@returnMsg = params
			@fdsStats = "이상"
			makeQuery(@fdsStats)
			ActiveRecord::Base.connection.execute(@query)
		end	

		render json: @returnMsg
	end
end