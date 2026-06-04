#!/bin/bash

migration_script="$1"
container_name="immy-db"

if [[ $migration_script -eq '' ]] then
	echo "You must define the migration filename (or part of its name)"
	exit 1
fi

container_id=$(docker container ls --filter "name=$container_name" --format "{{.ID}}")
migration_script_path=$(ls | grep "$migration_script" | awk '{ print $1 }')

echo "Performing migration of '$migration_script_path' to database with container id $container_id ..."

docker exec -i "$container_id" psql -U admin -d admin < "$migration_script_path"