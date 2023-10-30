CREATE SCHEMA IF NOT EXISTS "sec_kill";

DROP TABLE IF EXISTS sec_kill.product;
-- ----------------------------
-- Table structure for product
-- ----------------------------
CREATE TABLE sec_kill.product
(
    product_id   serial primary key not null,
    product_name varchar(50)        not null default '',
    total        bigint             NOT NULL DEFAULT 0,
    status       smallint           NOT NULL DEFAULT 0
);

-- ----------------------------
-- Records of product
-- ----------------------------
INSERT INTO sec_kill.product
VALUES ('1', 'banana', '100', '1');
INSERT INTO sec_kill.product
VALUES ('2', 'apple', '100', '1');
INSERT INTO sec_kill.product
VALUES ('3', 'peach', '100', '1');
INSERT INTO sec_kill.product
VALUES ('4', 'chocolate', '100', '1');
