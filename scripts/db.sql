
DROP DATABASE IF EXISTS db_scene_go;
CREATE DATABASE db_scene_go;
USE db_scene_go;


-- app表：记录接入的第三方app信息
CREATE TABLE `app` (
`app_id` VARCHAR(64) NOT NULL,
`name` varchar(128) NOT NULL DEFAULT '' COMMENT '应用名称',
`private_key` varchar(128) NOT NULL DEFAULT '' COMMENT '应用私钥',
PRIMARY KEY (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- space表：记录每个app创建的场景
CREATE TABLE `space` (
  `app_id` VARCHAR(64) NOT NULL COMMENT '应用id',
  `space_id` VARCHAR(64) NOT NULL COMMENT '应用的场景id',
  `grid_width` FLOAT unsigned NOT NULL COMMENT '场景九宫格宽',
  `grid_height` FLOAT unsigned NOT NULL COMMENT '场景九宫格高',
  PRIMARY KEY (`app_id`,`space_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- appId + userId -> spaceId 玩家最后一次登录的场景
CREATE TABLE `last_space` (
  `app_id` VARCHAR(64) NOT NULL COMMENT '应用id',
  `user_id` int(11) unsigned NOT NULL COMMENT '应用玩家id',
  `space_id` VARCHAR(64) NOT NULL COMMENT '最后一次登录的场景id',
  PRIMARY KEY (`app_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- appId + userId + spaceId -> x,y,angle	记录玩家在每个场景退出时的位置
CREATE TABLE `last_pos` (
  `app_id` VARCHAR(64) NOT NULL COMMENT '应用id',
  `user_id` int(11) unsigned NOT NULL COMMENT '应用玩家id',
  `space_id` VARCHAR(64) NOT NULL COMMENT '场景id',
  `x` FLOAT NOT NULL COMMENT '位置x',
  `y` FLOAT NOT NULL COMMENT '位置y',
  `angle` FLOAT NOT NULL COMMENT '位置角度',
  PRIMARY KEY (`app_id`,`space_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


