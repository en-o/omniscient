CREATE TABLE `jpid` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'java项目名',
                        `ports` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '运行端口,多个逗号隔开',
                        `pid` int NOT NULL COMMENT 'pid',
                        `catalog` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '运行目录',
                        `run` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原生启动命令',
                        `script` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'sh脚本启动命令',
                        `status` int DEFAULT '0' COMMENT '状态[1:启动，0:停止]',
                        `description` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '项目描述',
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='java项目详情';
