create table public.role_permission
(
    id         serial primary key,
    role_id    int not null, -- 角色ID
    permission_id int not null, -- 权限ID
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
