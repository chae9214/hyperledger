class CreateBlackLists < ActiveRecord::Migration[5.0]
  def change
    create_table :black_lists do |t|
      t.string :cid
      t.string :mac
      t.string :uuid
      t.date :transcationtime
      t.string :registeredby

      t.timestamps
    end
  end
end
