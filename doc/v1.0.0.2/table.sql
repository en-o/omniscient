ALTER TABLE `omniscient`.`jpid` ADD COLUMN `way` int NULL DEFAULT 2 COMMENT '启动方式[1:docker, 2:jdk]' AFTER `description`
