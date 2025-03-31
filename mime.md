# MIME Engine Specification

## Overview
MIME is a declarative DSL-driven system that enables defining entities, routes, and validation rules in a structured manner. It functions as an SQLite-backed execution engine, allowing users to define and interact with data models via a REPL or predefined routes.

## Architecture
### 1. Core Data Handling
- **Schema Storage & Management**
  - Store entity definitions, routes, and validation rules in memory.
  - Load schema from the DSL or an existing SQLite database.
  
- **SQLite Abstraction Layer**
  - Interface for executing queries against SQLite.
  - Support for CRUD operations.
  - Manage schema migrations and transactions.

### 2. Execution Layer
- **Validation System**
  - Ensure payloads match expected types.
  - Enforce constraints (e.g., `number {increment}`, `text <> ["male", "female"]`).
  - Support custom validation rules.

- **Routing & Execution**
  - Map routes to execution paths (e.g., `POST /employees` → insert payload).
  - Handle data lookups (`GET /employees/:id` → return record).
  - Format responses (e.g., JSON output).

- **REPL Integration**
  - Accept and execute raw commands interactively.
  - Support for basic operations (`get /employees/1`, `restart`).
  - Provide validation feedback.

### 3. Persistence & Migrations
- **Schema Diffing**
  - Compare stored schema vs. latest declaration.
  - Generate safe migration plans.
  - `restart` should confirm changes before applying migrations.

- **Backups & Recovery**
  - Automatically create backups before schema modifications.
  - Provide rollback mechanisms in case of failures.

### 4. Extensibility & Future-Proofing
- **Plugin System (Future Consideration)**
  - Enable custom validation rules.
  - Potential support for alternative storage backends.

- **Exporting & API Integration**
  - Support `mime serve` for running as a local API.
  - JSON/HTTP response handling.

## Example DSL
### Entity Definition
```plaintext
entity student ->
    id: number {increment}
    dob: text
    age: number
    created_at: timestamp
    gender: text <> ["male", "female"]
```

### Route Definition
```plaintext
routes ->
    GET /students/:id -> self.id
    POST /students -> payload
    POST /students -> self
```

### Alter Payload
```plaintext
alter ref student.payload ->
    gender: text
    age: number
    dob: text
```

## Commands
- `restart` → Reload schema and migrate database after confirmation.
- `get <endpoint>` → Retrieve data.
- `post <endpoint>` → Insert data.
- `validate` → Check schema consistency without modifying the database.

## Notes
- The DSL should remain minimal yet expressive.
- The REPL will allow on-the-fly interaction without requiring full restarts.
- The SQLite layer should ensure atomicity and safe migrations.

