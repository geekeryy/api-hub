-- 用户信息表
create table member_info
(
    id                 bigint auto_increment,
    member_uuid        varchar(255) not null unique,            -- 用户ID
    nickname           varchar(255) not null default '',        -- 昵称
    avatar             varchar(255) not null default '',        -- 头像
    gender             int not null default 3,                  -- 性别 1: 男 2: 女 3: 未知
    birthday           date not null default '1970-01-01',      -- 生日
    phone              varchar(255) not null default '',        -- 手机号
    email              varchar(255) not null default '',        -- 邮箱
    status             int not null default 1,                  -- 状态 1: 正常 2: 禁用
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 用户身份表
create table member_identity
(
    id                 bigint auto_increment,
    member_uuid        varchar(255) not null,            -- 用户UUID
    identity_type      int not null,                     -- 身份类型 1: 手机号 2: 邮箱 3: 用户名 4: wechat 5: google 6: facebook 7: github
    identifier         varchar(255) not null,            -- 标识符 手机号/邮箱/用户名/google_id/facebook_id/github_id
    credential         varchar(255) not null default '', -- 凭证 密码
    status             int not null default 1,           -- 状态 1: 正常 2: 禁用
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);

create index idx_member_identity_member_uuid_identity_type on member_identity (member_uuid, identity_type);