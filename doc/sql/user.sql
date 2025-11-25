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



-- jwks密钥表
create table jwks
(
    id         bigint auto_increment,
    service    varchar(255) not null, -- 服务名称
    kid        varchar(255) not null, -- 密钥ID
    public_key text not null, -- 公钥
    private_key text not null, -- 私钥
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 刷新令牌表
create table refresh_token
(
    id                   bigint auto_increment,
    refresh_token_hash   text(3000) not null ,         -- 刷新令牌哈希
    member_id            varchar(255) not null,                -- 用户ID
    status               int not null,                         -- 状态 1: 正常 2: 禁用
    expired_at           datetime not null, -- 过期时间
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 令牌刷新记录表
create table token_refresh_record
(
    id                   bigint auto_increment,
    refresh_token_hash   text(3000) not null,                        -- 刷新令牌哈希
    token                text(3000) not null ,                 -- 令牌
    kid                  varchar(255) not null,                        -- 密钥ID
    ip                   varchar(255) not null default '',     -- IP地址
    user_agent           varchar(255) not null default '',     -- 用户代理
    expired_at           datetime not null, -- 过期时间
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 菜单表
create table menus
(
    id         bigint auto_increment,
    name       varchar(255) not null, -- 菜单名称
    path       varchar(255) not null, -- 菜单路径 用于前端鉴权
    icon       varchar(255) not null default '', -- 菜单图标
    sort       int not null default 0, -- 排序优先级 用于菜单排序
    pid        bigint not null, -- 父级ID
    description varchar(255) not null default '', -- 菜单描述
    status     int not null default 1, -- 状态 1: 正常 2: 禁用
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 权限表
create table permissions
(
    id         bigint auto_increment,
    name       varchar(255) not null, -- 权限名称
    type       int not null, -- 权限类型 1: 菜单 2: 按钮 3: 接口
    params     json not null, -- 权限参数 不同类型对应不同的权限参数
    apis       json not null, -- 接口列表 用于后端鉴权
    pid  bigint not null, -- 父级ID
    description varchar(255) not null default '', -- 权限描述
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 角色表
create table roles
(
    id         bigint auto_increment,
    name       varchar(255) not null, -- 角色名称
    description varchar(255) not null default '', -- 角色描述
    status     int not null, -- 状态 1: 正常 2: 禁用
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);


-- 角色权限表
create table role_permissions
(
    id         bigint auto_increment,
    role_id    bigint not null, -- 角色ID
    permission_id bigint not null, -- 权限ID
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);

-- 用户角色表
create table user_roles
(
    id         bigint auto_increment,
    user_id    bigint not null, -- 用户ID
    role_id    bigint not null, -- 角色ID
    created_at datetime default current_timestamp not null, -- 创建时间
    updated_at datetime default current_timestamp on update current_timestamp not null, -- 更新时间
    primary key (id)
);
