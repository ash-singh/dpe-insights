USE `dpe_insights`;

CREATE TABLE `extracted_releases`
(
    `id`                    int(11) NOT NULL AUTO_INCREMENT,
    `release_id`            INT(10) UNSIGNED NOT NULL,
    `title`                 varchar(200),
    `tag_name`              VARCHAR(100) NULL DEFAULT NULL,
    `body`                  longtext DEFAULT NULL,
    `repository_id`         INT(10) UNSIGNED NOT NULL,
    `author_login`          VARCHAR(100) NULL DEFAULT NULL,
    `created_at`            TIMESTAMP NULL DEFAULT NULL,
    `published_at`          TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `release_id_UNIQUE` (`release_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4
;

