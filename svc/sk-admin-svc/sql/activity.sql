CREATE SCHEMA IF NOT EXISTS "sec_kill";

DROP TABLE IF EXISTS sec_kill.activity;

-- ----------------------------
-- Table structure for activity
-- ----------------------------
CREATE TABLE sec_kill.activity
(
    activity_id   serial primary key not null,
    activity_name varchar(50)        not null default '',
    product_id    bigint             NOT NULL,
    start_time    bigint             NOT NULL DEFAULT 0,
    end_time      bigint             NOT NULL DEFAULT 0,
    total         bigint             NOT NULL DEFAULT 0,
    status        smallint           NOT NULL DEFAULT 0,
    sec_speed     int                NOT NULL DEFAULT 0,
    buy_limit     int                NOT NULL DEFAULT 0,
    buy_rate      real               NOT NULL DEFAULT 0.00
);

-- ----------------------------
-- Records of activity
-- ----------------------------
INSERT INTO sec_kill.activity (activity_id, activity_name, product_id, start_time, end_time, total, status, sec_speed, buy_limit, buy_rate)
VALUES (1, 'banana activity', 1, 530871061, 530871061, 20, 0, 1, 1, 0.20),
       (2, 'apple activity', 2, 530871061, 530871061, 20, 0, 1, 1, 0.20),
       (3, 'peach activity', 3, 1530928052, 1530928052, 20, 0, 1, 1, 0.20),
       (4, 'chocolate activity', 4, 1530928052, 1530928052, 20, 0, 1, 1, 0.20);