/*
 Navicat Premium Data Transfer

 Source Server         : jiaen
 Source Server Type    : MariaDB
 Source Server Version : 100504
 Source Host           : localhost:3306
 Source Schema         : chatdb

 Target Server Type    : MariaDB
 Target Server Version : 100504
 File Encoding         : 65001

 Date: 15/07/2020 19:48:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for conversation
-- ----------------------------
DROP TABLE IF EXISTS `conversation`;
CREATE TABLE `conversation`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `userlist` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户会话',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '会话表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of conversation
-- ----------------------------
INSERT INTO `conversation` VALUES (1, ',1,2,');
INSERT INTO `conversation` VALUES (2, ',1,3,');

-- ----------------------------
-- Table structure for conversation_user
-- ----------------------------
DROP TABLE IF EXISTS `conversation_user`;
CREATE TABLE `conversation_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `cvsid` int(11) NULL DEFAULT NULL COMMENT '对应会话id',
  `userid` int(11) NULL DEFAULT NULL COMMENT '关联的用户id',
  `lastid` int(11) NULL DEFAULT NULL COMMENT '最后阅读消息的id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `index_cvsid`(`cvsid`) USING BTREE COMMENT '会话索引',
  INDEX `index_userid`(`userid`) USING BTREE COMMENT '用户索引'
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户对应的会话表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of conversation_user
-- ----------------------------
INSERT INTO `conversation_user` VALUES (1, 1, 1, 3);
INSERT INTO `conversation_user` VALUES (2, 1, 2, 2);
INSERT INTO `conversation_user` VALUES (3, 2, 3, 1);
INSERT INTO `conversation_user` VALUES (4, 2, 1, 1);

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `cvsid` int(11) NULL DEFAULT NULL COMMENT '会话id',
  `source` int(255) NULL DEFAULT NULL COMMENT '消息发送者ID',
  `content` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '消息内容',
  `sendtime` int(11) NULL DEFAULT NULL COMMENT '消息发送时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `index_cvsid`(`cvsid`) USING BTREE COMMENT '会话索引'
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of message
-- ----------------------------
INSERT INTO `message` VALUES (1, 1, 2, '在不在不', 1594652539);
INSERT INTO `message` VALUES (2, 1, 1, '我不在', 1594652639);
INSERT INTO `message` VALUES (3, 1, 2, '那你在哪呢！', 1594655639);
INSERT INTO `message` VALUES (4, 2, 2, '小宝贝', 1594755639);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户姓名',
  `usertel` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '手机号',
  `img` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户头像',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '好家伙', '10000000000', 'https://ss2.bdstatic.com/70cFvnSh_Q1YnxGkpoWK1HF6hhy/it/u=1399335390,2996130648&fm=26&gp=0.jpg');
INSERT INTO `user` VALUES (2, '老李头', '12222222222', 'https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=2599028949,1720538774&fm=26&gp=0.jpg');
INSERT INTO `user` VALUES (3, 'vvv', '13333333333', 'https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1594812971108&di=a80febf3eed4020912b284b9bf2688a3&imgtype=0&src=http%3A%2F%2Fimgsrc.baidu.com%2Fforum%2Fw%3D580%2Fsign%3D446bf26975f40ad115e4c7eb672c1151%2Fe6ac922a2834349b4d2babd4c7ea15ce36d3be30.jpg');

SET FOREIGN_KEY_CHECKS = 1;
