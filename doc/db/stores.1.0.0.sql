DROP DATABASE IF EXISTS `stores`;
CREATE DATABASE `stores` CHARSET utf8 COLLATE utf8_general_ci;
USE `stores`;

/*
 * drop tables if exists

DROP TABLE IF EXISTS `stores_node`;
DROP TABLE IF EXISTS `stores_user`;
DROP TABLE IF EXISTS `stores_op_log`;
DROP TABLE IF EXISTS `stores_role`;
DROP TABLE IF EXISTS `stores_role_node`;
DROP TABLE IF EXISTS `stores_role_user`;

 */

/* 
 * create tables
 */
/* 
 * 树形菜单节点表 
 */
CREATE TABLE `stores_node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '节点名称',
  `path` varchar(256) NOT NULL DEFAULT '' COMMENT '节点对应的url子路径',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '节点层级，从1开始递增',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '父节点id，第一层节点的父节点均为0',
  `menu` tinyint(4) NOT NULL DEFAULT '0' COMMENT '树型菜单展示状态，0不展示，1展示',
  `auth` tinyint(4) DEFAULT '1' COMMENT '是否需要认证， 0不需要，1需要',
  `icon`  varchar(256) NULL DEFAULT NULL COMMENT '节点图标',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `idx_pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 第一级菜单 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`, `menu`)   VALUES ('1', '商品', '', '1', '0', '0', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`, `menu`)   VALUES ('2', '统计', '', '1', '0', '0', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`, `menu`)   VALUES ('9', '后台设置', '', '1', '0', '0', '1');
/* 第二级菜单 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('10', '商品采购入库', '/pages/commodity/warehousing', '2', '1', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('11', '商品销售出库', '/pages/commodity/sale', '2', '1', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('12', '商品库存', '/pages/commodity/stores', '2', '1', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('20', '月报', '/pages/statistics/stores', '2', '2', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('21', '自定义报表', '/pages/statistics/custom', '2', '2', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('90', '商品设置', '/pages/admin/commodities', '2', '9', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('91', '库房设置', '/pages/admin/locations', '2', '9', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('92', '商品规格设置', '/pages/admin/specification', '2', '9', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('97', '用户管理', '/pages/admin/users', '2', '9', '1');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `menu`)           VALUES ('98', '系统操作日志', '/pages/admin/operations', '2', '9', '1');
/* 用户权限列表操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`, `auth`)           VALUES ('99', '后台用户权限列表', '/admin/user/access', '2', '9', '0');
/* 商品关联操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('100', '商品入库列表', '/commodity/purchases/list', '3', '10');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('101', '新增商品入库', '/commodity/purchases/add', '4', '100');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('102', '商品入库修改', '/commodity/purchases/update', '4', '100');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('103', '删除商品入库', '/commodity/purchases/delete', '4', '100');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('110', '商品销售列表', '/commodity/sales/list', '3', '11');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('111', '新增商品销售', '/commodity/sales/add', '4', '110');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('112', '修改商品销售', '/commodity/sales/update', '4', '110');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('113', '商品销售退货', '/commodity/sales/cancel', '4', '110');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('120', '商品库存列表', '/commodity/stores/list', '3', '12');
/* 商品统计操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('200', '月报列表', '/commodity/statistics/list', '3', '20');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('201', '核算月报', '/commodity/statistics/addorupdate', '4', '200');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('202', '月报设置', '/commodity/update/setting', '4', '100');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('210', '自定义报表核算', '/commodity/sales/list', '3', '11');
/* 商品设置操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('900', '商品列表', '/admin/commodities/list', '3', '90');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('901', '新增商品', '/admin/commodities/add', '4', '900');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('902', '商品修改', '/admin/commodities/update', '4', '900');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('903', '删除商品', '/admin/commodities/delete', '4', '900');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('904', '商品规格列表', '/admin/commodity/specifications/list', '3', '90');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('905', '新增商品规格', '/admin/commodity/specifications/add', '4', '904');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('906', '修改商品规格', '/admin/commodity/specifications/update', '4', '904');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('907', '删除商品规格', '/admin/commodity/specifications/delete', '4', '904');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('9000', '商品SKU列表', '/admin/commodities/sku/list', '4', '900');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('9001', '新增商品SKU', '/admin/commodities/sku/add', '5', '9000');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('9002', '商品SKU修改', '/admin/commodities/sku/update', '5', '9000');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('9003', '删除商品SKU', '/admin/commodities//sku/delete', '5', '9000');
/* 库房设置操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('910', '库房列表', '/admin/depot/list', '3', '91');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('911', '新增库房', '/admin/depot/add', '4', '910');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('912', '库房修改', '/admin/depot/update', '4', '910');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('913', '删除库房', '/admin/depot/delete', '4', '910');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('914', '货架列表', '/admin/depot/shelf/list', '5', '910');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('915', '新增货架', '/admin/depot/shelf/add', '5', '912');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('916', '修改货架', '/admin/depot/shelf/update', '5', '912');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('917', '删除货架', '/admin/depot/shelf/delete', '5', '912');
/* 后台用户关联操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('970', '账号列表', '/admin/users/list', '3', '97');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('971', '新增账号', '/admin/users/add', '4', '970');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('972', '账号修改', '/admin/users/update', '4', '970');
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('973', '删除账号', '/admin/users/delete', '4', '970');
/* 操作日志关联操作 */
INSERT INTO `stores_node` (`id`, `title`, `path`, `level`, `pid`)                   VALUES ('980', '操作日志列表', '/admin/operations/list', '3', '98');

/* 
 * 后台用户表
 */
