json.extract! transaction, :id, :seq, :stats, :cid, :mac, :accountnum, :txtime, :created_at, :updated_at
json.url transaction_url(transaction, format: :json)