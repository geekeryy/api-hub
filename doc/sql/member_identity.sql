create table public.member_identity
(
    id                 serial primary key,
    member_id          varchar(255) not null,            -- 用户ID
    identity_type      int not null,                     -- 身份类型 1: 手机号 2: 邮箱 3: 用户名 4: wechat 5: google 6: facebook 7: github
    identifier         varchar(255) not null,            -- 标识符 手机号/邮箱/用户名/google_id/facebook_id/github_id
    credential         varchar(255) not null default '', -- 凭证 密码
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
