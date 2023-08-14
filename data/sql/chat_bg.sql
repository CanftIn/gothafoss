-- +migrate Up
create table `chat_bg`(
  id integer not null primary key AUTO_INCREMENT,
  cover varchar(100) not null default '',
  -- 封面
  url varchar(100) not null default '',
  -- 图片地址
  is_svg smallint not null default 1,
  -- 是否为svg图片
  created_at timeStamp not null DEFAULT CURRENT_TIMESTAMP,
  -- 创建时间
  updated_at timeStamp not null DEFAULT CURRENT_TIMESTAMP -- 更新时间
);