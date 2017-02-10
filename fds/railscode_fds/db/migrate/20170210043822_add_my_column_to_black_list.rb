class AddMyColumnToBlackList < ActiveRecord::Migration[5.0]
  def change
    add_column :black_lists, :test, :string
  end
end
