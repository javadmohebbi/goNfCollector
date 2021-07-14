### installation on ubuntu 20
```
sudo apt install postgresql-12 postgresql-client-12 pgadmin3
```

### create user
```
sudo -u postgres createuser --pwprompt oi24_netflow
```


### create db
```
sudo -u postgres pgsql


### Then In postgres shell
CREATE DATABASE oi24_netflow_db

GRANT ALL PRIVILEGES ON DATABASE oi24_netflow_db TO oi24_netflow;

ALTER USER oi24_netflow Superuser;

### EXIT psql using CTRL+D


# EDIT pg_hba
sudo nano /etc/postgresql/12/main/pg_hba.conf

local   all             oi24_netflow                            md5

# restart service
sudo service postgresql restart

# test
psql -U oi24_netflow oi24_netflow_db


```


# INSTALL timescaleDB
```
# Add our PPA
sudo add-apt-repository ppa:timescale/timescaledb-ppa
sudo apt-get update

# Now install appropriate package for PG version
sudo apt install timescaledb-2-postgresql-13



sudo timescaledb-tune


sudo service postgresql restart



# log to psql
psql -U oi24_netflow oi24_netflow_db

# run this command
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

# check extensions
\dx



```





# Remove/Recreate Schema
```
DROP SCHEMA public CASCADE;CREATE SCHEMA public;
```