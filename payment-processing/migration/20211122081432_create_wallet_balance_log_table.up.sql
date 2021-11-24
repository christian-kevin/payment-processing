CREATE TABLE wallet_balance_log (
  id BIGINT AUTO_INCREMENT NOT NULL,
  wallet_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY(id),
  KEY k_wallet_id_on_wallet_balance_log(wallet_id) USING BTREE
);