class RemoveseqFromtransactions < ActiveRecord::Migration[5.0]
def self.up
  remove_column :transactions, :seq
end
end
