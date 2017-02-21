class OrderController < ApplicationController
	def purchase

		# Request Parameter 파싱 작업 
		params.to_json
		puts "#{params[:order][:orderNum]}"
		puts "#{params[:order][:product][:productCode]}"
		puts "#{params[:order][:product][:productPrice]}"
		puts "#{params[:order][:accountID]}"
		puts "#{params[:customer][:name]}"
		puts "#{params[:customer][:birth]}"
		puts "#{params[:customer][:cardnum]}"
		puts "#{params[:customer][:UUID]}"
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
