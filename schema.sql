drop table if exists properties;

create table properties
(
    id              bigserial primary key,
    parcel          text not null default '',
    address         text      not null,
    city            text      not null,
    coordinates     point not null default '(0,0)',
    lotsqft         integer not null default 0,
    sqft            integer not null default 0,
    state           text  not null default '',
    zipcode         text  not null default '',
    use_code        text not null default '',
    total_rooms     integer not null default 0,
    basement        text not null default '',
    style           text not null default '',
    bedrooms        integer not null default 0,
    grade           text not null default '',
    stories         integer not null default 0,
    full_baths      integer not null default 0,
    half_baths      integer not null default 0,
    condition       text not null default '',
    year_built      integer not null default 0,
    fireplaces      integer not null default 0,
    exterior_finish text not null default '',
    heating_cooling text not null default '',
    basement_garage integer not null default 0,
    roof_type       text not null default '',

    created_at      timestamp not null default now(),
    updated_at      timestamp not null default now()
);