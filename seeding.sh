#!/bin/bash
set -e

for file in ./seeders/*.sql; do
  echo "Running seed: $file"
  docker compose exec -T postgres psql -U postgres -d barcode_generator < "$file"
done
