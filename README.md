Note : words in CAPITAL_CASE are meant to be set by the user.
## Setting Up Database
Use PostgresQL database with version greater than 12.


Create a database in PostgresQL via commandline
sudo -u postgres createdb --owner=USERNAME DATABASE_NAME
example:
```
sudo -u postgres createdb --owner=santosh Commerce
```

```
[Just For Knowledge]
Drop a database in PostgresQL via commandline
dropdb -h localhost -p 5432 -U USERNAME DATABASE_NAME
example:
dropdb -h localhost -p 5432 -U santosh Commerce
```

Copy your DB_URI to application.yml file
DB_URI: "postgresql://USERNAME:PASSWORD@localhost:5432/DATABASE_NAME?sslmode=disable"
example: 
```
"postgresql://santosh:root@localhost:5432/Commerce?sslmode=disable"
```


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
psql Commerce < /home/santosh/Desktop/Josh/InternProject/go-e-commerce/migration.sql
```

You should see that all records have been inserted successfully, if not either you messed up or there is version problem or our code is broke.


### Run
```
./go-e-commerce start
```

### For documentation, refer the docs folder
docs only include Product_api documentation