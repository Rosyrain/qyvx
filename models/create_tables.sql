DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
                        `user_id` bigint(20) NOT NULL COMMENT '后端生成id',
                        `github_id` bigint(20) DEFAULT NULL COMMENT 'github id',
                        `qyvx_id` varchar(255) NOT NULL COMMENT '企业微信id',
                        `github_name` varchar(255) DEFAULT NULL COMMENT 'github name',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE,
                        UNIQUE KEY `idx_qyvx_id` (`qyvx_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

show tables;
