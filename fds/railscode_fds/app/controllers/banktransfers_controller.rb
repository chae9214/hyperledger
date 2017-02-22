class BanktransfersController < ApplicationController
  before_action :set_banktransfer, only: [:show, :edit, :update, :destroy]

  # GET /banktransfers
  # GET /banktransfers.json
  def index
    @banktransfers = Banktransfer.all
  end

  # GET /banktransfers/1
  # GET /banktransfers/1.json
  def show
  end

  # GET /banktransfers/new
  def new
    @banktransfer = Banktransfer.new
  end

  # GET /banktransfers/1/edit
  def edit
  end

  # POST /banktransfers
  # POST /banktransfers.json
  def create
    @banktransfer = Banktransfer.new(banktransfer_params)

    respond_to do |format|
      if @banktransfer.save
        format.html { redirect_to @banktransfer, notice: 'Banktransfer was successfully created.' }
        format.json { render :show, status: :created, location: @banktransfer }
      else
        format.html { render :new }
        format.json { render json: @banktransfer.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /banktransfers/1
  # PATCH/PUT /banktransfers/1.json
  def update
    respond_to do |format|
      if @banktransfer.update(banktransfer_params)
        format.html { redirect_to @banktransfer, notice: 'Banktransfer was successfully updated.' }
        format.json { render :show, status: :ok, location: @banktransfer }
      else
        format.html { render :edit }
        format.json { render json: @banktransfer.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /banktransfers/1
  # DELETE /banktransfers/1.json
  def destroy
    @banktransfer.destroy
    respond_to do |format|
      format.html { redirect_to banktransfers_url, notice: 'Banktransfer was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_banktransfer
      @banktransfer = Banktransfer.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def banktransfer_params
      params.fetch(:banktransfer, {})
    end
end
