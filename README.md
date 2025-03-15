# GO CMS Service

A lightweight CMS service built with Go and GORM, using SQLite to manage dynamic tables from JSON schemas. Designed for development with schema flexibility and data query/CRUD API support.

## ✨ Features
- Dynamic table creation from JSON schemas.
- Query data from schema-generated tables.
- CRUD APIs for each table.


## 🛠️ Tech Stack
- Go
- GORM
- SQLite (3.35.0+ required for DROP COLUMN)

## 🚀 Run
`go run cmd/main.go`

### 📦 Definition of Value Types 
| Schema Value Type | SQLite Type |
|-------------------|------------|
| `STRING`          | `TEXT`     |
| `INT`             | `INTEGER`  |
| `DOUBLE`          | `REAL`     |
| `FLOAT`           | `REAL`     |
| `BOOLEAN`         | `INTEGER` (`0 = false, 1 = true`) |

## 📌 Notes
- Drops and recreates columns on type changes (data not preserved, dev-only).
- Ensure SQLite version supports DROP COLUMN.
