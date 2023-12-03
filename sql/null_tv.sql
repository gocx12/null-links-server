CREATE DATABASE IF NOT EXISTS `null_tv`;
USE `null_tv`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
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
  `is_follow` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否关注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `author_id` bigint(20) unsigned NOT NULL COMMENT '作者id',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `play_url` varchar(255) NOT NULL COMMENT '播放地址',
  `cover_url` varchar(255) NOT NULL COMMENT '封面地址',
  `play_cnt` int(11) NOT NULL COMMENT '播放量',
  `like_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `favorite_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
  `comment_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `is_favorite` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否收藏',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
  `content` varchar(255) NOT NULL COMMENT '评论内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `follow_id` bigint(20) unsigned NOT NULL COMMENT '关注id',
  `follower_id` bigint(20) unsigned NOT NULL COMMENT '粉丝id',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `from_user_id` bigint(20) unsigned NOT NULL COMMENT '发送者id',
  `to_user_id` bigint(20) unsigned NOT NULL COMMENT '接收者id',
  `content` varchar(255) NOT NULL COMMENT '消息内容',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;