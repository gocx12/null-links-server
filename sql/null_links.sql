CREATE DATABASE IF NOT EXISTS `db_null_links`;
USE `db_null_links`;

DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `email` varchar(64) NOT NULL COMMENT '邮箱地址',
  `password` varchar(36) NOT NULL COMMENT '密码', 
  `avatar_url` varchar(255) NOT NULL COMMENT '头像地址',
  `background_url` varchar(255) NOT NULL COMMENT '背景地址',
  `signature` varchar(255) NOT NULL COMMENT '个性签名',
  `follow_count` int(11) NOT NULL DEFAULT '0' COMMENT '关注数',
  `follower_count` int(11) NOT NULL DEFAULT '0' COMMENT '粉丝数',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '用户信息表';

DROP TABLE IF EXISTS `t_favorite`;
CREATE TABLE `t_favorite`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `webset_id` bigint(20) unsigned NOT NULL COMMENT '网页单id',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '收藏表';

DROP TABLE IF EXISTS `t_like`;
CREATE TABLE `t_like`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `webset_id` bigint(20) unsigned NOT NULL COMMENT '网页单id',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '点赞状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '点赞表';

DROP TABLE IF EXISTS `t_comment`;
CREATE TABLE `t_comment`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `content` varchar(255) NOT NULL COMMENT '评论内容',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '评论表';

DROP TABLE IF EXISTS `t_relation`;
CREATE TABLE `t_relation`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `follow_id` bigint(20) unsigned NOT NULL COMMENT '关注id',
  `follower_id` bigint(20) unsigned NOT NULL COMMENT '粉丝id',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '关注表';

DROP TABLE IF EXISTS `t_message`;
CREATE TABLE `t_message`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `from_user_id` bigint(20) unsigned NOT NULL COMMENT '发送者id',
  `to_user_id` bigint(20) unsigned NOT NULL COMMENT '接收者id',
  `content` varchar(255) NOT NULL COMMENT '消息内容',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '私信表';

DROP TABLE IF EXISTS `t_chat`;
CREATE TABLE `t_chat`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `webset_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '网页单id',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '消息内容',
  `type` varchar(3) NOT NULL DEFAULT '0' COMMENT '消息类型',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX idx_webset_id (`webset_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '聊天表';

DROP TABLE IF EXISTS `t_webset`;
CREATE TABLE `t_webset`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `author_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作者id',
  `describe` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `cover_url` varchar(255) NOT NULL DEFAULT '' COMMENT '封面地址',
  `category` tinyint(5) NOT NULL DEFAULT 0 COMMENT '分区',
  `view_cnt` int NOT NULL DEFAULT 0 COMMENT '观看数',
  `like_cnt` int NOT NULL DEFAULT 0 COMMENT '点赞数',
  `favorite_cnt` int NOT NULL DEFAULT 0 COMMENT '收藏数',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '网页单表';

DROP TABLE IF EXISTS `t_weblink`;
CREATE TABLE `t_weblink`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `link_id` bigint(20) unsigned NOT NULL COMMENT '网页id',
  `author_id` bigint(20) unsigned NOT NULL  COMMENT '添加者id',
  `describe`  varchar(255) NOT NULL COMMENT '描述',
  `url`       text NOT NULL COMMENT '网址',
  `cover_url` varchar(255) NOT NULL COMMENT '封面地址',
  `click_cnt` int NOT NULL DEFAULT 0 COMMENT '点击数',
  `status`    tinyint(3) NOT NULL DEFAULT '0' COMMENT '在库状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '网页单链接表';