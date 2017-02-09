class CreateSlaContracts < ActiveRecord::Migration[5.0]
  def change
    create_table :sla_contracts do |t|
      t.string :SlaContractRegId
      t.string :SlaContractName
      t.string :SlaContractKind
      t.string :SlaContractStaDate
      t.string :SlaContractEndDate
      t.string :SlaContractClient
      t.string :SlaContractClientPerson
      t.string :SlaContractClientPersonTel
      t.string :SlaContractAssessDate
      t.string :SlaContractAssessYn

      t.timestamps
    end
  end
end
