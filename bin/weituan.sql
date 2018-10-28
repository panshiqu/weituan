/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50718
 Source Host           : localhost
 Source Database       : weituan

 Target Server Type    : MySQL
 Target Server Version : 50718
 File Encoding         : utf-8

 Date: 10/28/2018 23:22:47 PM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `sku`
-- ----------------------------
DROP TABLE IF EXISTS `sku`;
CREATE TABLE `sku` (
  `SkuID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品编号',
  `Name` varchar(255) NOT NULL COMMENT '名称',
  `Price` decimal(10,2) NOT NULL COMMENT '价格',
  `MinPrice` decimal(10,2) NOT NULL COMMENT '底价',
  `Bargain` tinyint(4) NOT NULL COMMENT '砍价（0不支持砍价 +n随机砍N次 -n等值砍N次）',
  `Intro` varchar(1024) NOT NULL COMMENT '介绍',
  `Images` varchar(255) NOT NULL COMMENT '图片',
  `WeChatID` varchar(255) NOT NULL COMMENT '微信号（卖家）',
  `Deadline` timestamp NULL DEFAULT NULL COMMENT '截止时间',
  `PublishTimestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  PRIMARY KEY (`SkuID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `UserID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户编号',
  `OpenID` varchar(255) NOT NULL COMMENT '用户唯一标识',
  PRIMARY KEY (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
