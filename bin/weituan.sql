/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 172.16.10.177
 Source Database       : weituan

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : utf-8

 Date: 10/29/2018 15:42:43 PM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `bargain`
-- ----------------------------
DROP TABLE IF EXISTS `bargain`;
CREATE TABLE `bargain` (
  `ShareID` int(10) unsigned NOT NULL COMMENT '分享编号',
  `UserID` int(10) unsigned NOT NULL COMMENT '用户编号',
  `BargainPrice` decimal(12,2) NOT NULL COMMENT '砍价',
  `BargainTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '砍价时间',
  PRIMARY KEY (`ShareID`,`UserID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `share`
-- ----------------------------
DROP TABLE IF EXISTS `share`;
CREATE TABLE `share` (
  `ShareID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '分享编号',
  `UserID` int(10) unsigned NOT NULL COMMENT '用户编号',
  `SkuID` int(10) unsigned NOT NULL COMMENT '商品编号',
  `ShareTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分享时间',
  PRIMARY KEY (`ShareID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
--  Table structure for `sku`
-- ----------------------------
DROP TABLE IF EXISTS `sku`;
CREATE TABLE `sku` (
  `SkuID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品编号',
  `UserID` int(10) unsigned NOT NULL COMMENT '用户编号',
  `Name` varchar(255) NOT NULL COMMENT '名称',
  `Price` decimal(12,2) NOT NULL COMMENT '价格',
  `MinPrice` decimal(12,2) NOT NULL COMMENT '底价',
  `Bargain` tinyint(4) NOT NULL COMMENT '砍价（0不支持砍价 +n随机砍N次 -n等值砍N次）',
  `Intro` varchar(1020) NOT NULL COMMENT '介绍',
  `Images` varchar(1020) NOT NULL COMMENT '图片',
  `WeChatID` varchar(255) NOT NULL COMMENT '微信号（卖家）',
  `Deadline` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '截止时间',
  `PublishTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
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
