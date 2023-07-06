#! /bin/sh

database_url=$(jq .db.url -r $CONFIG_PATH)
echo "starting migrations in path $MIGRATIONS_PATH for $database_url"
cd $MIGRATIONS_PATH
cd ..
sql-migrate up -env="production" -config=dbconfig.yml
cd $MIGRATIONS_PATH
ls