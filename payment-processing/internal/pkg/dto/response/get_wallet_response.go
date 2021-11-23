package response

type GetWallet struct {
	WalletID int64  `json:"wallet_id"`
	Balance  string `json:"balance"`
	Country  string `json:"country"`
}
