CREATE DATABASE IF NOT EXISTS `db_null_links`;
USE `db_null_links`;

DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `email` varchar(64) NOT NULL COMMENT '邮箱地址',
  `password` varchar(64) NOT NULL COMMENT '密码', 
  `avatar_url` varchar(255) NOT NULL COMMENT '头像地址',
  `background_url` varchar(255) NOT NULL COMMENT '背景地址',
  `signature` varchar(255) NOT NULL COMMENT '个性签名',
  `follow_count` int(11) NOT NULL DEFAULT '0' COMMENT '关注数',
  `follower_count` int(11) NOT NULL DEFAULT '0' COMMENT '粉丝数',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_username` (`username`),
  UNIQUE KEY `uidx_email` (`email`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '用户信息表';

DROP TABLE IF EXISTS `t_favorite`;
CREATE TABLE `t_favorite`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `webset_id` bigint(20) NOT NULL COMMENT '网页单id',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX idx_user_id_webset_id (`user_id`, `webset_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '收藏表';

DROP TABLE IF EXISTS `t_like`;
CREATE TABLE `t_like`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `webset_id` bigint(20) NOT NULL COMMENT '网页单id',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '点赞状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY uidx_user_id_webset_id (`user_id`, `webset_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '点赞表';

DROP TABLE IF EXISTS `t_relation`;
CREATE TABLE `t_relation`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `follow_id` bigint(20) NOT NULL COMMENT '关注id',
  `follower_id` bigint(20) NOT NULL COMMENT '粉丝id',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '关注表';

DROP TABLE IF EXISTS `t_chat`;
CREATE TABLE `t_chat`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `chat_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '聊天消息id',
  `user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户id',
  `webset_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '网页单id',
  `topic_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '话题id',
  `content` text NOT NULL COMMENT '消息内容',
  `type` varchar(3) NOT NULL DEFAULT '0' COMMENT '消息类型',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX idx_webset_id_created_at_status (`webset_id`, `created_at`, `status`),
  INDEX idx_topic_id_created_at_status (`topic_id`, `created_at`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '聊天表';

DROP TABLE IF EXISTS `t_topic`;
CREATE TABLE `t_topic`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `topic_title` varchar(255) NOT NULL DEFAULT '' COMMENT '话题标题',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '话题表';

DROP TABLE IF EXISTS `t_webset`;
CREATE TABLE `t_webset`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `author_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '作者id',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `cover_url` varchar(255) NOT NULL DEFAULT '' COMMENT '封面地址',
  `category` tinyint(5) NOT NULL DEFAULT 0 COMMENT '分区',
  `view_cnt` int NOT NULL DEFAULT 0 COMMENT '观看数',
  `like_cnt` int NOT NULL DEFAULT 0 COMMENT '点赞数',
  `favorite_cnt` int NOT NULL DEFAULT 0 COMMENT '收藏数',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX idx_created_at (`created_at`),
  INDEX idx_view_cnt (`view_cnt`),
  INDEX idx_like_cnt (`like_cnt`),
  INDEX idx_author_id_created_at (`author_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '网页单表';

DROP TABLE IF EXISTS `t_weblink`;
CREATE TABLE `t_weblink`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `link_id` bigint(20) NOT NULL COMMENT '网页id',
  `webset_id` bigint(20) NOT NULL COMMENT '网页单id',
  `author_id` bigint(20) NOT NULL  COMMENT '添加者id',
  `description`  varchar(255) NOT NULL COMMENT '描述',
  `url`       text NOT NULL COMMENT '网址',
  `cover_url` text NOT NULL COMMENT '封面地址',
  `click_cnt` int NOT NULL DEFAULT 0 COMMENT '点击数',
  `status`    tinyint(3) NOT NULL DEFAULT 0 COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX idx_webset_id_link_id_status (`webset_id`, `link_id`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '网页单链接表';

DROP TABLE IF EXISTS `t_balance`;
CREATE TABLE `t_balance`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(64) NOT NULL COMMENT '用户名',
  `amount` int(11) NOT NULL COMMENT '金额',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX idx_user_id_created_at (`user_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '用户余额表';

DROP TABLE IF EXISTS `t_pay_history`;
CREATE TABLE `t_pay_history`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(64) NOT NULL COMMENT '用户名',
  `amount` int(11) NOT NULL COMMENT '金额',
  `business_id` varchar(64) NOT NULL COMMENT '业务id', 
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  INDEX idx_user_id_created_at (`user_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '用户支付历史表';

DROP TABLE IF EXISTS `t_business`;
CREATE TABLE `t_user`
(
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `business` varchar(64) NOT NULL COMMENT '业务名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '业务表';
