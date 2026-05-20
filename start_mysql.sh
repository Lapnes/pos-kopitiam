#!/bin/bash
set -e

# Setup dirs
rm -rf /tmp/mysql_data /tmp/mysql.sock /tmp/mariadb.err /tmp/mariadb.pid
mkdir -p /tmp/mysql_data

# Initialize DB
mariadb-install-db --datadir=/tmp/mysql_data --tmpdir=/tmp

# Start mariadbd in the background with nohup
nohup mariadbd --datadir=/tmp/mysql_data \
               --tmpdir=/tmp \
               --port=3306 \
               --socket=/tmp/mysql.sock \
               --log-error=/tmp/mariadb.err \
               --pid-file=/tmp/mariadb.pid \
               --skip-networking=0 >/dev/null 2>&1 &

# Disown the process so it persists
disown

# Wait for MariaDB to start
echo "Waiting for MariaDB to start..."
for i in {1..30}; do
    if mysqladmin --socket=/tmp/mysql.sock ping >/dev/null 2>&1; then
        echo "MariaDB is up!"
        break
    fi
    sleep 0.5
done

if ! mysqladmin --socket=/tmp/mysql.sock ping >/dev/null 2>&1; then
    echo "MariaDB failed to start. Logs:"
    cat /tmp/mariadb.err
    exit 1
fi

# Create database and grant privileges
mariadb --socket=/tmp/mysql.sock -e "CREATE DATABASE IF NOT EXISTS pos_kopitiam;"
mariadb --socket=/tmp/mysql.sock -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'popoi';"
mariadb --socket=/tmp/mysql.sock -e "CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY 'popoi';"
mariadb --socket=/tmp/mysql.sock -e "GRANT ALL PRIVILEGES ON pos_kopitiam.* TO 'root'@'127.0.0.1';"
mariadb --socket=/tmp/mysql.sock -e "FLUSH PRIVILEGES;"

echo "Database initialized successfully!"
