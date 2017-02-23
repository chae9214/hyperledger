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

ActiveRecord::Schema.define(version: 20170223075510) do

  create_table "banktransfers", force: :cascade do |t|
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "black_lists", force: :cascade do |t|
    t.string   "cid"
    t.string   "mac"
    t.string   "uuid"
    t.date     "transcationtime"
    t.string   "registeredby"
    t.datetime "created_at",      null: false
    t.datetime "updated_at",      null: false
    t.string   "test"
  end

  create_table "fdsbanks", force: :cascade do |t|
    t.string   "tid"
    t.string   "cid"
    t.string   "mac"
    t.string   "uuid"
    t.string   "transactiondate"
    t.string   "transactiontime"
    t.string   "fdsauthresult"
    t.string   "fdsproducedby"
    t.string   "fdsregistreason"
    t.string   "ipaddress"
    t.string   "posid"
    t.string   "mobileyn"
    t.string   "accountnum"
    t.datetime "created_at",      null: false
    t.datetime "updated_at",      null: false
    t.string   "fdsstatus"
  end

  create_table "malltransfers", force: :cascade do |t|
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "rests", force: :cascade do |t|
    t.string   "name"
    t.integer  "age"
    t.string   "job"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "transactions", force: :cascade do |t|
    t.string   "stats"
    t.string   "cid"
    t.string   "mac"
    t.string   "accountnum"
    t.datetime "txtime"
    t.datetime "created_at",      null: false
    t.datetime "updated_at",      null: false
    t.string   "uuid"
    t.string   "cardnum"
    t.string   "ordernum"
    t.string   "correspondentid"
    t.string   "posid"
  end

end
