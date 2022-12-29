-- --------------------------------------------------------
-- 호스트:                          127.0.0.1
-- 서버 버전:                        10.5.8-MariaDB - mariadb.org binary distribution
-- 서버 OS:                        Win64
-- HeidiSQL 버전:                  11.2.0.6213
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- dev_bp_test 데이터베이스 구조 내보내기
CREATE DATABASE IF NOT EXISTS `dev_bp_test` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `dev_bp_test`;

-- 테이블 dev_bp_test.bp_gen_test_schema 구조 내보내기
CREATE TABLE IF NOT EXISTS `bp_gen_test_schema` (
  `u8_seq` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `is_bool` tinyint(4) DEFAULT 0,
  `s_str` varchar(100) COLLATE utf8_bin DEFAULT '',
  `n1_num` tinyint(4) DEFAULT 0,
  `n2_num` smallint(6) DEFAULT 0,
  `n4_num` int(11) DEFAULT 0,
  `n8_num` bigint(20) DEFAULT 0,
  `u1_num` tinyint(3) unsigned DEFAULT 0,
  `u2_num` smallint(5) unsigned DEFAULT 0,
  `u4_num` int(10) unsigned DEFAULT 0,
  `u8_num` bigint(20) unsigned DEFAULT 0,
  `bt_bin` varbinary(128) DEFAULT '',
  `sn_snum` longblob DEFAULT '',
  `f_num` float DEFAULT 0,
  `d_num` double DEFAULT 0,
  `js_str` longtext COLLATE utf8_bin DEFAULT '',
  `dt_time` int(11) unsigned DEFAULT 0,
  `dtn_time` bigint(20) unsigned DEFAULT 0,
  PRIMARY KEY (`u8_seq`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC;

-- 테이블 데이터 dev_bp_test.bp_gen_test_schema:~10 rows (대략적) 내보내기
DELETE FROM `bp_gen_test_schema`;
/*!40000 ALTER TABLE `bp_gen_test_schema` DISABLE KEYS */;
INSERT INTO `bp_gen_test_schema` (`u8_seq`, `is_bool`, `s_str`, `n1_num`, `n2_num`, `n4_num`, `n8_num`, `u1_num`, `u2_num`, `u4_num`, `u8_num`, `bt_bin`, `sn_snum`, `f_num`, `d_num`, `js_str`, `dt_time`, `dtn_time`) VALUES
	(1, 0, 'a', 127, 32767, 2147483647, 9223372036854775807, NULL, NULL, NULL, 0, _binary 0x61, _binary 0x88056BC75E2D63100000, 0, 0, 'a', 0, 1534435),
	(2, 1, 'b', -127, -32767, -2147483647, -9223372036854775807, NULL, NULL, NULL, 0, _binary 0x62, _binary 0x77FA9438A1D29CEFFFFF, 1, 1, 'b', 1, 15),
	(3, 0, 'c', 3, 3, 3, 3, NULL, NULL, NULL, 0, _binary 0x63, _binary 0x8B033B2E373409220F84F00000, 0.00001, 0.00001, 'c', 2, 153),
	(4, 1, 'd', 4, 4, 4, 4, NULL, NULL, NULL, 0, _binary 0x64, _binary 0x74FCC4D1C8CBF6DDF07B0FFFFF, -0.00001, -0.00001, 'd', 3, 1534),
	(5, 0, 'e', 5, 5, 5, 5, NULL, NULL, NULL, 0, _binary 0x65, _binary 0x8000, 94194900, 94194939.12314345, 'e', 4, 15344),
	(6, 0, 'a', 127, 32767, 2147483647, 9223372036854775807, NULL, NULL, NULL, 0, _binary 0x61, _binary 0x88056BC75E2D63100000, 0, 0, 'a', 0, 1534435),
	(7, 1, 'b', -127, -32767, -2147483647, -9223372036854775807, NULL, NULL, NULL, 0, _binary 0x62, _binary 0x77FA9438A1D29CEFFFFF, 1, 1, 'b', 1, 15),
	(8, 0, 'c', 3, 3, 3, 3, NULL, NULL, NULL, 0, _binary 0x63, _binary 0x8B033B2E373409220F84F00000, 0.00001, 0.00001, 'c', 2, 153),
	(9, 1, 'd', 4, 4, 4, 4, NULL, NULL, NULL, 0, _binary 0x64, _binary 0x74FCC4D1C8CBF6DDF07B0FFFFF, -0.00001, -0.00001, 'd', 3, 1534),
	(10, 0, 'e', 5, 5, 5, 5, NULL, NULL, NULL, 0, _binary 0x65, _binary 0x8000, 94194900, 94194939.12314345, 'e', 4, 15344);
/*!40000 ALTER TABLE `bp_gen_test_schema` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
