class CreateFdsbanks < ActiveRecord::Migration[5.0]
  def change
    create_table :fdsbanks do |t|
      t.string :tid
      t.string :cid
      t.string :mac
      t.string :uuid
      t.string :identitynum
      t.string :transactiondate
      t.string :transactiontime
      t.string :fdsauthresult
      t.string :fdsproducedby
      t.string :fdsregistreason
      t.string :ipaddress
      t.string :posid
      t.string :mobileyn
      t.string :accountnum

      t.timestamps
    end
  end
end
