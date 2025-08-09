-- 用户认证表（存储登录凭证信息）
CREATE TABLE `sso_user_auth` (
  `user_id` BIGINT NOT NULL COMMENT '用户唯一标识（关联UserProfile）',
  `mobile` VARCHAR(20) NULL DEFAULT NULL COMMENT '手机号(可选)',
  `email` VARCHAR(100) NULL DEFAULT NULL COMMENT '邮箱(可选)',
  `account` VARCHAR(50) NULL DEFAULT NULL COMMENT '账号(可选)',
  `password` VARCHAR(100) NOT NULL COMMENT '密码(加密存储)',
  `salt` VARCHAR(32) NOT NULL COMMENT '密码盐值',
  `platform` VARCHAR(20) NULL DEFAULT NULL COMMENT '第三方平台: wechat/qq(可选)',
  `open_id` VARCHAR(64) NULL DEFAULT NULL COMMENT '第三方平台OpenID(可选)',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '用户状态: 0-正常, 1-禁用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间戳',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
  `last_login_at` DATETIME NULL DEFAULT NULL COMMENT '最后登录时间戳',
  `last_login_ip` VARCHAR(45) NULL DEFAULT NULL COMMENT '最后登录IP',
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `idx_account` (`account`) COMMENT '账号唯一索引（登录用）',
  INDEX `idx_platform_openid` (`platform`, `open_id`) COMMENT '第三方登录联合索引',
  INDEX `idx_status` (`status`) COMMENT '用户状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户认证信息表';
