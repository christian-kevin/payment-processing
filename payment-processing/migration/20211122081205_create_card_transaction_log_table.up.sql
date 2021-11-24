CREATE TABLE card_transaction_log (
  id BIGINT AUTO_INCREMENT NOT NULL,
  card_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY(id),
  KEY k_card_id_on_card_transaction_log(card_id) USING BTREE
);