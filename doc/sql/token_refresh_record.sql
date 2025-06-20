create table public.token_refresh_record
(
    id                   serial primary key,
    refresh_token_hash   text not null,                        -- 刷新令牌哈希
    token                text not null unique,                 -- 令牌
    kid                  text not null,                        -- 密钥ID
    ip                   varchar(255) not null,                -- IP地址
    expired_at           timestamp(6) with time zone not null, -- 过期时间
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
