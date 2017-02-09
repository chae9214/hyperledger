# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20170207064051) do

  create_table "sla_contracts", force: :cascade do |t|
    t.string   "SlaContractRegId"
    t.string   "SlaContractName"
    t.string   "SlaContractKind"
    t.string   "SlaContractStaDate"
    t.string   "SlaContractEndDate"
    t.string   "SlaContractClient"
    t.string   "SlaContractClientPerson"
    t.string   "SlaContractClientPersonTel"
    t.string   "SlaContractAssessDate"
    t.string   "SlaContractAssessYn"
    t.datetime "created_at",                 null: false
    t.datetime "updated_at",                 null: false
  end

end
