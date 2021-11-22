CREATE DATABASE IF NOT EXISTS `app`;
USE `app`;

CREATE TABLE user (
  id BIGINT AUTO_INCREMENT NOT NULL,
  username VARCHAR(100) NOT NULL,
  password VARCHAR(200) NOT NULL,
  country VARCHAR(10) NOT NULL,
  PRIMARY KEY(id),
  KEY k_username_on_user(username) USING BTREE
);

CREATE TABLE wallet (
  id BIGINT AUTO_INCREMENT NOT NULL,
  parent_id BIGINT NOT NULL,
  parent_type INT NOT NULL,
  balance BIGINT NOT NULL DEFAULT 0,
  country VARCHAR(10) NOT NULL,
  PRIMARY KEY(id),
  KEY k_parent_type_parent_id_on_wallet(parent_type, parent_id) USING BTREE
);

CREATE TABLE card (
  id BIGINT AUTO_INCREMENT NOT NULL,
  wallet_id BIGINT NOT NULL,
  card_number VARCHAR(100) NOT NULL,
  expiry_date VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY(id),
  KEY k_wallet_id_on_card(wallet_id) USING BTREE
);

CREATE TABLE card_transaction_log (
  id BIGINT AUTO_INCREMENT NOT NULL,
  card_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY(id),
  KEY k_card_id_on_card_transaction_log(card_id) USING BTREE
);

CREATE TABLE wallet_balance_log (
  id BIGINT AUTO_INCREMENT NOT NULL,
  wallet_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY(id),
  KEY k_wallet_id_on_wallet_balance_log(wallet_id) USING BTREE
);

CREATE TABLE limits (
  id BIGINT AUTO_INCREMENT NOT NULL,
  parent_id BIGINT NOT NULL,
  parent_type INT NOT NULL,
  amount BIGINT NOT NULL,
  type INT NOT NULL,
  PRIMARY KEY(id),
  KEY k_parent_type_parent_id_on_limits(parent_type, parent_id) USING BTREE
);