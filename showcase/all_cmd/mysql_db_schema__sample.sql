-- --------------------------------------------------------
-- 호스트:                          127.0.0.1
-- 서버 버전:                        10.6.4-MariaDB - mariadb.org binary distribution
-- 서버 OS:                        Win64
-- HeidiSQL 버전:                  11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- dev_bp_sample 데이터베이스 구조 내보내기
CREATE DATABASE IF NOT EXISTS `dev_bp_sample` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */;
USE `dev_bp_sample`;

-- 테이블 dev_bp_sample.test_table 구조 내보내기
CREATE TABLE IF NOT EXISTS `test_table` (
  `u8_seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `is_bool` tinyint(4) DEFAULT 0,
  `s_str` varchar(100) COLLATE utf8mb3_bin DEFAULT '',
  `bt_bin` varbinary(128) DEFAULT '',
  `n1_num` tinyint(4) DEFAULT 0,
  `n2_num` smallint(6) DEFAULT 0,
  `n4_num` int(11) DEFAULT 0,
  `n8_num` bigint(20) DEFAULT 0,
  `u1_num` tinyint(3) unsigned DEFAULT 0,
  `u2_num` smallint(5) unsigned DEFAULT 0,
  `u4_num` int(10) unsigned DEFAULT 0,
  `u8_num` bigint(20) unsigned DEFAULT 0,
  `f4_num` float DEFAULT 0,
  `f8_num` double DEFAULT 0,
  `bt_snum` varbinary(66) DEFAULT '',
  PRIMARY KEY (`u8_seq`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin ROW_FORMAT=DYNAMIC;

-- 내보낼 데이터가 선택되어 있지 않습니다.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
