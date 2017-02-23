class AddMyColumnToFdsbanks < ActiveRecord::Migration[5.0]
  def change
    add_column :fdsbanks, :fdsstatus, :string
  end
end