CREATE TABLE `stores_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mobile` bigint(20) NOT NULL COMMENT '用户手机',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '用户密码',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `deletable` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否可删除',
  `deleted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已删除',
  `last_login_time` datetime DEFAULT NULL,
  `created_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_mobi` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1000;
/* insert super administrator with password @dminPwd */
INSERT INTO `stores_user` (`id`, `mobile`, `nickname`,`password`,`deletable`,`created_time`) VALUES (-1, 12345678901, 'admin', '5106e9dd4f30c7a042569a4e3d42b4a4', 0, NOW());

/* 
 * 后台操作日志表
 */
CREATE TABLE `stores_op_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '执行用户的用户ID',
  `nickname` varchar(128) NOT NULL COMMENT '执行用户的用户名',
  `from` varchar(128) NOT NULL COMMENT '执行用户的来源',
  `action` varchar(128) NOT NULL COMMENT '执行动作',
  `detail` text DEFAULT NULL COMMENT '操作详情',
  `created_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 用户角色授权访问信息
 */
CREATE TABLE `stores_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(128) NOT NULL COMMENT '角色类型',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `deletable` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否可删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `stores_role` (`role_name`, `remark`, `deletable`, `id`) VALUES ('超级管理员', '后台超级管理员', 0, -1);
INSERT INTO `stores_role` (`role_name`, `remark`) VALUES ('管理员', '后台管理员');
INSERT INTO `stores_role` (`role_name`, `remark`) VALUES ('库管', '库管员工');
INSERT INTO `stores_role` (`role_name`, `remark`) VALUES ('销售', '销售员工');

/* 
 * 用户角色授权访问信息
 */
CREATE TABLE `stores_role_node` (
  `role_id` bigint(20) NOT NULL COMMENT '用户角色类型ID',
  `node_id` bigint(20) NOT NULL COMMENT '授权访问节点ID',
  PRIMARY KEY (`role_id`, `node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 用户角色授权访问信息
 */
CREATE TABLE `stores_role_user` (
  `role_id` bigint(20) NOT NULL COMMENT '用户角色类型ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`role_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `stores_role_user` (`role_id`, `user_id`) VALUES (-1, -1);

/* 
 * 商品仓库
 */
CREATE TABLE `stores_location_depot` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '仓库名称',
  `detail` varchar(64) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_un` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品货架
 */
CREATE TABLE `stores_location_shelf` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `depot_id` bigint(20) NOT NULL COMMENT '仓库ID',
  `name` varchar(64) NOT NULL COMMENT '货架名称',
  `layers` tinyint(4) NOT NULL DEFAULT 1 COMMENT '层数',
  `detail` varchar(64) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uni` (`depot_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品信息
 */
CREATE TABLE `stores_commodity` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `commodity_name` varchar(128) NOT NULL COMMENT '商品名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品规格
 */
CREATE TABLE `stores_commodity_spec` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `spec_name` varchar(64) NOT NULL COMMENT '规格名称',
  `detail` varchar(128) NOT NULL COMMENT '描述',
  `segmentable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否可拆分',
  `segment_id` bigint(20) NULL COMMENT '可拆分的下级ID',
  `segment_amount` int(11) NULL COMMENT '可拆分的下级规格数量',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品SKU
 */
CREATE TABLE `stores_commodity_sku` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `commodity_id` bigint(20) NOT NULL COMMENT '所属商品ID',
  `sku_name` varchar(64) NOT NULL COMMENT 'SKU名称',
  `barcode` varchar(64) NOT NULL DEFAULT '' COMMENT '条码',
  `spec_id` bigint(20) NOT NULL COMMENT '关联规格ID',
  `profit` int(11) NOT NULL DEFAULT 20 COMMENT '利润率，基准销售价为 成本价 * (100 + profit)',
  `discount` tinyint(4) NULL DEFAULT 0 COMMENT '在基准销售价基础上默认折扣幅度, 0~100',
  `max_discount` tinyint(4) NULL DEFAULT 0 COMMENT '在基准销售价基础上允许的最大折扣幅度, 0~100',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_bc` (`barcode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品SKU属性表
 */
CREATE TABLE `stores_sku_property` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `property` varchar(64) NOT NULL COMMENT 'SKU属性名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品SKU属性值表
 */
CREATE TABLE `stores_sku_property_value` (
  `sku_id` bigint(20) NOT NULL,
  `property_id` bigint(20) NOT NULL COMMENT 'SKU属性ID',
  `value` varchar(64) NOT NULL COMMENT 'SKU属性值',
  PRIMARY KEY (`sku_id`, `property_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品库存表
 */
CREATE TABLE `stores_sku_stock` (
  `sku_id` bigint(20) NOT NULL COMMENT 'SKU ID',
  `shelf_id` bigint(20) NOT NULL COMMENT '货架',
  `price`  double(10,2) NOT NULL COMMENT '存货进价平均数',
  `amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '库存数量',
  PRIMARY KEY (`sku_id`, `shelf_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 
 * 商品库存表变动表
 */
CREATE TABLE `stores_sku_stock_change` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sku_id` bigint(20) NOT NULL COMMENT 'SKU ID',
  `shelf_id` bigint(20) NOT NULL COMMENT '货架',
  `layer` varchar(64) NULL COMMENT '货架层数',
  `price` double(10,2) NOT NULL COMMENT '入库、拆包入库时为入库价格；销售出库时为售价；拆包出库时为存货进价平均数',
  `amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '库存数量',
  `type` tinyint NOT NULL COMMENT '变动类型：1采购入库；2销售出库；3拆包出库；4拆包入库',
  `operator_id` bigint(20) NOT NULL COMMENT '操作人',
  `created_time` datetime NOT NULL,
  `detail` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`),
  KEY `idx_ct` (`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;