# Database Migrations

This directory contains SQL migration scripts for the AWSomeShop database.

## Migration Files

Migrations are numbered sequentially:

1. `001_create_users_table.sql` - Creates the users table
2. `002_create_products_table.sql` - Creates the products table
3. `003_create_redemption_orders_table.sql` - Creates the redemption_orders table
4. `004_create_points_transactions_table.sql` - Creates the points_transactions table
5. `005_create_product_price_history_table.sql` - Creates the product_price_history table

## Running Migrations

### Using GORM AutoMigrate (Development)

The application automatically runs migrations on startup using GORM's AutoMigrate feature.

```bash
go run cmd/api/main.go
```

### Using SQL Scripts (Production)

For production environments, you can run the SQL scripts directly:

```bash
mysql -u username -p database_name < migrations/001_create_users_table.sql
mysql -u username -p database_name < migrations/002_create_products_table.sql
mysql -u username -p database_name < migrations/003_create_redemption_orders_table.sql
mysql -u username -p database_name < migrations/004_create_points_transactions_table.sql
mysql -u username -p database_name < migrations/005_create_product_price_history_table.sql
```

Or run all migrations at once:

```bash
cat migrations/*.sql | mysql -u username -p database_name
```

## Adding New Migrations

When adding a new migration:

1. Create a new file with the next sequential number
2. Use descriptive names: `XXX_description.sql`
3. Include both UP and DOWN migrations if needed
4. Test the migration on a development database first
5. Update this README with the new migration

## Notes

- Migrations are idempotent (can be run multiple times safely)
- Always backup your database before running migrations in production
- Foreign key constraints are enforced
- Character set is utf8mb4 for full Unicode support
