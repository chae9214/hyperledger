Rails.application.routes.draw do
  resources :fdscards
  resources :malltransfers
  resources :banktransfers
  # get 'transfer/transfer'
  # post 'transfer/transfer' => 'transfer#transfer'
  # match 'transfer/transfer', to: 'transfer#hh', via: [:get, :post]
  # match 'tra', to: 'tra#index', via: :all

  match 'transfer', to: 'transfer#trans', via: :all
  resources :malls
  resources :fdsbanks
  resources :transactions
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
