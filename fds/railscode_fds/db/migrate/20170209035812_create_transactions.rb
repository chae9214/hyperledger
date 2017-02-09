class CreateTransactions < ActiveRecord::Migration[5.0]
  def change
    create_table :transactions do |t|
      t.integer :seq
      t.string :stats
      t.string :cid
      t.string :mac
      t.string :accountnum
      t.timestamp :txtime

      t.timestamps
    end
  end
end
