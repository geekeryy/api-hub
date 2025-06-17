create table public.jwks
(
    id         serial primary key,
    kid        text not null,
    public_key text not null,
    private_key text not null,
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    deleted_at timestamp(6) with time zone
);
