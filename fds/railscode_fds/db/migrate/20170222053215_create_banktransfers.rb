class CreateBanktransfers < ActiveRecord::Migration[5.0]
  def change
    create_table :banktransfers do |t|

      t.timestamps
    end
  end
end
