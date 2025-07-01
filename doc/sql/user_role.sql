create table public.user_role
(
    id         serial primary key,
    user_id    int not null, -- 用户ID
    role_id    int not null, -- 角色ID
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
