#!/usr/bin/env bash
set -e

# Load environment variables from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Ensure backups directory exists (absolute path)
mkdir -p /backups
# Remove backups older than 30 days
find /backups -type f -name "*.sql" -mtime +30 -delete

# Timestamp for backup file
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
# Use absolute backup directory
BACKUP_FILE="/backups/${DB_NAME:-barcode_generator}_$TIMESTAMP.sql"

echo "Backing up database '$DB_NAME' to '$BACKUP_FILE'..."

# pass password via env and dump over the network
export PGPASSWORD="${DB_PASSWORD}"
pg_dump -h postgres -U "${DB_USER:-postgres}" -d "${DB_NAME:-barcode_generator}" > "$BACKUP_FILE"

# Ensure the backup file is owned by non-root user
chown ${UID:-1000}:${GID:-1000} "$BACKUP_FILE"

echo "Backup completed successfully."
