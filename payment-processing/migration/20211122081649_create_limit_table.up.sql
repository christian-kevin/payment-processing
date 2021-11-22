CREATE TABLE limits (
    id BIGINT AUTO_INCREMENT NOT NULL,
    parent_id BIGINT NOT NULL,
    parent_type INT NOT NULL,
    amount BIGINT NOT NULL,
    type INT NOT NULL,
    PRIMARY KEY(id),
    KEY k_parent_type_parent_id_on_limits(parent_type, parent_id) USING BTREE
);