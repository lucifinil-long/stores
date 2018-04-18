DROP DATABASE IF EXISTS `stores`;
CREATE DATABASE `stores` CHARSET utf8 COLLATE utf8_general_ci;
USE `stores`;

/*
 * drop tables if exists

DROP TABLE IF EXISTS `stores_node`;
DROP TABLE IF EXISTS `stores_user`;
DROP TABLE IF EXISTS `stores_user_node`;
DROP TABLE IF EXISTS `stores_op_log`;

 */

/* 
 * create tables
 */
 /* 
 * 树形菜单节点表
 */
CREATE TABLE `stores_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '节点名称',
  `path` varchar(256) NOT NULL DEFAULT '' COMMENT '节点对应的url子路径',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '节点层级，从1开始递增',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父节点id，第一层节点的父节点均为0',
  `menu` tinyint(4) NOT NULL DEFAULT '0' COMMENT '树型菜单展示状态，0不展示，1展示',
  `auth` tinyint(4) DEFAULT '1' COMMENT '是否需要认证， 0不需要，1需要',
  `icon`  varchar(256) NULL DEFAULT NULL COMMENT '节点图标',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `idx_pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 第一级菜单 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`, `menu`)   VALUES ('1', '后台设置', '', '1', '0', '0', '1');
/* 第二级菜单 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('2', '用户管理', '/pages/admin/users', '2', '1', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('3', '系统操作日志', '/pages/admin/operations', '2', '1', '1');
/* 后台用户关联操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('4', '账号列表', '/admin/users/list', '3', '2');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('5', '新增账号', '/admin/users/add', '4', '4');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('6', '账号修改', '/admin/users/update', '4', '4');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('7', '删除账号', '/admin/users/delete', '4', '4');
/* 操作日志关联操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('8', '操作日志列表', '/admin/operations/list', '3', '3');
/* 用户权限列表操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`)           VALUES ('9', '后台用户权限列表', '/admin/user/access', '2', '1', '0');


/* 
 * 后台用户表
 */
CREATE TABLE `stores_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(128) NOT NULL DEFAULT '' COMMENT '用户登录名',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '用户密码',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `mobile` varchar(128) DEFAULT '' COMMENT '用户手机',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '用户状态，0为禁用，1为启用',
  `level` tinyint(4) NOT NULL DEFAULT '1' COMMENT '用户级别，-1为超级管理员，1为普通管理员',
  `last_login_time` datetime DEFAULT NULL,
  `created_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_un` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/* insert super administrator with password @dminPwd */
INSERT INTO `stores_user` (`username`,`password`,`nickname`,`level`, `created_time`) VALUES ('admin', '5106e9dd4f30c7a042569a4e3d42b4a4', 'super administrator', '-1', NOW());

/* 
 * 用户授权访问信息
 */
CREATE TABLE `stores_user_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `node_id` int(11) NOT NULL COMMENT '授权访问节点ID',
  PRIMARY KEY (`id`),
  KEY `idx_nid` (`node_id`),
  KEY `idx_uid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 后台操作日志表
 */
CREATE TABLE `stores_op_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '执行用户的用户ID',
  `username` varchar(128) NOT NULL COMMENT '执行用户的用户名',
  `from` varchar(128) NOT NULL COMMENT '执行用户的来源',
  `action` varchar(128) NOT NULL COMMENT '执行动作',
  `detail` text DEFAULT NULL COMMENT '操作详情',
  `created_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

