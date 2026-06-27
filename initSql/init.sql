-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: go_chat
-- ------------------------------------------------------
-- Server version	9.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `conversations`
--

DROP TABLE IF EXISTS `conversations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `conversations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `peer_id` bigint unsigned DEFAULT NULL,
  `last_msg` longtext,
  `last_time` datetime DEFAULT NULL,
  `unread_count` bigint unsigned DEFAULT NULL,
  `type` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_ab` (`user_id`,`peer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conversations`
--

LOCK TABLES `conversations` WRITE;
/*!40000 ALTER TABLE `conversations` DISABLE KEYS */;
INSERT INTO `conversations` VALUES (77,40,41,'??','2026-06-26 21:33:50',0,0),(78,41,40,'??','2026-06-26 21:33:50',0,0),(79,42,40,'http://127.0.0.1:8080/static/upload/0abdc2cfdf6ac3a2aa92528b1f42a13b.jpg','2026-06-27 12:53:34',0,0),(80,40,42,'http://127.0.0.1:8080/static/upload/0abdc2cfdf6ac3a2aa92528b1f42a13b.jpg','2026-06-27 12:53:34',1,0),(81,42,5,'??','2026-06-27 12:53:13',0,1),(82,40,5,'??','2026-06-27 12:53:13',0,1),(83,41,5,'??','2026-06-27 12:53:13',0,1);
/*!40000 ALTER TABLE `conversations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friend_reqs`
--

DROP TABLE IF EXISTS `friend_reqs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friend_reqs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `from_id` bigint unsigned DEFAULT NULL,
  `target_id` bigint unsigned DEFAULT NULL,
  `status` bigint unsigned DEFAULT NULL,
  `create_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friend_reqs`
--

