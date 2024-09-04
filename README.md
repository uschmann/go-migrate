# dbmigrate

Create and execute migrations for oracle db using sql plus.

## Usage

```bash
NAME:
   dbmigrate - Create and execute migrations for oracle db

USAGE:
   dbmigrate [global options] command [command options]

COMMANDS:
   status    List pending and executed migrations
   migrate   Execute all pending migrations
   rollback  Rollback the last batch of migrations
   make      Create a new migration
   test      Test configuration and print results
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value, -d value  Path to the folder that contains the migrations (default: "./sql")
   --help, -h             show help
```

## Configuration

Configuration is done using OS environment variables.

**dbmigrate** tries to read a .env file from the current working direcoty when started.

```bash
DB_USER=auschmann
DB_PASSWORD=secret
DB_HOST=localhost
DB_PORT=1522
DB_SERVICE=FREE
DB_MIGRATION_LOG_TABLE=auschmann.migration_logs
SQLPLUS_BIN=sqlplus
```