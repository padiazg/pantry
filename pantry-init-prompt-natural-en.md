# Family Pantry API — Go Project with Hexagonal Architecture

You have access to the **Hexago** MCP server, a scaffolding CLI for Go projects with
hexagonal architecture (Ports & Adapters). Use it to generate the complete structure
for the project described below. Work in order and validate the architecture at the end.

---

## Project Context

A backend API to manage a **family pantry**. Products come from purchases at wholesale
stores and are identified by their **EAN-13** barcode. The system must maintain a product
catalog and track their movement (stock entries and exits).

**Target directory:** `/home/pato/go/src/github.com/padiazg/tmp/pantry`

---

## 1. Initialization

Create a new project named `pantry` with Go module `github.com/padiazg/pantry`.
Use **Chi** as the HTTP framework. The project is a standard HTTP server.
Include support for Docker, database migrations, and observability (health checks / metrics).

---

## 2. Domain

### Entities

**Product** — an item stored in the pantry.
Fields: `ean13` (string, natural key), `name`, `description`, `unit` (string — e.g. "kg", "liters", "units"),
`minStock` (float64), `currentStock` (float64), `categoryID` (string), `active` (bool),
`createdAt` and `updatedAt` (time.Time).

**Category** — a grouping of products.
Fields: `id` (string), `name`, `description` (string), `createdAt` (time.Time).

**Movement** — a record of a stock entry or exit.
Fields: `id` (string), `productEan13` (string), `type` (string — see value object),
`quantity` (float64), `reason`, `notes`, `createdBy` (string), `createdAt` (time.Time).

### Value Objects

- **MovementType** — represents the movement type: `"in"` (entry) or `"out"` (exit).
  Must be implemented as typed constants in the domain.
- **StockLevel** — a quantity with positive-value validation (> 0).
- **EAN13** — a 13-digit barcode. Must validate that the check digit is correct
  according to the standard GS1 algorithm.

---

## 3. Business Services

Generate the following services in the business logic layer:

- **ManageProduct** — create, update, and activate/deactivate products.
- **GetProduct** — fetch a product by EAN-13 and list products with filters.
- **ManageCategory** — create and update categories.
- **RecordMovement** — record a stock entry or exit, enforcing business rules.
- **GetMovements** — query movement history with optional filters.
- **GetStockReport** — generate a current stock summary with low-stock alerts.

---

## 4. Adapters

### Output — Database Repositories

Generate a repository for each entity: `ProductRepository`, `CategoryRepository`,
and `MovementRepository`. All are secondary adapters of type database.

### Input — HTTP Handlers

Generate an HTTP handler for each resource: `ProductHandler`, `CategoryHandler`,
and `MovementHandler`. These are primary adapters.

---

## 5. Migrations

Create migrations with the following names, in this order:

1. `create_categories_table`
2. `create_products_table`
3. `create_movements_table`
4. `add_stock_indexes`

The reference schema to implement them is:

```sql
CREATE TABLE categories (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    ean13         CHAR(13) PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    description   TEXT,
    unit          VARCHAR(30) NOT NULL,
    min_stock     DECIMAL(10,3) NOT NULL DEFAULT 0,
    current_stock DECIMAL(10,3) NOT NULL DEFAULT 0,
    category_id   VARCHAR(36) REFERENCES categories(id),
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE movements (
    id            VARCHAR(36) PRIMARY KEY,
    product_ean13 CHAR(13) NOT NULL REFERENCES products(ean13),
    type          VARCHAR(3) NOT NULL CHECK (type IN ('in', 'out')),
    quantity      DECIMAL(10,3) NOT NULL CHECK (quantity > 0),
    reason        VARCHAR(200),
    notes         TEXT,
    created_by    VARCHAR(100),
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movements_product ON movements(product_ean13);
CREATE INDEX idx_movements_created ON movements(created_at);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active   ON products(active);
```

---

## 6. Infrastructure

Generate the following support tools:

- A **validator** named `PantryValidator`.
- A **mapper** named `ProductMapper`.
- A **mapper** named `MovementMapper`.

---

## 7. Validation

Once all components are generated, verify that the architecture is correct
and that there are no dependency violations between layers.

---

## 8. REST API (reference for handler implementation)

### Categories
```
GET    /api/v1/categories        # list all
POST   /api/v1/categories        # create
GET    /api/v1/categories/{id}   # get by ID
PUT    /api/v1/categories/{id}   # update
```

### Products
```
GET    /api/v1/products                    # list (?category=&active=&low_stock=)
POST   /api/v1/products                    # create
GET    /api/v1/products/{ean13}            # get by EAN-13
PUT    /api/v1/products/{ean13}            # update
DELETE /api/v1/products/{ean13}            # deactivate (soft delete)
GET    /api/v1/products/{ean13}/stock      # current stock + minimum level
GET    /api/v1/products/{ean13}/movements  # product movement history
```

### Movements
```
GET    /api/v1/movements         # history (?ean13=&type=&from=&to=)
POST   /api/v1/movements         # record a movement (in/out)
GET    /api/v1/movements/{id}    # movement detail
```

### Reports
```
GET    /api/v1/reports/stock      # current stock summary
GET    /api/v1/reports/low-stock  # products below minimum
```

---

## 9. Business Rules

1. **Non-negative stock**: when recording an exit (`out`), reject if `currentStock - quantity < 0`.
2. **Atomic stock update**: record the `Movement` and update `Product.currentStock`
   in the same database transaction.
3. **Soft delete**: deactivating a product sets `active = false`; the record is never deleted.
4. **Low-stock alert**: a product is low on stock when `currentStock <= minStock`.
5. **EAN-13 validation**: verify that the code has 13 digits and that the check digit
   is correct (GS1 algorithm) before persisting any product or movement.
6. **Movement IDs as UUID v4**: generated in the service layer, not delegated to the database.

---

## 10. Technology Stack

| Component   | Library                                  |
|-------------|------------------------------------------|
| HTTP        | `github.com/go-chi/chi/v5`               |
| Database    | PostgreSQL + `github.com/jmoiron/sqlx`   |
| Migrations  | `github.com/golang-migrate/migrate`      |
| UUID        | `github.com/google/uuid`                 |
| Config      | `github.com/spf13/viper` (already included) |
| Validation  | `github.com/go-playground/validator/v10` |

## 11. Environment Variables

```env
PANTRY_SERVER_PORT=8080
PANTRY_DATABASE_URL=postgres://user:pass@localhost:5432/pantry?sslmode=disable
PANTRY_LOG_LEVEL=info
PANTRY_LOG_FORMAT=json
```

---

When the scaffolding is complete, let me know and we'll start implementing the business
logic layer by layer.