LOCK TABLES `friend_reqs` WRITE;
/*!40000 ALTER TABLE `friend_reqs` DISABLE KEYS */;
INSERT INTO `friend_reqs` VALUES (29,40,41,1,'2026-06-17 17:18:13'),(30,42,40,1,'2026-06-27 12:22:19');
/*!40000 ALTER TABLE `friend_reqs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `friend_id` bigint unsigned DEFAULT NULL,
  `create_at` datetime DEFAULT NULL,
  `status` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_friend` (`user_id`,`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (29,40,41,'2026-06-17 17:18:32',0),(30,41,40,'2026-06-17 17:18:32',0),(31,42,40,'2026-06-27 12:22:31',0),(32,40,42,'2026-06-27 12:22:31',0);
/*!40000 ALTER TABLE `friends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_members`
--

DROP TABLE IF EXISTS `group_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group_members` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `role` bigint DEFAULT '0',
  `is_mute` bigint DEFAULT '0',
  `mute_end_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_members`
--

LOCK TABLES `group_members` WRITE;
/*!40000 ALTER TABLE `group_members` DISABLE KEYS */;
INSERT INTO `group_members` VALUES (23,5,42,2,0,NULL,'2026-06-27 12:24:37.353'),(24,5,40,0,0,NULL,'2026-06-27 12:43:11.046'),(25,5,41,0,0,NULL,'2026-06-27 12:43:31.809');
/*!40000 ALTER TABLE `group_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_models`
--

DROP TABLE IF EXISTS `group_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group_models` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(64) NOT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `total_count` int DEFAULT NULL,
  `notice` varchar(1000) DEFAULT NULL,
  `owner_id` bigint unsigned NOT NULL,
  `is_all_mute` bigint DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_models`
--

LOCK TABLES `group_models` WRITE;
/*!40000 ALTER TABLE `group_models` DISABLE KEYS */;
INSERT INTO `group_models` VALUES (5,'三国演义','',3,'',42,0,'2026-06-27 12:24:37','2026-06-27 12:24:37');
/*!40000 ALTER TABLE `group_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `msg_id` varchar(255) NOT NULL,
  `from_id` bigint unsigned DEFAULT NULL,
  `target_id` bigint unsigned DEFAULT NULL,
  `type` longtext,
  `media` bigint DEFAULT NULL,
  `content` longtext,
  `picture` longtext,
  `url` longtext,
  `create_at` datetime DEFAULT NULL,
  `msg_type` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `msg_id_UNIQUE` (`msg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=280 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (267,'ae36b477-13ae-49f0-b06a-3a4d6f6ee43d',40,41,'chat',0,'666',NULL,NULL,'2026-06-26 17:41:09',0),(268,'9b83f811-8154-496f-8806-1f6c33c20ce6',41,40,'chat',0,'hhh',NULL,NULL,'2026-06-26 17:41:33',0),(269,'fcb237e1-db29-4ec8-b653-c811854ab0f3',40,41,'chat',0,'你好',NULL,NULL,'2026-06-26 17:45:24',0),(270,'774a3b08-cb66-47e2-85f6-5986de36ec65',41,40,'chat',0,'??',NULL,NULL,'2026-06-26 21:33:50',0),(274,'bde02b06-3e25-4d84-8850-d4512fe90ee3',42,40,'chat',0,'你好',NULL,NULL,'2026-06-27 12:22:59',0),(275,'b656cf6f-23be-47c2-a619-b829b566626b',41,5,'groupMessage',0,'?',NULL,NULL,'2026-06-27 12:43:55',0),(276,'6b5b412e-a320-4438-b17b-b1a8ce6a9403',42,5,'groupMessage',0,'http://127.0.0.1:8080/static/upload/7e7a45d37d46db56a2a1a1c1bfecd83e.jpg',NULL,NULL,'2026-06-27 12:44:14',1),(277,'e4f15e4d-a41a-4810-867f-8f2e612dca16',40,5,'groupMessage',0,'http://127.0.0.1:8080/static/upload/40b29ca738c42d0f62c4ec9e93207b1a.jpg',NULL,NULL,'2026-06-27 12:44:41',1),(278,'098fcf8e-9582-4f21-959c-a15860e2c864',40,5,'groupMessage',0,'??',NULL,NULL,'2026-06-27 12:53:13',0),(279,'96627f31-3d0d-4f1c-bb2b-13dc067db791',42,40,'chat',0,'http://127.0.0.1:8080/static/upload/0abdc2cfdf6ac3a2aa92528b1f42a13b.jpg',NULL,NULL,'2026-06-27 12:53:34',1);
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_basic`
--

DROP TABLE IF EXISTS `user_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_basic` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(20) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `email` varchar(20) DEFAULT NULL,
  `phone` varchar(11) DEFAULT NULL,
  `avatar` varchar(100) DEFAULT NULL,
  `client_ip` longtext,
  `client_port` longtext,
  `login_time` datetime DEFAULT NULL,
  `heartbeat_time` datetime DEFAULT NULL,
  `login_out_time` datetime DEFAULT NULL,
  `is_logout` tinyint(1) DEFAULT NULL,
  `device_info` longtext,
  `create_at` datetime DEFAULT NULL,
  `update_at` datetime DEFAULT NULL,
  `delete_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_basic`
--

LOCK TABLES `user_basic` WRITE;
/*!40000 ALTER TABLE `user_basic` DISABLE KEYS */;
INSERT INTO `user_basic` VALUES (40,'关羽','$2a$11$Ge48h5zi4PkKwh58XKsHp.O0aiR0aRML/qXnLa.5R6MiJQMqM3ICm','','','','','',NULL,NULL,NULL,0,'','2026-06-17 17:14:02','2026-06-17 17:14:02',NULL),(41,'张飞','$2a$11$nE0b5iiC0U6pO83Uk3WIFOS0lptn2B/aNKPi3wO3eBLWB6TWGO11W','','','','','',NULL,NULL,NULL,0,'','2026-06-17 17:17:42','2026-06-17 17:17:42',NULL),(42,'孙尚香','$2a$11$f0Nq.5WsucJy6CvD47eZnOT39K3I4ofl8CA4sLogdv/IUzfyFUgY.','','','?','','',NULL,NULL,NULL,0,'','2026-06-27 12:20:57','2026-06-27 12:20:57',NULL);
/*!40000 ALTER TABLE `user_basic` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-27 17:36:42
