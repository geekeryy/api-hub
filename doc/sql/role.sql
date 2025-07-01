create table public.role
(
    id         serial primary key,
    name       varchar(255) not null, -- 角色名称
    description varchar(255) not null default '', -- 角色描述
    status     int not null, -- 状态 1: 正常 2: 禁用
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
