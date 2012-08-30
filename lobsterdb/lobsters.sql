create role golobs 
        login 
        password '';
create database lobsters 
                ENCODING = 'UTF8' 
                LC_COLLATE = 'en_US.UTF-8' 
                LC_CTYPE = 'en_US.UTF-8' 
                template = template0;
ALTER DATABASE lobsters OWNER TO golobs ;
\connect lobsters
create table posted (guid text unique not null, 
                     posted boolean default false not null);
create index posted_idx on posted;
