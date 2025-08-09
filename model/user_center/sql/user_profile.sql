-- 用户资料表（存储用户基本信息）
CREATE TABLE `sso_user_profile` (
  `user_id` BIGINT NOT NULL COMMENT '用户唯一标识（关联UserAuth）',
  `nickname` VARCHAR(50) NOT NULL COMMENT '用户昵称',
  `avatar` VARCHAR(255) NULL DEFAULT NULL COMMENT '头像URL',
  `bio` TEXT NULL DEFAULT NULL COMMENT '个人简介(可选)',
  `gender` TINYINT NOT NULL DEFAULT 0 COMMENT '性别: 0-未知, 1-男, 2-女',
  `region` VARCHAR(100) NULL DEFAULT NULL COMMENT '地区(可选，如"北京")',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间戳',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_user_profile_auth` FOREIGN KEY (`user_id`) REFERENCES `sso_user_auth` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基本资料表';
