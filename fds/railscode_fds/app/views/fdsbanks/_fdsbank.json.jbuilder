json.extract! fdsbank, :id, :tid, :cid, :mac, :uuid, :fdsstatus, :transactiondate, :transactiontime, :fdsauthresult, :fdsproducedby, :fdsregistreason, :ipaddress, :posid, :mobileyn, :accountnum, :created_at, :updated_at
json.url fdsbank_url(fdsbank, format: :json)