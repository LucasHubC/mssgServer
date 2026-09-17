/* ============================================================
 * mssgServer 完整数据库脚本
 * 两个库：
 *   mssg      — 业务库（玩家数据、背包、余额）
 *   game_log  — 日志库（操作日志，只写不改，量大会膨胀）
 * 四张表：
 *   mssg.users        — 用户账号
 *   mssg.balance      — 玩家余额（金币等）
 *   mssg.user_items   — 玩家背包
 *   game_log.item_op_log — 道具操作日志
 * ============================================================ */

-- ============================================================
-- 第一部分：业务库 mssg
-- ============================================================
CREATE DATABASE IF NOT EXISTS `mssg` 
  DEFAULT CHARACTER SET utf8mb4 
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `mssg`;

-- ------------------------------------------------------------
-- 表1: users 用户账号
-- ------------------------------------------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT  COMMENT '自增主键',
  `uid`          VARCHAR(64)  NOT NULL                 COMMENT '业务用户ID（UUID，对外暴露）',
  `username`     VARCHAR(64)  NOT NULL                 COMMENT '登录账号',
  `password`     VARCHAR(128) NOT NULL                 COMMENT '加密后的密码（bcrypt/sha256+盐）',
  `created_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid`      (`uid`),       -- 业务ID唯一
  UNIQUE KEY `uk_username` (`username`)   -- 登录名唯一
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户账号表';

-- ------------------------------------------------------------
-- 表2: balance 玩家余额
-- ------------------------------------------------------------
DROP TABLE IF EXISTS `balance`;
CREATE TABLE `balance` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT       NOT NULL                 COMMENT '关联 users.id',
  `coin`         BIGINT       NOT NULL DEFAULT 0       COMMENT '金币/货币数量',
  `created_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id`  (`user_id`),
  CONSTRAINT `fk_balance_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家余额表';

-- ------------------------------------------------------------
-- 表3: user_items 玩家背包
-- ------------------------------------------------------------
DROP TABLE IF EXISTS `user_items`;
CREATE TABLE `user_items` (
  `uid`          VARCHAR(64)  NOT NULL                 COMMENT '业务用户ID',
  `item_id`      VARCHAR(64)  NOT NULL                 COMMENT '道具ID（模板ID，如 sword_001）',
  `count`        BIGINT       NOT NULL DEFAULT 0       COMMENT '持有数量',
  `updated_at`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`, `item_id`),          -- 同一玩家同一道具只有一行
  KEY `idx_uid`  (`uid`)                   -- 按玩家查背包
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家背包（Redis 定时落盘）';

-- ============================================================
-- 第二部分：日志库 game_log
-- ============================================================
CREATE DATABASE IF NOT EXISTS `game_log` 
  DEFAULT CHARACTER SET utf8mb4 
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `game_log`;

-- ------------------------------------------------------------
-- 表4: item_op_log 道具操作日志
-- 注意：before / after 是 MySQL 保留字，列名改成 before_count / after_count
-- ------------------------------------------------------------
DROP TABLE IF EXISTS `item_op_log`;
CREATE TABLE `item_op_log` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT,
  `uid`           VARCHAR(64)  NOT NULL                 COMMENT '业务用户ID',
  `item_id`       VARCHAR(64)  NOT NULL                 COMMENT '道具ID',
  `op_type`       TINYINT      NOT NULL                 COMMENT '操作类型: 1=add 2=consume 3=drop',
  `count`         BIGINT       NOT NULL                 COMMENT '本次操作的数量',
  `before_count`  BIGINT       NOT NULL DEFAULT 0       COMMENT '操作前数量',
  `after_count`   BIGINT       NOT NULL DEFAULT 0       COMMENT '操作后数量',
  `source`        VARCHAR(32)  NOT NULL DEFAULT ''      COMMENT '来源: ws_cmd/http_api/kill_boss/shop/gm',
  `created_at`    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_uid_time`  (`uid`, `created_at`),   -- 按玩家查历史
  KEY `idx_item_time` (`item_id`, `created_at`), -- 按道具查变动
  KEY `idx_time`      (`created_at`)           -- 按时间清理旧分区
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='道具操作日志（只写不改，可按月分区）';

-- ============================================================
-- 第三部分：给 mssg_user 授权（阿里云 RDS 执行）
-- ============================================================
-- 如果还没创建用户，先创建：
-- CREATE USER 'mssg_user'@'%' IDENTIFIED BY '你的密码';

-- 授权 mssg 库全部权限
GRANT SELECT, INSERT, UPDATE, DELETE ON `mssg`.* TO 'mssg_user'@'%';

-- 授权 game_log 库（只需要写权限 + 查自己的日志，不需要 UPDATE）
GRANT SELECT, INSERT, DELETE ON `game_log`.* TO 'mssg_user'@'%';

FLUSH PRIVILEGES;

-- 验证
SHOW GRANTS FOR 'mssg_user'@'%';