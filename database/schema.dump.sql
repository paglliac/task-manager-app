-- MySQL dump 10.13  Distrib 5.7.27, for Linux (x86_64)
--
-- Host: localhost    Database: golang
-- ------------------------------------------------------
-- Server version       5.7.27

/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE = @@TIME_ZONE */;
/*!40103 SET TIME_ZONE = '+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS, UNIQUE_CHECKS = 0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES = @@SQL_NOTES, SQL_NOTES = 0 */;

--
-- Table structure for table `labels`
--

DROP TABLE IF EXISTS `labels`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `labels`
(
    `id`   int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(256) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statuses`
--

DROP TABLE IF EXISTS `statuses`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statuses`
(
    `name`        varchar(256)                         NOT NULL,
    `title`       varchar(256) DEFAULT NULL,
    `status_type` enum ('open','in_progress','closed') NOT NULL,
    PRIMARY KEY (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `task_comments`
--

DROP TABLE IF EXISTS `task_comments`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `task_comments`
(
    `id`         char(36) CHARACTER SET ascii NOT NULL,
    `task_id`    char(36) DEFAULT NULL,
    `message`    text,
    `created_at` datetime                     NOT NULL,
    `author`     varchar(255)                 NOT NULL,
    PRIMARY KEY (`id`),
    KEY `task_comments_tasks_id_fk` (`task_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `task_labels`
--

DROP TABLE IF EXISTS `task_labels`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `task_labels`
(
    `task_id`  int(11) NOT NULL,
    `label_id` int(11) DEFAULT NULL,
    KEY `task_labels_labels_id_fk` (`label_id`),
    KEY `task_labels_tasks_id_fk` (`task_id`),
    CONSTRAINT `task_labels_labels_id_fk` FOREIGN KEY (`label_id`) REFERENCES `labels` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tasks`
(
    `id`          char(36) CHARACTER SET ascii NOT NULL,
    `title`       varchar(70)                  NOT NULL,
    `description` text,
    `status`      varchar(256) DEFAULT NULL,
    `created_at`  datetime     DEFAULT NULL,
    `updated_at`  datetime     DEFAULT NULL,
    `author`      int(11)                      NOT NULL,
    PRIMARY KEY (`id`),
    KEY `tasks_statuses_name_fk` (`status`),
    KEY `tasks_users_id_fk` (`author`),
    CONSTRAINT `tasks_statuses_name_fk` FOREIGN KEY (`status`) REFERENCES `statuses` (`name`) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT `tasks_users_id_fk` FOREIGN KEY (`author`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tasks_events`
--

DROP TABLE IF EXISTS `tasks_events`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tasks_events`
(
    `id`          int(11)                      NOT NULL AUTO_INCREMENT,
    `task_id`     char(36) CHARACTER SET ascii NOT NULL,
    `event_type`  varchar(256)                 NOT NULL,
    `payload`     text COMMENT 'JSON encoded data for event description\n',
    `occurred_on` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `tasks_events_tasks_id_fk` (`task_id`),
    CONSTRAINT `tasks_events_tasks_id_fk` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 52
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users`
(
    `id`       int(11) NOT NULL AUTO_INCREMENT,
    `name`     varchar(256) DEFAULT NULL,
    `email`    varchar(256) DEFAULT NULL,
    `password` varchar(256) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE = @OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE = @OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS = @OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION = @OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES = @OLD_SQL_NOTES */;

-- Dump completed on 2020-02-17 14:59:29
