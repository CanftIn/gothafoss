-- +migrate Up

-- 机器人
create table `robot`
(
  id         bigint         not null primary key AUTO_INCREMENT,
  robot_id   VARCHAR(40)    not null default '', -- 机器人ID
  token      VARCHAR(100)    not null default '', -- toekn
  `version`    BIGINT         not null DEFAULT 0,  -- 同步版本号
  status     smallint       not null DEFAULT 1, -- 机器人状态0:禁用1:启用
  created_at timeStamp      not null DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  updated_at timeStamp      not null DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

CREATE UNIQUE INDEX `robot_id_robot_index` on `robot` (`robot_id`);

-- 机器人菜单
create table `robot_menu`
(
  id         bigint         not null primary key AUTO_INCREMENT,
  robot_id   VARCHAR(40)    not null default '', -- 机器人ID
  cmd        VARCHAR(100)   not null default '', -- 命令
  remark     VARCHAR(100)   not null default '', -- 命令说明
  type       VARCHAR(100)   not null default '', -- 命令类型
  created_at timeStamp      not null DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  updated_at timeStamp      not null DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);
CREATE INDEX `bot_id_robot_menu_index` on `robot_menu` (`robot_id`);

-- +migrate Up

ALTER TABLE `robot` ADD COLUMN inline_on smallint not null DEFAULT 0 comment '是否开启行内搜索';
ALTER TABLE `robot` ADD COLUMN placeholder VARCHAR(40) not null DEFAULT '' comment '输入框占位符，开启行内搜索有效';

-- +migrate Up

ALTER TABLE `robot` ADD COLUMN username VARCHAR(40) not null DEFAULT '' comment '机器人的username';

ALTER TABLE `robot` ADD COLUMN app_id VARCHAR(40) not null DEFAULT '' comment '机器人所属app id';
