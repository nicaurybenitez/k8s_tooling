#!/bin/bash

# Variables
NAMESPACE="postgres"
POD_NAME="postgresql-6999fb7c7c-v5xqp"
DB_NAME="flask-service"
DB_USER="test"
DB_PASSWORD="test@123"

# Función para ejecutar comandos en el pod
exec_command() {
    kubectl exec -it -n "$NAMESPACE" "$POD_NAME" -- bash -c "$1"
}


# Función para ejecutar una consulta SQL en la base de datos
execute_query() {
    query="$1"
    exec_command "psql -U $DB_USER -c \"$query\""
}

# Obtener todas las bases de datos
get_all_databases() {
    execute_query "\l"
}

# Comandos disponibles
help() {
    echo "Uso: ./pg_cli.sh [comando]"
    echo "Comandos disponibles:"
    echo "  databases               Muestra todas las bases de datos"
    echo "  help                    Muestra esta ayuda"
}

# Obtener argumentos del script
command="$1"
shift

# Manejar comandos
case "$command" in
    databases)
        get_all_databases
        ;;
    help)
        help
        ;;
    *)
        help
        ;;
esac
