CREATE DATABASE IF NOT EXISTS `dpe_insights` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
USE `dpe_insights`;
-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Database: dpe_insights
-- ------------------------------------------------------
-- Server version	5.6.50

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
-- Table structure for table `extracted_pull_request_count`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `extracted_pull_request_count` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `open` int(11) DEFAULT 0,
    `closed` int(11) DEFAULT 0,
    `total` int(11) DEFAULT 0,
    `date` datetime DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `date_UNIQUE` (`date`),
    KEY `time` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `extracted_pull_requests`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_pull_requests` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `pr_id` int(11) NOT NULL,
     `pr_number` int(11) NOT NULL COMMENT 'Pull request number',
     `comments` int(11) DEFAULT 0,
     `branch_name` VARCHAR(200) NOT NULL,
     `review_comments` int(11) DEFAULT 0,
     `commits` int(11) DEFAULT 0,
     `additions` int(11) DEFAULT 0,
     `deletions` int(11) DEFAULT 0,
     `changed_files` int(11) DEFAULT 0,
     `title` varchar(250) NOT NULL,
     `repository_name` varchar(100) NOT NULL,
     `body` longtext DEFAULT NULL,
     `labels` longtext NULL,
     `owner_login` varchar(100) NOT NULL,
     `owner_id` int(11) NOT NULL,
     `first_commit_at` datetime DEFAULT '0000-00-00 00:00:00',
     `pr_created_at` datetime NOT NULL DEFAULT current_timestamp(),
     `pr_merged_at` datetime DEFAULT '0000-00-00 00:00:00',
     `pr_updated_at` datetime DEFAULT '0000-00-00 00:00:00',
     `pr_closed_at` datetime DEFAULT '0000-00-00 00:00:00',
     `transform_at` datetime DEFAULT '0000-00-00 00:00:00',
     `transform_status` enum('pending','done') DEFAULT 'pending',
     PRIMARY KEY (`id`),
     UNIQUE KEY `pr_id_UNIQUE` (`pr_id`),
     KEY `transform_status_idx` (`transform_status`),
     KEY `pr_created_at_idx` (`pr_created_at`),
     KEY `pr_closed_at_idx` (`pr_closed_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transformed_pull_request_count`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transformed_pull_request_count` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `open` int(11) DEFAULT 0,
    `closed` int(11) DEFAULT 0,
    `total` int(11) DEFAULT 0,
    `date` datetime DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `date_UNIQUE` (`date`),
    KEY `time` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transformed_pull_request_data`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `transformed_pull_request_data` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `pr_number` int(11) NOT NULL,
   `pr_id` int(11) NOT NULL,
   `duration` int(11) DEFAULT 0,
   `status` enum('open','closed','merged') NOT NULL DEFAULT 'open' COMMENT '1 - Open, 2 - Closed, 3 - Merged',
   `title` varchar(200) NOT NULL,
   `branch_name` VARCHAR(200) NOT NULL,
   `body` longtext DEFAULT NULL,
   `repository_name` varchar(200) NOT NULL,
   `pr_created_at` datetime NOT NULL,
   `pr_closed_at` datetime DEFAULT '0000-00-00 00:00:00',
   `owner_login_name` varchar(200) NOT NULL,
   `owner_id` int(11) NOT NULL,
   `team_slug` varchar(200) NOT NULL,
   `labels` longtext NULL,
   `comments` int(11) DEFAULT 0,
   `review_comments` int(11) DEFAULT 0,
   `commits` int(11) DEFAULT 0,
   `additions` int(11) DEFAULT 0,
   `deletions` int(11) DEFAULT 0,
   `changed_files` int(11) DEFAULT 0,
   `first_commit_at` datetime DEFAULT '0000-00-00 00:00:00',
   `pr_updated_at` datetime DEFAULT '0000-00-00 00:00:00',
   `pr_merged_at` datetime DEFAULT '0000-00-00 00:00:00',
   PRIMARY KEY (`id`),
   UNIQUE KEY `pr_id_UNIQUE` (`pr_id`),
   KEY `status_idx` (`status`),
   KEY `pr_created_at_idx` (`pr_created_at`),
   KEY `pr_closed_at_idx` (`pr_closed_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

--
-- Table structure for table `transformed_team_pull_request_count`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `transformed_team_pull_request_count` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `open` int(11) DEFAULT 0,
    `closed` int(11) DEFAULT 0,
    `total` int(11) DEFAULT 0,
    `team_slug` varchar(200) NOT NULL,
    `date` datetime DEFAULT current_timestamp(),
    `close_total_ratio` decimal(3,2) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `team_date_UNIQUE` (`team_slug`,`date`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

-- Dump completed on 2020-12-22 14:03:42
