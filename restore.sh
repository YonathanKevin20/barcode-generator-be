#!/usr/bin/env bash
set -e

# Usage: ./restore.sh <backup_file>
if [ -z "$1" ]; then
  echo "Usage: $0 <backup_file>"
  exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
  echo "Backup file '$BACKUP_FILE' does not exist." 1>&2
  exit 1
fi

# Load environment variables from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# export password for psql commands
export PGPASSWORD="${DB_PASSWORD}"

# Reset schema to clean existing objects
echo "Resetting schema public in database '${DB_NAME}'..."
docker compose exec -T postgres psql -U ${DB_USER:-postgres} -d ${DB_NAME:-barcode_generator} -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Proceed with restore
echo "Restoring database '$DB_NAME' from '$BACKUP_FILE'..."

docker compose exec -T postgres psql -U ${DB_USER:-postgres} -d ${DB_NAME:-barcode_generator} < "$BACKUP_FILE"

echo "Restore completed successfully."
