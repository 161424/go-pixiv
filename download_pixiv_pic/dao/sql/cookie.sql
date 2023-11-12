CREATE TABLE IF NOT EXISTS `cookie`
(
    `uname`           varchar(16) PRIMARY KEY  COMMENT '主键ID',
    `app`             varchar(16) NOT NULL COMMENT '应用名',
    `status`          boolean NOT NULL COMMENT 'cookie状态 1表示过期',
    `cookie`          text NOT NULL COMMENT 'cookie',
    `deleted_at`      timestamp      DEFAULT NULL COMMENT '删除时间',
    `created_at`      timestamp      NOT NULL COMMENT '创建时间',
    `updated_at`      timestamp      DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

) ;