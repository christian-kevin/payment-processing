CREATE TABLE card (
    id BIGINT AUTO_INCREMENT NOT NULL,
    wallet_id BIGINT NOT NULL,
    card_number VARCHAR(100) NOT NULL,
    expiry_date VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at BIGINT NOT NULL,
    is_deleted TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    KEY k_wallet_id_on_card(wallet_id) USING BTREE,
    KEY k_card_number_expiry_date_on_card(card_number, expiry_date) USING BTREE
);