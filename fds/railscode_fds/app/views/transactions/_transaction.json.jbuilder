json.extract! transaction, :id, :stats, :cid, :mac, :accountnum, :txtime, :created_at, :updated_at, :cardnum, :ordernum, :posid, :correspondentid
json.url transaction_url(transaction, format: :json)