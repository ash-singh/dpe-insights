USE `dpe_insights`;

CREATE TABLE `extracted_pagerduty_incidents`
(
    `id`                    INT(10) UNSIGNED NOT NULL,
    `created_at`            TIMESTAMP NULL DEFAULT NULL,
    `last_status_change_at` TIMESTAMP NULL DEFAULT NULL,
    `service_name`          VARCHAR(255) NULL DEFAULT NULL,
    `duration`              INT(10) UNSIGNED NULL DEFAULT NULL,
    `status`                VARCHAR(50) NULL DEFAULT NULL,
    `urgency`               VARCHAR(50) NULL DEFAULT NULL,
    `title`                 TEXT NULL DEFAULT NULL,
    `description`           TEXT NULL DEFAULT NULL,
    `false_positive`        TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX                   `created_at` (`created_at`),
    INDEX                   `service_name` (`service_name`),
    INDEX                   `status` (`status`)
) ENGINE=InnoDB
;
