-- +migrate Up
create table `app_config`(
  id integer not null primary key AUTO_INCREMENT,
  rsa_private_key varchar(4000) not null default '',
  -- 系统私钥 (使用来加密cmd类消息内容 防止前端模拟发送)
  rsa_public_key varchar(4000) not null default '',
  -- 系统公钥
  `version` integer not null default 0,
  -- 数据版本
  super_token varchar(40) not null default '',
  -- 超级token 用于操作一些系统api的安全校验
  super_token_on smallint not null default 0,
  -- 是否禁用super_token  0.禁用 1.开启 如果禁用 则一些需要super_token的API将不能使用 默认为禁用
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

ALTER TABLE
  `app_config`
ADD
  COLUMN revoke_second smallint not null DEFAULT 0 COMMENT '消息可撤回时长';

ALTER TABLE
  `app_config`
ADD
  COLUMN welcome_message varchar(2000) not null DEFAULT '' COMMENT '登录欢迎语';

ALTER TABLE
  `app_config`
ADD
  COLUMN new_user_join_system_group smallint not null DEFAULT 1 COMMENT '注册用户是否默认加入系统群';

ALTER TABLE
  `app_config`
ADD
  COLUMN search_by_phone smallint not null DEFAULT 0 COMMENT '是否可通过手机号搜索';
