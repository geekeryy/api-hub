create table public.refresh_token
(
    id                   serial primary key,
    refresh_token_hash   text not null unique,                 -- 刷新令牌哈希
    member_id            varchar(255) not null,                -- 用户ID
    status               int not null,                         -- 状态 1: 正常 2: 禁用
    expired_at           timestamp(6) with time zone not null, -- 过期时间
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
