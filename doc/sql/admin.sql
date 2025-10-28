create table admin_info
(
    id                 bigint auto_increment,
    admin_id          varchar(255) not null unique,            -- 管理员ID
    nickname           varchar(255) not null default '',        -- 昵称
    avatar             varchar(255) not null default '',        -- 头像
    gender             int not null default 3,                  -- 性别 1: 男 2: 女 3: 未知
    birthday           date not null default '1970-01-01',      -- 生日
    phone              varchar(255) not null default '',        -- 手机号
    email              varchar(255) not null default '',        -- 邮箱
    status             int not null default 1,                  -- 状态 1: 正常 2: 禁用
    created_at datetime default current_timestamp not null,     -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);