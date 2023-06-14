CREATE TABLE IF NOT EXISTS `auth`
(
    `uname`      varchar(16) unsigned NOT NULL  COMMENT '用户名',
    `password`   varchar(32) NOT NULL COMMENT '密码',
    `ip_1`       cidr  COMMENT 'ip1'
    `ip_2`       inet  COMMENT 'ip2'
    `Mac`        macaddr  COMMENT 'Mac'
    `output`     text DEFAULT NULL COMMENT '执行结果',
    `run_timer`  timestamp     NOT NULL COMMENT '执行时间',
    `cost_time`  int(8) DEFAULT NULL COMMENT '执行耗时',
    `status`     int(4) NOT NULL COMMENT '当前状态',
    `created_at` timestamp     NOT NULL COMMENT '创建时间',
    `updated_at` timestamp     NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`uname`) USING BTREE COMMENT '用户索引',
);