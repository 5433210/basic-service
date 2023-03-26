/*
 Navicat Premium Data Transfer
 
 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : localhost:3306
 Source Schema         : authn
 
 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001
 
 Date: 18/10/2021 19:07:28
 */
SET
    NAMES utf8mb4;

SET
    FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE authn;

USE authn;

-- ----------------------------
-- Table structure for tb_credential
-- ----------------------------
DROP TABLE IF EXISTS `tb_credential`;

CREATE TABLE `tb_credential` (
    `id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `config` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `identity_id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `credential_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `credential_idx` (`identity_id`),
    KEY `fk_credential_credential_type_id` (`credential_type`),
    CONSTRAINT `fk_credential_credential_type_id` FOREIGN KEY (`credential_type`) REFERENCES `tb_credential_type` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_credential_identity_id` FOREIGN KEY (`identity_id`) REFERENCES `tb_identity` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_credential
-- ----------------------------
BEGIN;

COMMIT;

-- ----------------------------
-- Table structure for tb_credential_type
-- ----------------------------
DROP TABLE IF EXISTS `tb_credential_type`;

CREATE TABLE `tb_credential_type` (
    `id` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `descr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_credential_type
-- ----------------------------
BEGIN;

INSERT INTO
    `tb_credential_type`
VALUES
    (
        'P',
        'password as credential',
        '2021-10-18 18:30:28',
        '2021-10-18 18:30:31'
    );

COMMIT;

-- ----------------------------
-- Table structure for tb_domain
-- ----------------------------
DROP TABLE IF EXISTS `tb_domain`;

CREATE TABLE `tb_domain` (
    `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `descr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_domain
-- ----------------------------
BEGIN;

INSERT INTO
    `tb_domain`
VALUES
    (
        'wailik.com',
        'a domain example',
        '2021-10-18 18:26:58',
        '2021-10-18 18:27:04'
    );

COMMIT;

-- ----------------------------
-- Table structure for tb_identifier
-- ----------------------------
DROP TABLE IF EXISTS `tb_identifier`;

CREATE TABLE `tb_identifier` (
    `id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `identifier` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `identitiy_id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `domain_id` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `identifier_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `identifiers_idx` (`identifier`, `domain_id`, `identifier_type`) USING BTREE,
    KEY `fk_identifier_domain_id` (`domain_id`),
    KEY `fk_identifier_identifier_type` (`identifier_type`),
    CONSTRAINT `fk_identifier_domain_id` FOREIGN KEY (`domain_id`) REFERENCES `tb_domain` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_identifier_identifier_type` FOREIGN KEY (`identifier_type`) REFERENCES `tb_identifier_type` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_identifier_identity_id` FOREIGN KEY (`identitiy_id`) REFERENCES `tb_identity` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_identifier
-- ----------------------------
BEGIN;

COMMIT;

-- ----------------------------
-- Table structure for tb_identifier_credential
-- ----------------------------
DROP TABLE IF EXISTS `tb_identifier_credential`;

CREATE TABLE `tb_identifier_credential` (
    `identifier_id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `credential_id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`identifier_id`),
    KEY `fk_identifier_credential_credential_id` (`credential_id`),
    CONSTRAINT `fk_identifier_credential_credential_id` FOREIGN KEY (`credential_id`) REFERENCES `tb_credential` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_identifier_credential_identifier_id` FOREIGN KEY (`identifier_id`) REFERENCES `tb_identifier` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_identifier_credential
-- ----------------------------
BEGIN;

COMMIT;

-- ----------------------------
-- Table structure for tb_identifier_type
-- ----------------------------
DROP TABLE IF EXISTS `tb_identifier_type`;

CREATE TABLE `tb_identifier_type` (
    `id` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `is_verifiable` tinyint(1) NOT NULL,
    `descr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_identifier_type
-- ----------------------------
BEGIN;

INSERT INTO
    `tb_identifier_type`
VALUES
    (
        'E',
        1,
        'email as an identifier',
        '2021-10-18 18:28:52',
        '2021-10-18 18:28:55'
    );

INSERT INTO
    `tb_identifier_type`
VALUES
    (
        'M',
        1,
        'mobile phone as an identifier',
        '2021-10-18 18:29:25',
        '2021-10-18 18:29:28'
    );

INSERT INTO
    `tb_identifier_type`
VALUES
    (
        'A',
        1,
        'account as an identifier',
        '2021-10-18 18:29:25',
        '2021-10-18 18:29:28'
    );

COMMIT;

-- ----------------------------
-- Table structure for tb_identity
-- ----------------------------
DROP TABLE IF EXISTS `tb_identity`;

CREATE TABLE `tb_identity` (
    `id` char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `domain_id` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `stat` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `stat_changed_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `fk_identify_domain_id` (`domain_id`),
    CONSTRAINT `fk_identify_domain_id` FOREIGN KEY (`domain_id`) REFERENCES `tb_domain` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

-- ----------------------------
-- Records of tb_identity
-- ----------------------------
BEGIN;

COMMIT;

SET
    FOREIGN_KEY_CHECKS = 1;