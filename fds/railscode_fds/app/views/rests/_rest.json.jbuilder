json.extract! rest, :id, :name, :age, :job, :created_at, :updated_at
json.url rest_url(rest, format: :json)