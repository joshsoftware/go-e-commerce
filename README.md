Assumption 1: You are simran.
Assumption 2: Words in CAPITAL_CASE are meant to be set by simran.
Assumption 3: Port 33001 must be free, else specify other port in application.yml .
Assumption 4: Postgresql server is running on port 5432, if on other port then makes appropriate changes in application.yml. 
Assumption 5: All the source code is pulled from git and you are currently on product_api branch

## Setting Up Database
Use PostgresQL database with version greater than 12.


Create a database in PostgresQL via commandline
sudo -u postgres createdb --owner=USERNAME DATABASE_NAME
example:
```
sudo -u postgres createdb --owner=simran Commerce
```

```
[Just For Knowledge]
Drop a database in PostgresQL via commandline
dropdb -h localhost -p PORT_NUMBER -U USERNAME DATABASE_NAME
example:
dropdb -h localhost -p 5432 -U simran Commerce
```

Copy your DB_URI to application.yml file
DB_URI: "postgresql://USERNAME:PASSWORD@localhost:PORT_NUMBER/DATABASE_NAME?sslmode=disable"
example: 
```
"postgresql://simran:root@localhost:5432/Commerce?sslmode=disable"
```
### golang version
golang version 1.14.10

### Build the Product_api
```
go build
```

You should see some executable getting created with name go-e-commerce

### Migrate
```
./go-e-commerce migrate
```

This will create tables in the DATABASE_NAME you specified in application.yml

### Copying Dump to Database
Copy the migration.sql dump to this db with following command.
psql DATABASE_NAME < PATH_TO_migration.sql
example : 
```
psql Commerce < /home/simran/Desktop/Josh/InternProject/go-e-commerce/migration.sql
```

You should see that all records have been inserted successfully, if not either you messed up or there is version problem or our code is broke.


### Run directly on host
```
./go-e-commerce start
```

### Run on Docker

Dependency 1: You must have docker installed with sudo rights.
Dependency 2: You must have all migrations already created(See Migrate).

## First Way
Build Locally and run
```
docker build -t joshsoftware/go-e-commerce_product-api:v1 .


docker run -it -p 33001:33001  \
 --network=host  \
 joshsoftware/go-e-commerce_product-api:v1
```

## Second Way
Directly Pull the docker image from https://hub.docker.com/repository/docker/skavhar1998/go-e-commerce_product-api and then run.

```
docker pull skavhar1998/go-e-commerce_product-api:v1

docker run -it -p 33001:33001  \
 --network=host  \
 skavhar1998/go-e-commerce_product-api:v1
```

### For documentation, refer the docs folder
docs only include Product_api documentation