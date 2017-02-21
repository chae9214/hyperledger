class RestsController < ApplicationController
  before_action :set_rest, only: [:show, :edit, :update, :destroy]

  # GET /rests
  # GET /rests.json
  def index
    @rests = Rest.all

    logger.debug "START HERE\n"
    test_httparty
    logger.debug "END HERE\n"
  end

  # GET /rests/1
  # GET /rests/1.json
  def show
  end

  # GET /rests/new
  def new
    @rest = Rest.new
  end

  # GET /rests/1/edit
  def edit
  end

  # POST /rests
  # POST /rests.json
  def create
    @rest = Rest.new(rest_params)

    respond_to do |format|
      if @rest.save
        format.html { redirect_to @rest, notice: 'Rest was successfully created.' }
        format.json { render :show, status: :created, location: @rest }
      else
        format.html { render :new }
        format.json { render json: @rest.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /rests/1
  # PATCH/PUT /rests/1.json
  def update
    respond_to do |format|
      if @rest.update(rest_params)
        format.html { redirect_to @rest, notice: 'Rest was successfully updated.' }
        format.json { render :show, status: :ok, location: @rest }
      else
        format.html { render :edit }
        format.json { render json: @rest.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /rests/1
  # DELETE /rests/1.json
  def destroy
    @rest.destroy
    respond_to do |format|
      format.html { redirect_to rests_url, notice: 'Rest was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_rest
      @rest = Rest.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def rest_params
      params.require(:rest).permit(:name, :age, :job)
    end

    def test_httparty
      response = HTTParty.get('http://localhost:3000')
      logger.debug response.body, response.code, response.message, response.headers.inspect
    end
end
