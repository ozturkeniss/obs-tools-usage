#!/bin/bash

# Database management script
# This script handles database operations like migrations, seeding, and backups

set -e

echo "🗄️  Database Management Script..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Navigate to project root
cd "$(dirname "$0")/.."

# Default values
ACTION="migrate"
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="product_service"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        migrate)
            ACTION="migrate"
            shift
            ;;
        seed)
            ACTION="seed"
            shift
            ;;
        backup)
            ACTION="backup"
            shift
            ;;
        restore)
            ACTION="restore"
            shift
            ;;
        reset)
            ACTION="reset"
            shift
            ;;
        status)
            ACTION="status"
            shift
            ;;
        --host)
            DB_HOST="$2"
            shift 2
            ;;
        --port)
            DB_PORT="$2"
            shift 2
            ;;
        --user)
            DB_USER="$2"
            shift 2
            ;;
        --password)
            DB_PASSWORD="$2"
            shift 2
            ;;
        --database)
            DB_NAME="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [ACTION] [OPTIONS]"
            echo ""
            echo "ACTIONS:"
            echo "  migrate   Run database migrations (default)"
            echo "  seed      Seed database with initial data"
            echo "  backup    Create database backup"
            echo "  restore   Restore database from backup"
            echo "  reset     Reset database (drop and recreate)"
            echo "  status    Show database status"
            echo ""
            echo "OPTIONS:"
            echo "  --host HOST        Database host (default: localhost)"
            echo "  --port PORT        Database port (default: 5432)"
            echo "  --user USER        Database user (default: postgres)"
            echo "  --password PASS    Database password (default: password)"
            echo "  --database DB      Database name (default: product_service)"
            echo "  -h, --help         Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Set environment variables
export DB_HOST=$DB_HOST
export DB_PORT=$DB_PORT
export DB_USER=$DB_USER
export DB_PASSWORD=$DB_PASSWORD
export DB_NAME=$DB_NAME
export DB_SSL_MODE=disable

# Function to check if PostgreSQL is running
check_postgres() {
    print_status "Checking PostgreSQL connection..."
    if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
        print_error "PostgreSQL is not running or not accessible!"
        print_status "Please start PostgreSQL or check your connection settings."
        exit 1
    fi
    print_success "PostgreSQL is running!"
}

# Function to run migrations
run_migrations() {
    print_status "Running database migrations..."
    
    # Build the application
    go build -o bin/product-service cmd/product/main.go
    
    # Run migrations (this would be implemented in the application)
    print_status "Migrations would be run here..."
    print_warning "Migration functionality needs to be implemented in the application"
}

# Function to seed database
seed_database() {
    print_status "Seeding database with initial data..."
    
    # Build the application
    go build -o bin/product-service cmd/product/main.go
    
    # Run seeding (this would be implemented in the application)
    print_status "Database seeding would be run here..."
    print_warning "Seeding functionality needs to be implemented in the application"
}

# Function to create backup
create_backup() {
    print_status "Creating database backup..."
    
    BACKUP_FILE="backups/product_service_$(date +%Y%m%d_%H%M%S).sql"
    mkdir -p backups
    
    pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME > $BACKUP_FILE
    
    if [ $? -eq 0 ]; then
        print_success "Backup created: $BACKUP_FILE"
        ls -lh $BACKUP_FILE
    else
        print_error "Failed to create backup!"
        exit 1
    fi
}

# Function to restore database
restore_database() {
    print_status "Restoring database from backup..."
    
    # Find the latest backup
    LATEST_BACKUP=$(ls -t backups/product_service_*.sql 2>/dev/null | head -n1)
    
    if [ -z "$LATEST_BACKUP" ]; then
        print_error "No backup files found in backups/ directory!"
        exit 1
    fi
    
    print_status "Using backup: $LATEST_BACKUP"
    
    # Restore the database
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME < $LATEST_BACKUP
    
    if [ $? -eq 0 ]; then
        print_success "Database restored successfully!"
    else
        print_error "Failed to restore database!"
        exit 1
    fi
}

# Function to reset database
reset_database() {
    print_warning "This will DROP and RECREATE the database!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Operation cancelled."
        exit 0
    fi
    
    print_status "Dropping database..."
    dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || true
    
    print_status "Creating database..."
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    
    if [ $? -eq 0 ]; then
        print_success "Database reset successfully!"
    else
        print_error "Failed to reset database!"
        exit 1
    fi
}

# Function to show database status
show_status() {
    print_status "Database Status:"
    echo "  Host: $DB_HOST"
    echo "  Port: $DB_PORT"
    echo "  User: $DB_USER"
    echo "  Database: $DB_NAME"
    echo ""
    
    # Check connection
    if pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
        print_success "Connection: OK"
    else
        print_error "Connection: FAILED"
        exit 1
    fi
    
    # Show database size
    DB_SIZE=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));" 2>/dev/null | xargs)
    if [ ! -z "$DB_SIZE" ]; then
        echo "  Size: $DB_SIZE"
    fi
    
    # Show table count
    TABLE_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
    if [ ! -z "$TABLE_COUNT" ]; then
        echo "  Tables: $TABLE_COUNT"
    fi
}

# Main execution
case $ACTION in
    migrate)
        check_postgres
        run_migrations
        ;;
    seed)
        check_postgres
        seed_database
        ;;
    backup)
        check_postgres
        create_backup
        ;;
    restore)
        check_postgres
        restore_database
        ;;
    reset)
        check_postgres
        reset_database
        ;;
    status)
        show_status
        ;;
esac
