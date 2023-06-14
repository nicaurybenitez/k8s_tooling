# k8s_tooling

## PostgreSQL Command Line Script

This command line script allows you to execute SQL commands on a Kubernetes cluster that contains a PostgreSQL container. The script is written in Go and utilizes the `os/exec` package to interact with the command line.

### Prerequisites

Make sure you have Go installed on your system and properly configured. Additionally, you will need access to a Kubernetes cluster that contains a PostgreSQL container and the `kubectl` and `psql` command line tools.

### Usage
Please Generate a encryption Key 
pg_cli generate_key
pg_cli config 
Run the script as follows:

```
go run pg_cli.go [command]
```

Where `[command]` is one of the following:

- `databases`: Show the list of databases.
- `create_db [db_name]`: Create a new database with the specified name.
- `drop_db [db_name]`: Drop an existing database with the specified name.
- `tables [db_name]`: Show the list of tables in the specified database.
- `table_info [db_name] [table_name]`: Show detailed information about a specific table in the specified database.
- `select [db_name] [table_name]`: Execute a SELECT query on the specified table.
- `insert [db_name] [table_name] [values]`: Insert a new row into the specified table with the provided values.
- `update [db_name] [table_name] [set_clause] [condition]`: Update rows in the specified table based on the provided SET clause and condition.
- `delete [db_name] [table_name] [condition]`: Delete rows from the specified table that match the provided condition.

**Note:** You need to provide all the requested data at the beginning of the script, such as the `Namespace`, `Pod Name`, `Database User`, and `Database Password`.

### Example

Here's an example of how to execute the `databases` command:

```
go run pg_cli.go databases
```
This will show the list of databases in the specified Kubernetes cluster and PostgreSQL container.
