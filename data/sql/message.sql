-- +migrate Up
-- 管理员发送消息记录
create table `send_history`(
  id integer not null primary key AUTO_INCREMENT,
  receiver VARCHAR(40) not null default '',
  -- 接受者uid
  receiver_name varchar(100) not null default '',
  -- 接受者
  receiver_channel_type smallint not null default 0,
  -- 接受者频道类型
  sender varchar(40) not null default '',
  -- 发送者uid
  sender_name varchar(100) not null default '',
  -- 发送者名字
  handler_uid varchar(40) not null default '',
  -- 操作者uid
  handler_name VARCHAR(100) not null default '',
  -- 操作者名字
  content TEXT,
  -- 发送内容   
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- +migrate Up
-- 消息扩展表
create table `message_extra` (
  id bigint not null primary key AUTO_INCREMENT,
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  message_seq bigint not null default 0,
  -- 消息序列号(严格递增)
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  from_uid VARCHAR(40) not null default '',
  -- 发送者uid
  `revoke` smallint not null default 0,
  -- 是否撤回
  revoker VARCHAR(40) not null default '',
  -- 是否撤回
  clone_no VARCHAR(40) not null default '',
  -- 未读编号
  -- voice_status smallint not null default 0, -- 语音状态 0.未读 1.已读
  `version` bigint not null default 0,
  -- 数据版本
  readed_count integer not null default 0,
  -- 已读数量
  is_deleted smallint not null default 0,
  -- 是否已删除
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX from_uid_idx on `message_extra` (from_uid);

CREATE INDEX channel_idx on `message_extra` (channel_id, channel_type);

CREATE UNIQUE INDEX message_id on `message_extra` (message_id);

-- +migrate Up
ALTER TABLE
  `message_extra`
ADD
  COLUMN content_edit TEXT COMMENT '编辑后的正文';

ALTER TABLE
  `message_extra`
ADD
  COLUMN content_edit_hash varchar(255) not null default '' COMMENT '编辑正文的hash值，用于重复判断';

ALTER TABLE
  `message_extra`
ADD
  COLUMN edited_at integer not null default 0 COMMENT '编辑时间 时间戳（秒）';

-- 成员已读列表
CREATE TABLE `member_readed`(
  id bigint not null primary key AUTO_INCREMENT,
  clone_no VARCHAR(40) not null default '',
  -- 克隆成员唯一编号
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  uid VARCHAR(40) not null default '',
  -- 已读用户uid
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX channel_idx on `member_readed` (channel_id, channel_type);

CREATE INDEX uid_idx on `member_readed` (uid);

CREATE UNIQUE INDEX message_uid_idx on `member_readed` (message_id, uid);

-- 成员克隆列表(TODO: 此表已作废)
CREATE TABLE `member_clone`(
  id bigint not null primary key AUTO_INCREMENT,
  clone_no VARCHAR(40) not null default '',
  -- 克隆成员唯一编号
  channel_id VARCHAR(40) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  uid VARCHAR(40) not null default '',
  -- 已读用户uid
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX clone_no_idx on `member_clone` (clone_no);

CREATE INDEX channel_idx on `member_clone` (channel_id, channel_type);

-- 频道成员变化记录
CREATE TABLE `member_change`(
  id bigint not null primary key AUTO_INCREMENT,
  clone_no VARCHAR(40) not null default '',
  -- 未读编号
  channel_id VARCHAR(40) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  max_version bigint not null default 0,
  -- 当前最大版本
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 最近会话扩展表
CREATE TABLE `conversation_extra`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '' comment '所属用户',
  channel_id VARCHAR(100) not null default '' comment '频道ID',
  channel_type smallint not null default 0 comment '频道类型',
  browse_to bigint not null default 0 comment '预览到的位置，与会话保持位置不同的是 预览到的位置是用户读到的最大的messageSeq。跟未读消息数量有关系',
  keep_message_seq bigint not null default 0 comment '会话保持的位置',
  keep_offset_y integer not null default 0 comment '会话保持的位置的偏移量',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '创建时间',
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '更新时间'
);

CREATE UNIQUE INDEX uid_channel_idx on `conversation_extra` (uid, channel_id, channel_type);

CREATE INDEX uid_idx on `conversation_extra` (uid);

-- +migrate Up
ALTER TABLE
  `conversation_extra`
ADD
  COLUMN draft varchar(1000) not null default '' COMMENT '草稿';

ALTER TABLE
  `conversation_extra`
ADD
  COLUMN `version` bigint not null default 0 COMMENT '数据版本';

-- +migrate Up
-- 回应用户
CREATE TABLE `reaction_users`(
  id bigint not null primary key AUTO_INCREMENT,
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  seq bigint not null default 0,
  --  回复递增序号（可以用此序号做递增操作）
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  uid varchar(40) not null default '',
  -- 回应的用户uid
  name varchar(40) not null default '',
  -- 回应的用户名
  emoji varchar(20) not null default '',
  -- 回应的emoji
  is_deleted smallint not null default 0,
  -- 是否已删除
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX `reaction_user_message_channel` on `reaction_users` (`message_id`, uid, `emoji`);

-- +migrate Up
--  用户独立对消息的扩充
CREATE TABLE `message_user_extra`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  message_seq bigint not null default 0,
  -- 消息序列号(严格递增)
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  voice_readed smallint not null default 0,
  -- 语音是否已读
  message_is_deleted smallint not null default 0,
  -- 消息是否已删除
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_message_idx on `message_user_extra` (uid, message_id);

CREATE TABLE `message_user_extra1`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  message_seq bigint not null default 0,
  -- 消息序列号(严格递增)
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  voice_readed smallint not null default 0,
  -- 语音是否已读
  message_is_deleted smallint not null default 0,
  -- 消息是否已删除
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_message_idx on `message_user_extra1` (uid, message_id);

CREATE TABLE `message_user_extra2`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  message_id VARCHAR(20) not null default '',
  -- 消息唯一ID（全局唯一）
  message_seq bigint not null default 0,
  -- 消息序列号(严格递增)
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  voice_readed smallint not null default 0,
  -- 语音是否已读
  message_is_deleted smallint not null default 0,
  -- 消息是否已删除
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_message_idx on `message_user_extra2` (uid, message_id);

-- 频道偏移表 （每个用户针对于频道的偏移位置）
CREATE TABLE `channel_offset`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  message_seq bigint not null default 0,
  -- 偏移的消息序号
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_channel_idx on `channel_offset` (uid, channel_id, channel_type);

CREATE TABLE `channel_offset1`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  message_seq bigint not null default 0,
  -- 偏移的消息序号
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_channel_idx on `channel_offset1` (uid, channel_id, channel_type);

CREATE TABLE `channel_offset2`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  message_seq bigint not null default 0,
  -- 偏移的消息序号
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid_channel_idx on `channel_offset2` (uid, channel_id, channel_type);

-- +migrate Up
-- 设备消息偏移量
CREATE TABLE `device_offset`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  device_uuid VARCHAR(40) not null default '',
  -- 设备唯一ID
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  message_seq bigint not null default 0,
  -- 偏移的消息序号
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX uid_device_offset_idx on `device_offset` (uid, device_uuid);

CREATE UNIQUE INDEX uid_device_offset_unidx on `device_offset` (uid, device_uuid, channel_id, channel_type);

-- 用户消息最新偏移量
CREATE TABLE `user_last_offset`(
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 编辑用户唯一ID
  channel_id VARCHAR(100) not null default '',
  -- 频道ID
  channel_type smallint not null default 0,
  -- 频道类型
  message_seq bigint not null default 0,
  -- 偏移的消息序号
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX uid_user_last_offset_idx on `user_last_offset` (uid);

CREATE UNIQUE INDEX uid_user_last_offset_unidx on `user_last_offset` (uid, channel_id, channel_type);

-- +migrate Up
CREATE TABLE `reminders`(
  id bigint not null primary key AUTO_INCREMENT,
  channel_id VARCHAR(100) not null default '' COMMENT '频道ID',
  channel_type smallint not null default 0 COMMENT '频道类型',
  reminder_type integer not null default 0 COMMENT '提醒类型 1.有人@我 2.草稿',
  uid varchar(40) not null default '' COMMENT '提醒的用户uid，如果此字段为空则表示 提醒项为整个频道内的成员',
  `text` varchar(255) not null default '' COMMENT '提醒内容',
  `data` varchar(1000) not null default '' COMMENT '自定义数据',
  is_locate smallint not null default 0 COMMENT ' 是否需要定位',
  message_seq bigint not null default 0 COMMENT '消息序列号',
  message_id VARCHAR(20) not null default '' COMMENT '消息唯一ID（全局唯一）',
  `version` bigint not null default 0 COMMENT ' 数据版本',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX channel_uid_uidx on `reminders` (uid, channel_id, channel_type);

-- +migrate Up
ALTER TABLE
  `reminders`
ADD
  COLUMN `client_msg_no` varchar(40) not null default '' COMMENT '消息client msg no';

ALTER TABLE
  `reminders`
ADD
  COLUMN `is_deleted` smallint not null default 0 COMMENT '是否被删除';

ALTER TABLE
  `reminders`
ADD
  COLUMN `publisher` varchar(40) not null default '' COMMENT '提醒项发布者uid';

CREATE TABLE `reminder_done`(
  id bigint not null primary key AUTO_INCREMENT,
  reminder_id bigint not null default 0 COMMENT '提醒事项的id',
  uid varchar(40) not null default '' COMMENT '完成的用户uid',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX reminder_id_uidx on `reminder_done` (reminder_id, uid);

-- +migrate Up
create table `prohibit_words`(
  id integer not null primary key AUTO_INCREMENT,
  is_deleted smallint not null default 0,
  -- 是否删除
  `version` bigint not null default 0,
  content TEXT,
  -- 内容   
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);