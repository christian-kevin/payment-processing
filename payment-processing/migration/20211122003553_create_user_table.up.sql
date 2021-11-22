CREATE TABLE user (
     id BIGINT AUTO_INCREMENT NOT NULL,
     username VARCHAR(100) NOT NULL,
     password VARCHAR(200) NOT NULL,
     country VARCHAR(10) NOT NULL,
     PRIMARY KEY(id),
     KEY k_username_on_user(username) USING BTREE
);