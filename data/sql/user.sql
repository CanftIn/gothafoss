-- +migrate Up
-- 用户表
create table `user` (
  id integer not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 用户唯一ID
  name VARCHAR(100) not null default '',
  -- 用户的名字
  short_no VARCHAR(40) not null default '',
  -- 短编码
  short_status smallint not null default 0,
  -- 短编码 0.未修改 1.已修改
  sex smallint not null default 0,
  -- 性别 0.女 1.男
  robot smallint not null default 0,
  -- 机器人 0.否1.是
  category VARCHAR(40) not null default '',
  -- 用户分类  service:客服
  role VARCHAR(40) not null default '',
  -- 用户角色  admin:管理员 superAdmin
  username VARCHAR(40) not null default '',
  -- 用户名
  password VARCHAR(40) not null default '',
  -- 密码
  zone VARCHAR(40) not null default '',
  -- 手机区号
  phone VARCHAR(20) not null default '',
  -- 手机号
  chat_pwd VARCHAR(40) not null default '',
  -- 聊天密码
  lock_screen_pwd varchar(40) not null default '',
  -- 锁屏密码
  lock_after_minute integer not null default 0,
  -- 在几分钟后锁屏 0 表示立即
  vercode VARCHAR(100) not null default '',
  -- 验证码 加好友来源
  is_upload_avatar smallint not null default 0,
  -- 是否上传过头像 1:上传0:未上传
  qr_vercode VARCHAR(100) not null default '',
  -- 二维码验证码 加好友来源
  device_lock smallint not null DEFAULT 0,
  -- 是否开启设备锁
  search_by_phone smallint not null default 1,
  -- 是否可用通过手机号搜索到本人0.否1.是
  search_by_short smallint not null default 1,
  -- 是否可以通过短编号搜索0.否1.是
  new_msg_notice smallint not null default 1,
  -- 新消息通知0.否1.是
  msg_show_detail smallint not null default 1,
  -- 新消息通知详情0.否1.是
  voice_on smallint not null default 1,
  -- 是否开启声音0.否1.是
  shock_on smallint not null default 1,
  -- 是否开启震动0.否1.是
  mute_of_app smallint not null default 0,
  -- app是否禁音（当pc登录的时候app可以设置禁音，当pc登录后有效）
  offline_protection smallint not null default 0,
  -- 离线保护，断网屏保
  `version` bigint not null DEFAULT 0,
  -- 数据版本 
  status smallint not null DEFAULT 1,
  -- 用户状态 0.禁用 1.可用
  bench_no VARCHAR(40) not null default '',
  -- 性能测试批次号，性能测试幂等用
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX uid on `user` (uid);

CREATE UNIQUE INDEX short_no_udx on `user` (short_no);

-- -- +migrate StatementBegin
-- CREATE TRIGGER user_updated_at
--   BEFORE UPDATE
--   ON `user` for each row 
--   BEGIN
--     set NEW.updated_at = NOW();
--   END;
-- -- +migrate StatementEnd
-- 创建系统账号
INSERT INTO
  `user` (
    uid,
    name,
    short_no,
    phone,
    zone,
    search_by_phone,
    search_by_short,
    new_msg_notice,
    voice_on,
    shock_on,
    msg_show_detail,
    status,
    is_upload_avatar,
    category,
    robot
  )
VALUES
  (
    'u_10000',
    '系统账号',
    10000,
    '13000000000',
    '0086',
    0,
    0,
    0,
    0,
    0,
    0,
    1,
    1,
    'system',
    1
  );

INSERT INTO
  `user` (
    uid,
    name,
    short_no,
    phone,
    zone,
    search_by_phone,
    search_by_short,
    new_msg_notice,
    voice_on,
    shock_on,
    msg_show_detail,
    status,
    is_upload_avatar,
    category
  )
VALUES
  (
    'fileHelper',
    '文件传输助手',
    20000,
    '13000000001',
    '0086',
    0,
    0,
    0,
    0,
    0,
    0,
    1,
    1,
    'system'
  );

-- 创建后台管理平台超级管理员账号 admin/admiN123456
-- INSERT INTO `user` (uid,name,short_no,username,password,role,phone,zone,search_by_phone,search_by_short,new_msg_notice,voice_on,shock_on,msg_show_detail,status,is_upload_avatar,category) VALUES ('admin','超级管理员',30000,'admin','14c3a0db22308e34ca7dacb1806c0bdf','superAdmin','13000000002','0086',0,0,0,0,0,0,1,0,'system');
-- 用户设置
create table `user_setting` (
  id integer not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 用户UID
  to_uid VARCHAR(40) not null default '',
  -- 对方uid
  mute smallint not null DEFAULT 0,
  --  是否免打扰
  top smallint not null DEFAULT 0,
  -- 是否置顶
  blacklist smallint not null DEFAULT 0,
  -- 是否黑名单 0:正常1:黑名单
  chat_pwd_on smallint not null DEFAULT 0,
  -- 是否开启聊天密码
  screenshot smallint not null DEFAULT 1,
  -- 截屏通知
  revoke_remind smallint not null DEFAULT 1,
  -- 撤回通知
  receipt smallint not null default 1,
  -- 消息是否回执
  version BIGINT not null DEFAULT 0,
  -- 版本
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE INDEX uid_idx on `user_setting` (uid);

-- -- +migrate StatementBegin
-- CREATE TRIGGER user_setting_updated_at
--   BEFORE UPDATE
--   ON `user_setting` for each row 
--   BEGIN
--     set NEW.updated_at = NOW();
--   END;
-- -- +migrate StatementEnd
-- 用户设备
create table `device` (
  id integer not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 设备所属用户uid                     
  device_id VARCHAR(40) not null default '',
  -- 设备唯一ID          
  device_name VARCHAR(100) not null default '',
  -- 设备名称                  
  device_model VARCHAR(100) not null default '',
  -- 设备型号              
  last_login integer not null DEFAULT 0,
  -- 最后一次登录时间(时间戳 10位)
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE unique INDEX device_uid_device_id on `device` (uid, device_id);

CREATE INDEX device_uid on `device` (uid);

CREATE INDEX device_device_id on `device` (device_id);

-- 好友表
create table `friend` (
  id integer not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '' comment '用户UID',
  to_uid VARCHAR(40) not null default '' comment '好友uid',
  remark varchar(100) not null default '' comment '对好友的备注 TODO: 此字段不再使用，已经迁移到user_setting表',
  flag smallint not null default 0 comment '好友标示',
  `version` bigint not null default 0 comment '版本号',
  vercode VARCHAR(100) not null default '' comment '验证码 加好友来源',
  source_vercode varchar(100) not null default '' comment '好友来源',
  is_deleted smallint not null default 0 comment '是否已删除',
  is_alone smallint not null default 0 comment '单项好友',
  initiator smallint not null default 0 comment '加好友发起方',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '创建时间',
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '更新时间'
);

-- -- +migrate StatementBegin
-- CREATE TRIGGER friend_updated_at
--   BEFORE UPDATE
--   ON `friend` for each row 
--   BEGIN
--    set NEW.updated_at = NOW();
--   END;
-- -- +migrate StatementEnd
-- +migrate Up
-- 登录日志
CREATE TABLE IF NOT EXISTS login_log(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  uid VARCHAR(40) DEFAULT '' NOT NULL COMMENT '用户OpenId',
  login_ip VARCHAR(40) DEFAULT '' NOT NULL COMMENT '最后一次登录ip',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '创建时间',
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP comment '更新时间'
) CHARACTER SET utf8mb4;

-- +migrate Up
ALTER TABLE
  `user`
ADD
  COLUMN app_id VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'app id';

-- +migrate Up
ALTER TABLE
  `user`
ADD
  COLUMN email VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'email地址';

-- +migrate Up
-- 消息表
create table `user_online` (
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 用户uid
  device_flag smallint not null default 0,
  -- 设备flag 0.APP 1. WEB
  last_online integer not null DEFAULT 0,
  -- 最后一次在线时间
  last_offline integer not null DEFAULT 0,
  -- 最后一次离线时间
  online tinyint(1) not null default 0,
  -- 用户是否在线
  `version` bigint not null default 0,
  -- 数据版本
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX `uid_device` on `user_online` (`uid`, device_flag);

CREATE INDEX `online_idx` on `user_online` (`online`);

CREATE INDEX `uid_idx` on `user_online` (`uid`);

-- +migrate Up
-- 用户身份表 （signal protocol使用）
create table `signal_identities` (
  id bigint not null primary key AUTO_INCREMENT,
  uid varchar(40) not null DEFAULT '',
  --  用户uid
  registration_id bigint not null DEFAULT 0,
  -- 身份ID
  identity_key text not null,
  -- 用户身份公钥
  signed_prekey_id integer not null DEFAULT 0,
  -- 签名key的id
  signed_pubkey text not null,
  -- 签名key的公钥
  signed_signature text not null,
  -- 由身份密钥签名的signed_pubkey
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX identities_index_id ON signal_identities(uid);

-- 一次性公钥
create table `signal_onetime_prekeys` (
  id bigint not null primary key AUTO_INCREMENT,
  uid varchar(40) not null DEFAULT '',
  -- 用户uid
  key_id integer not null DEFAULT 0,
  pubkey text not null,
  -- 公钥
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX key_id_uid_index_id ON signal_onetime_prekeys(uid, key_id);

-- +migrate Up
-- 手机联系人
create table `user_maillist` (
  id bigint not null primary key AUTO_INCREMENT,
  uid VARCHAR(40) not null default '',
  -- 用户uid
  phone VARCHAR(40) not null default '',
  -- 手机号
  zone VARCHAR(40) not null default '',
  -- 区号
  name VARCHAR(40) not null default '',
  -- 名字
  vercode VARCHAR(100) not null default '',
  -- 验证码 加好友来源
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX `uid_maillist_index` on `user_maillist` (`uid`, `zone`, `phone`);

-- +migrate Up
CREATE unique INDEX to_uid_uid on `user_setting` (uid, to_uid);

CREATE unique INDEX to_uid_uid on `friend` (uid, to_uid);

-- +migrate Up
ALTER TABLE
  `user`
ADD
  COLUMN is_destroy smallint not null default 0 COMMENT '是否已销毁';

ALTER TABLE
  `user`
MODIFY
  COLUMN zone VARCHAR(20);

ALTER TABLE
  `user`
MODIFY
  COLUMN phone VARCHAR(100);

-- +migrate Up
ALTER TABLE
  `user_setting`
ADD
  COLUMN remark VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户备注';

-- 迁移备注数据
insert into
  user_setting(uid, to_uid, remark)
select
  uid,
  to_uid,
  remark
from
  friend
where
  remark <> '' on duplicate key
update
  remark =
values
(remark);

-- +migrate Up
ALTER TABLE
  `user`
ADD
  COLUMN wx_openid VARCHAR(100) NOT NULL DEFAULT '' COMMENT '微信openid';

ALTER TABLE
  `user`
ADD
  COLUMN wx_unionid VARCHAR(100) NOT NULL DEFAULT '' COMMENT '微信unionid';

-- +migrate Up
ALTER TABLE
  `user_setting`
ADD
  COLUMN `flame` smallint not null default 0 COMMENT '阅后即焚是否开启 1.开启 0.未开启';

ALTER TABLE
  `user_setting`
ADD
  COLUMN `flame_second` smallint not null default 0 COMMENT '阅后即焚销毁秒数';

-- +migrate Up
-- 短编号
create table `shortno` (
  id bigint not null primary key AUTO_INCREMENT,
  shortno VARCHAR(40) not null default '' COMMENT '唯一短编号',
  used smallint not null default 0 COMMENT '是否被用',
  hold smallint not null default 0 COMMENT '保留，保留的号码将不会再被分配',
  locked smallint not null default 0 COMMENT '是否被锁定，锁定了的短编号将不再被分配,直到解锁',
  business VARCHAR(40) not null default '' COMMENT '被使用的业务，比如 user',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX `udx_shortno` on `shortno` (`shortno`);

-- +migrate Up
--  设备标识
create table `device_flag` (
  id bigint not null primary key AUTO_INCREMENT,
  device_flag smallint not null default 0 COMMENT '设备标记 0. app 1.Web 2.PC',
  `weight` integer not null default 0 COMMENT '设备权重 值越大越优先',
  remark VARCHAR(100) not null default '' COMMENT '备注',
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

CREATE UNIQUE INDEX `udx_device_flag` on `device_flag` (`device_flag`);

insert into
  device_flag(device_flag, `weight`, remark)
values
(2, '80000', 'PC');

insert into
  device_flag(device_flag, `weight`, remark)
values
(1, '70000', 'Web');

insert into
  device_flag(device_flag, `weight`, remark)
values
(0, '90000', '手机');

-- +migrate Up
--  gitee_user 用户信息
CREATE TABLE IF NOT EXISTS gitee_user(
  id BIGINT PRIMARY KEY DEFAULT 0 COMMENT '用户 ID',
  `login` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户名',
  name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户姓名',
  email VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  bio VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户简介',
  avatar_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户头像 URL',
  blog VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户博客 URL',
  events_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户事件 URL',
  followers INT NOT NULL DEFAULT 0 COMMENT '用户粉丝数',
  followers_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户粉丝 URL',
  following INT NOT NULL DEFAULT 0 COMMENT '用户关注数',
  following_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户关注 URL',
  gists_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户 Gist URL',
  html_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户主页 URL',
  member_role VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户角色',
  organizations_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户组织 URL',
  public_gists INT NOT NULL DEFAULT 0 COMMENT '用户公开 Gist 数',
  public_repos INT NOT NULL DEFAULT 0 COMMENT '用户公开仓库数',
  received_events_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户接收事件 URL',
  remark VARCHAR(100) NOT NULL DEFAULT '' COMMENT '企业备注名',
  repos_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户仓库 URL',
  stared INT NOT NULL DEFAULT 0 COMMENT '用户收藏数',
  starred_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户收藏 URL',
  subscriptions_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户订阅 URL',
  url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户 URL',
  watched INT NOT NULL DEFAULT 0 COMMENT '用户关注的仓库数',
  weibo VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '用户微博 URL',
  `type` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '用户类型',
  `gitee_created_at` VARCHAR(30) NOT NULL DEFAULT '' COMMENT 'gitee用户创建时间',
  `gitee_updated_at` VARCHAR(30) NOT NULL DEFAULT '' COMMENT 'gitee用户更新时间',
  created_at timeStamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at timeStamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

CREATE unique INDEX gitee_user_login on `gitee_user` (`login`);

ALTER TABLE
  `user`
ADD
  COLUMN gitee_uid VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'gitee的用户id';

-- gthub用户
CREATE TABLE IF NOT EXISTS github_user (
  id BIGINT PRIMARY KEY DEFAULT 0 COMMENT '用户 ID',
  login VARCHAR(100) NOT NULL COMMENT '登录名',
  node_id VARCHAR(255) NOT NULL COMMENT '节点ID',
  avatar_url VARCHAR(1000) NOT NULL COMMENT '头像URL',
  gravatar_id VARCHAR(1000) NOT NULL COMMENT 'Gravatar ID',
  url VARCHAR(1000) NOT NULL COMMENT 'GitHub URL',
  html_url VARCHAR(1000) NOT NULL COMMENT 'GitHub HTML URL',
  followers_url VARCHAR(1000) NOT NULL COMMENT '关注者URL',
  following_url VARCHAR(1000) NOT NULL COMMENT '被关注者URL',
  gists_url VARCHAR(1000) NOT NULL COMMENT '代码片段URL',
  starred_url VARCHAR(1000) NOT NULL COMMENT '收藏URL',
  subscriptions_url VARCHAR(1000) NOT NULL COMMENT '订阅URL',
  organizations_url VARCHAR(1000) NOT NULL COMMENT '组织URL',
  repos_url VARCHAR(1000) NOT NULL COMMENT '仓库URL',
  events_url VARCHAR(1000) NOT NULL COMMENT '事件URL',
  received_events_url VARCHAR(1000) NOT NULL COMMENT '接收事件URL',
  `type` VARCHAR(100) NOT NULL COMMENT '用户类型',
  site_admin BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否为管理员',
  name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '姓名',
  company VARCHAR(100) NOT NULL DEFAULT '' COMMENT '公司',
  blog VARCHAR(100) NOT NULL DEFAULT '' COMMENT '博客',
  location VARCHAR(255) NOT NULL DEFAULT '' COMMENT '所在地',
  email VARCHAR(100) NOT NULL DEFAULT '' COMMENT '电子邮件',
  hireable BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否可被雇佣',
  bio VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '个人简介',
  twitter_username VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'Twitter 用户名',
  public_repos INT NOT NULL DEFAULT 0 COMMENT '公共仓库数量',
  public_gists INT NOT NULL DEFAULT 0 COMMENT '公共代码片段数量',
  followers INT NOT NULL DEFAULT 0 COMMENT '关注者数量',
  following INT NOT NULL DEFAULT 0 COMMENT '被关注者数量',
  github_created_at VARCHAR(30) NOT NULL DEFAULT '' COMMENT '创建时间',
  github_updated_at VARCHAR(30) NOT NULL DEFAULT '' '更新时间',
  private_gists INT NOT NULL DEFAULT 0 COMMENT '私有代码片段数量',
  total_private_repos INT NOT NULL DEFAULT 0 COMMENT '私有仓库总数',
  owned_private_repos INT NOT NULL DEFAULT 0 COMMENT '拥有的私有仓库数量',
  disk_usage INT NOT NULL DEFAULT 0 COMMENT '磁盘使用量',
  collaborators INT NOT NULL DEFAULT 0 COMMENT '协作者数量',
  two_factor_authentication BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否启用两步验证',
  created_at timeStamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at timeStamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间'
);

CREATE unique INDEX github_user_login on `github_user` (`login`);

ALTER TABLE
  `user`
ADD
  COLUMN github_uid VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'github的用户id';