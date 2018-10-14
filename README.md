#Using Golang to explose ltree function in postgres db

Setup db: 
1. docker network create --driver bridge postgres-network
2. docker run --name postgres1 --network postgres-network -it -p 5432:5432 -e POSTGRES_DB=postgres1 -e DB_HOST=postgres1 -e POSTGRES_PASSWORD=postgres -v /Users/andreaswestberg/work/src/github.com/skiarn/ltree/pgdata:/var/lib/postgressql/data -d postgres
3. docker exec -it containerid /bin/bash
4. createuser testuser -P --createdb -h postgres1 -U postgres
5. psql -h postgres1 -U postgres
6. create extension ltree;
7.create table hierarchy (
    id serial primary key,
    nodeid uuid,
    path ltree
);
8. create index hierarchy_path_idx on hierarchy using gist (path);
9. GRANT SELECT, INSERT, UPDATE, DELETE ON hierarchy TO testuser;
10. GRANT USAGE, SELECT ON SEQUENCE hierarchy_id_seq TO testuser;

