begin;

create table if not exists public.competitors (
    domain_id          bigint      not null,
    domain_name        text        not null,
    competitor_names   text[]      not null,
    competition_levels real[]      not null,
    location           char(2)     not null,
    created_time       date        not null,
    updated_time       date
);

end;
