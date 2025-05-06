CREATE TABLE `jpid` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'java项目名',
                        `port` int NOT NULL COMMENT '运行端口',
                        `pid` int NOT NULL COMMENT 'pid',
                        `catalog` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '运行目录',
                        `run` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '运行脚本（sh命令',
                        `description` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '项目描述',
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='java项目详情';
