class CreateFdscards < ActiveRecord::Migration[5.0]
  def change
    create_table :fdscards do |t|
      t.string :tid
      t.string :cid
      t.string :mac
      t.string :uuid
      t.string :customername
      t.string :transactiondate
      t.string :transactiontime
      t.string :fdsproducedby
      t.string :fdsregistreason
      t.string :ordernum
      t.string :fraudproductcode
      t.string :fdsstatus
      t.string :correspondid
      t.string :ipaddr
      t.string :mobileyn
      t.string :cardnum

      t.timestamps
    end
  end
end
