Rails.application.routes.draw do
  resources :banktransfers
  # get 'transfer/transfer'
  # post 'transfer/transfer' => 'transfer#transfer'
  # match 'transfer/transfer', to: 'transfer#hh', via: [:get, :post]
  # match 'tra', to: 'tra#index', via: :all

  match 'transfer', to: 'transfer#trans', via: :all
  post 'order', to: 'order#purchase', via: :all
  resources :rests
  resources :transactions
  resources :black_lists
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
