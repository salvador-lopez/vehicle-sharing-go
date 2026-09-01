CREATE DATABASE IF NOT EXISTS inventory;
CREATE DATABASE IF NOT EXISTS inventory_test;
-- Must match the MYSQL_MESSAGE_RELAY_DATABASE default in pkg/Makefile,
-- used by the pkg/message-relay binlog integration test.
CREATE DATABASE IF NOT EXISTS message_relay_test;

CREATE USER IF NOT EXISTS 'inventory'@'%' IDENTIFIED BY 'inventory';

GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `inventory`.* TO 'inventory'@'%';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `inventory_test`.* TO 'inventory'@'%';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `message_relay_test`.* TO 'inventory'@'%';

-- Required so the go-mysql canal client can stream binlog events as a replica.
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'inventory'@'%';