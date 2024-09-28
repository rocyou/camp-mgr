CREATE DATABASE camp;

USE camp;

CREATE TABLE IF NOT EXISTS `campaigns` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_name` varchar(20) NOT NULL,
    `campaign_id` BIGINT(20) NOT NULL,
    `message_template` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `scheduled_time` bigint(20) NOT NULL,
    `csv_file_path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `camp_id` (`campaign_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `recipients` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `phone` (`phone`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `messages_0` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_0` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_0` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_1` (
                                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_1` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_1` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_2` (
                                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_2` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_2` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_3` (
                                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_3` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_3` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_4` (
                                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_4` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_4` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_5` (
                                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_5` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_5` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `messages_6` (
                                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_6` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_6` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages_7` (
                                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `campaign_id` bigint(20) NOT NULL,
    `name` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `send_at` bigint(20) NOT NULL,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '0: not start, 1: send success, 2: send failed',
    `deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `camp_id` (`campaign_id`),
    UNIQUE KEY `camp_phone` (`campaign_id`, `phone`),
    CONSTRAINT `fk_campaign_id_7` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`campaign_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_phone_7` FOREIGN KEY (`phone`) REFERENCES `recipients` (`phone`) ON DELETE CASCADE ON UPDATE CASCADE

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

commit;

