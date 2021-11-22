CREATE TABLE wallet (
      id BIGINT AUTO_INCREMENT NOT NULL,
      parent_id BIGINT NOT NULL,
      parent_type INT NOT NULL,
      balance BIGINT NOT NULL DEFAULT 0,
      country VARCHAR(10) NOT NULL,
      PRIMARY KEY(id),
      KEY k_parent_type_parent_id_on_wallet(parent_type, parent_id) USING BTREE
);