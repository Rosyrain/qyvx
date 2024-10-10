DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `qyvx_id` varchar(255) NOT NULL COMMENT '企业微信id',
    `github_id` varchar(50) DEFAULT NULL COMMENT 'github id',
    `name` varchar(50) DEFAULT NULL COMMENT '姓名',
    `github_name` varchar(255) DEFAULT NULL COMMENT 'github name',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1-在职，2-离职',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`qyvx_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `shifts`;
CREATE TABLE  `shifts` (
   `start_date` DATE NOT NULL COMMENT '开始值班日期',
   `end_date` DATE NOT NULL COMMENT '结束值班日期',
   `oncallers` varchar(100) NOT NULL COMMENT '值班人员',
   PRIMARY KEY (`start_date`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
insert into `shifts` (`start_date`, `end_date`, `oncallers`) values ('2023-01-01', '2023-01-07', '值班人员1,值班人员2,值班人员3');
insert into `shifts` (`start_date`, `end_date`, `oncallers`) values ('2024-01-01', '2024-01-07', '值班人员4,值班人员5,值班人员6');

show tables;
