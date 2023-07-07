#! /bin/sh

database_url=$(jq .db.url -r $CONFIG_PATH)
echo "starting migrations in path $MIGRATIONS_PATH for $database_url"
cd $MIGRATIONS_PATH
ls
cat 015_create_prices_table.up.sql
migrate --path $MIGRATIONS_PATH -database $database_url -verbose up