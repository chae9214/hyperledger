class CreateRests < ActiveRecord::Migration[5.0]
  def change
    create_table :rests do |t|
      t.string :name
      t.integer :age
      t.string :job

      t.timestamps
    end
  end
end
