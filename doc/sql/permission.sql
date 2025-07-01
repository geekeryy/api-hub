create table public.permission
(
    id         serial primary key,
    name       varchar(255) not null, -- 权限名称
    type       int not null, -- 权限类型 1: 菜单 2: 按钮 3: 接口
    path       varchar(255) not null, -- 权限路径 用于前端鉴权
    apis       jsonb not null, -- 接口列表 用于后端鉴权
    parent_id  int not null, -- 父级ID
    sort       int not null, -- 排序 用于菜单排序
    description varchar(255) not null default '', -- 权限描述
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
