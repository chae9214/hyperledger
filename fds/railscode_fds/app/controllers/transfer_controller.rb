class TransferController < ApplicationController
	def trans
		
		# Request Parameter 파싱 작업 
		params.to_json
		puts "#{params[:ordernum]}"
		puts "#{params[:customer][:CID]}"
		puts "#{params[:customer][:name]}"
		puts "#{params[:customer][:bankcode]}"
		puts "#{params[:customer][:account]}"
		puts "#{params[:depositAccount][:bankCode]}"
		puts "#{params[:depositAccount][:customerName]}"
		puts "#{params[:depositAccount][:account]}"
		puts "#{params[:depositAccount][:money]}"
		puts "#{params[:device][:deviceID]}"
		puts "#{params[:device][:MAC]}"
		puts "#{params[:device][:IPAddr]}"
		puts "#{params[:device][:Type]}"

		# 처리결과 응답 
		@respJson = '{
			"orderNum":"N16021810001",							
			"result":"Y"
		}'

		render json: @respJson
	end
end
