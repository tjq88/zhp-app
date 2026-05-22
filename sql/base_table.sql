CREATE TABLE `zp_xxx` (
                                    `id` int NOT NULL,
                                    `create_user` varchar(255) NOT NULL,
                                    `update_user` datetime NOT NULL,
                                    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    `tenant_code` varchar(255) NOT NULL,
                                    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;