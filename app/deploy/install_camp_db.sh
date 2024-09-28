#!/bin/bash

# Configure database connection information
DB_USER="root"        # MySQL username
DB_PASSWORD="123456"  # MySQL password (make sure to replace with the actual password)
DB_HOST="localhost"   # MySQL host address
DB_PORT="3306"        # MySQL port, default is 3306
SQL_FILE="v1.0-camp.sql"  # Path to SQL file

# Check if the SQL file exists
if [ ! -f "$SQL_FILE" ]; then
  echo "SQL file not found: $SQL_FILE"
  exit 1
fi

# Execute the SQL file
echo "Executing SQL file $SQL_FILE on the database..."
mysql -u $DB_USER -p$DB_PASSWORD -h $DB_HOST -P $DB_PORT < $SQL_FILE

# Check if execution was successful
if [ $? -eq 0 ]; then
  echo "Database installation successful!"
else
  echo "Database installation failed!"
  exit 1
fi