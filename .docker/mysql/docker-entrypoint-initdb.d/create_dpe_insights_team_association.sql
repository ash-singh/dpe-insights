-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Database: dpe_insights
-- ------------------------------------------------------
-- Server version	5.5.5-10.5.8-MariaDB-1:10.5.8+maria~focal

USE `dpe_insights`;

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
-- Table structure for table `extracted_repositories`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_repositories` (
    `id` int(11) NOT NULL COMMENT 'Github repository id',
    `name` varchar(200) NOT NULL,
    `is_private` tinyint(4) NOT NULL DEFAULT 1,
    `description` longtext DEFAULT NULL,
    `size` int(11) NOT NULL,
    `open_issues` int(11) DEFAULT NULL,
    `language` varchar(100) NOT NULL COMMENT 'Source code primary language',
    `is_archived` tinyint(4) NOT NULL DEFAULT 0,
    `is_disabled` tinyint(4) NOT NULL DEFAULT 0,
    `updated_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `extracted_teams`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_teams` (
    `id` int(11) NOT NULL COMMENT 'Github team id',
    `name` varchar(30) NOT NULL,
    `slug` varchar(30) NOT NULL,
    `description` varchar(150) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `extracted_team_repositories`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_team_repositories` (
     `github_team_id` int(11) NOT NULL,
     `github_repository_id` int(11) NOT NULL,
     PRIMARY KEY (`github_team_id`,`github_repository_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `extracted_team_users`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_team_users` (
    `github_team_id` int(11) NOT NULL,
    `github_user_id` int(11) NOT NULL,
    PRIMARY KEY (`github_team_id`,`github_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `extracted_users`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `extracted_users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `github_id` int(11) NOT NULL COMMENT 'Github user id',
    `github_login_name` varchar(60) NOT NULL COMMENT 'Github login/username ',
    `github_user_type` varchar(30) NOT NULL,
    `is_site_admin` tinyint(4) DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `github_id_UNIQUE` (`github_id`),
    UNIQUE KEY `github_login_name_UNIQUE` (`github_login_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

-- Dump completed on 2020-12-02 18:27:27
