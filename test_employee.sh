#!/bin/bash
# Exit on error
set -e

# Setup dirs
rm -rf /tmp/mysql_data /tmp/mysql.sock /tmp/mariadb.err /tmp/mariadb.pid
mkdir -p /tmp/mysql_data

# Initialize DB
mariadb-install-db --datadir=/tmp/mysql_data --tmpdir=/tmp >/dev/null 2>&1

# Start mariadbd in the background
mariadbd --datadir=/tmp/mysql_data \
         --tmpdir=/tmp \
         --port=3306 \
         --socket=/tmp/mysql.sock \
         --log-error=/tmp/mariadb.err \
         --pid-file=/tmp/mariadb.pid \
         --skip-networking=0 &

# Wait for MariaDB to start
for i in {1..30}; do
    if mysqladmin --socket=/tmp/mysql.sock ping >/dev/null 2>&1; then
        break
    fi
    sleep 0.5
done

# Create database and user
mariadb --socket=/tmp/mysql.sock -e "CREATE DATABASE IF NOT EXISTS pos_kopitiam;"
mariadb --socket=/tmp/mysql.sock -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'popoi';"
mariadb --socket=/tmp/mysql.sock -e "CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY 'popoi';"
mariadb --socket=/tmp/mysql.sock -e "GRANT ALL PRIVILEGES ON pos_kopitiam.* TO 'root'@'127.0.0.1';"
mariadb --socket=/tmp/mysql.sock -e "FLUSH PRIVILEGES;"

# Apply SQL view migrations
mysql -h 127.0.0.1 -u root -ppopoi pos_kopitiam < migrations/Structure.sql
mysql -h 127.0.0.1 -u root -ppopoi pos_kopitiam < migrations/002_create_analytics_views.sql

echo "DB migration & views ready."

# Seed data
# Wait! seeds/seed.go uses the DB_HOST, DB_USER, DB_PASSWORD from .env
# Let's run it
GOCACHE="/home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf (2)/.go_cache" \
GOPATH="/home/tenzly/Documents/00_GOLang/02_KopiTiam_Project-1/pos-kopitiam-v3nf (2)/.go_path" \
go run seeds/seed.go

echo "Seeding completed."

# Start backend
./main > /tmp/backend.log 2>&1 &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Login as manager
echo "Logging in..."
LOGIN_RES=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@kopitiam.id","password":"password123"}')

echo "Login Response: $LOGIN_RES"

TOKEN=$(echo $LOGIN_RES | grep -oP '"access_token":"\K[^"]+')
BRANCH_ID=$(echo $LOGIN_RES | grep -oP '"branch_id":"\K[^"]+')

echo "Token: $TOKEN"
echo "Branch ID: $BRANCH_ID"

# Get employees
echo "Getting employees..."
curl -s -X GET http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer $TOKEN"

# Create new employee
echo -e "\nCreating new employee..."
CREATE_RES=$(curl -s -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"branch_id\":\"$BRANCH_ID\",\"name\":\"Test Staff\",\"email\":\"teststaff@kopitiam.id\",\"password\":\"password123\",\"pin_code\":\"999999\",\"role\":\"cashier\",\"is_active\":true}")

echo "Create Response: $CREATE_RES"

NEW_EMP_ID=$(echo $CREATE_RES | grep -oP '"id":"\K[^"]+')
echo "New Employee ID: $NEW_EMP_ID"

# Update employee
echo -e "\nUpdating employee..."
UPDATE_RES=$(curl -s -X PUT http://localhost:8080/api/v1/employees/$NEW_EMP_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Staff Updated\",\"email\":\"teststaff@kopitiam.id\",\"role\":\"cashier\",\"is_active\":true}")

echo "Update Response: $UPDATE_RES"

# Clean up
kill $BACKEND_PID
mysqladmin --socket=/tmp/mysql.sock shutdown
