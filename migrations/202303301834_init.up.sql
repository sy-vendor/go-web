-- create "coin_market_cap_stats" table
CREATE TABLE `user` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, `sex` tinyint unsigned NOT NULL DEFAULT 0, `age` int NOT NULL DEFAULT 0, `account` varchar(255) NOT NULL, `password` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `account` (`account`)) CHARSET utf8mb4;