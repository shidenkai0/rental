# Rental Service API

## Introduction

## Database migrations

### Create a migration

To create a migration called `$migration_name` run the following command:

```
migrate create -ext sql -dir db/migrations -seq $migration_name
```

This should create files called `$migration_name.up.sql` and `$migration_name.down.sql` in the `db/migrations` directory. Fill in the necessary SQL statements to create the tables and columns you need. Do not forget to fill the "down" migration with the SQL statements to drop the tables and columns you created: running the "up" and "down" files sequentially should result in a no-op.
