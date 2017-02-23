class MalltransfersController < ApplicationController
  before_action :set_malltransfer, only: [:show, :edit, :update, :destroy]

  # GET /malltransfers
  # GET /malltransfers.json
  def index
    @malltransfers = Malltransfer.all
  end

  # GET /malltransfers/1
  # GET /malltransfers/1.json
  def show
  end

  # GET /malltransfers/new
  def new
    @malltransfer = Malltransfer.new
  end

  # GET /malltransfers/1/edit
  def edit
  end

  # POST /malltransfers
  # POST /malltransfers.json
  def create
    @malltransfer = Malltransfer.new(malltransfer_params)

    respond_to do |format|
      if @malltransfer.save
        format.html { redirect_to @malltransfer, notice: 'Malltransfer was successfully created.' }
        format.json { render :show, status: :created, location: @malltransfer }
      else
        format.html { render :new }
        format.json { render json: @malltransfer.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /malltransfers/1
  # PATCH/PUT /malltransfers/1.json
  def update
    respond_to do |format|
      if @malltransfer.update(malltransfer_params)
        format.html { redirect_to @malltransfer, notice: 'Malltransfer was successfully updated.' }
        format.json { render :show, status: :ok, location: @malltransfer }
      else
        format.html { render :edit }
        format.json { render json: @malltransfer.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /malltransfers/1
  # DELETE /malltransfers/1.json
  def destroy
    @malltransfer.destroy
    respond_to do |format|
      format.html { redirect_to malltransfers_url, notice: 'Malltransfer was successfully destroyed.' }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_malltransfer
      @malltransfer = Malltransfer.find(params[:id])
    end

    # Never trust parameters from the scary internet, only allow the white list through.
    def malltransfer_params
      params.fetch(:malltransfer, {})
    end
end
