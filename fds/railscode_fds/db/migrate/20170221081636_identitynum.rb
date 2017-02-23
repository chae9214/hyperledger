class Identitynum < ActiveRecord::Migration[5.0]
def self.up
  remove_column :fdsbanks, :identitynum
end
end
