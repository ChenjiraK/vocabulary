# Setup Golang + Gin + GORM + Database Migration

เอกสารนี้อธิบายขั้นตอนติดตั้งและตั้งค่าโปรเจกต์ Go ที่ใช้ Gin, GORM, PostgreSQL และ migration ด้วย `golang-migrate`

## 1. ติดตั้ง Go

ตรวจสอบว่าเครื่องมี Go แล้ว:

```bash
go version
```

ถ้ายังไม่มี ให้ติดตั้งจาก:

```text
https://go.dev/doc/install
```

## 2. สร้าง Go Module

ถ้าโปรเจกต์ยังไม่มี `go.mod` ให้รัน:

```bash
go mod init vocabulary
```

โปรเจกต์นี้ใช้ module name:

```text
vocabulary
```

## 3. ติดตั้ง Gin + GORM

ติดตั้ง Gin:

```bash
go get github.com/gin-gonic/gin
```

ติดตั้ง GORM:

```bash
go get gorm.io/gorm
```

ติดตั้ง PostgreSQL driver สำหรับ GORM:

```bash
go get gorm.io/driver/postgres
```

ติดตั้งตัวอ่านไฟล์ `.env`:

```bash
go get github.com/joho/godotenv
```

ติดตั้ง migration library:

```bash
go get github.com/golang-migrate/migrate/v4
```

จัด dependency:

```bash
go mod tidy
```

## 4. ตั้งค่า Environment

สร้างไฟล์ `.env` ที่ root project:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=test
DB_PASSWORD=password
DB_NAME=db_vocabulary
DB_SSLMODE=disable
APP_PORT=8000
```

ความหมายของค่า:

- `DB_HOST`: host ของ database
- `DB_PORT`: port ของ PostgreSQL
- `DB_USER`: username สำหรับต่อ database
- `DB_PASSWORD`: password สำหรับต่อ database
- `DB_NAME`: ชื่อ database
- `DB_SSLMODE`: local development ใช้ `disable`
- `APP_PORT`: port สำหรับ run API server

## 5. สร้าง Database

ถ้ายังไม่มี database ให้สร้างก่อน:

```bash
createdb -U test db_vocabulary
```

ถ้า user PostgreSQL เป็นคนละชื่อ ให้เปลี่ยน `test` ให้ตรงกับเครื่องของคุณ

## 6. โครงสร้าง Migration File

โปรเจกต์ใช้ `golang-migrate` ดังนั้นไฟล์ migration ต้องแยกเป็น `.up.sql` และ `.down.sql`

ตัวอย่าง version 1:

```text
migrations/000001_create_vocabulary_table.up.sql
migrations/000001_create_vocabulary_table.down.sql
```

ไฟล์ `up` ใช้สร้าง table:

```sql
CREATE TABLE IF NOT EXISTS vocabulary (
    id BIGSERIAL PRIMARY KEY,
    eng VARCHAR(255) NOT NULL,
    parts_of_speech VARCHAR(100) NOT NULL,
    thai VARCHAR(255) NOT NULL,
    meaning TEXT,
    synonyms TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

ไฟล์ `down` ใช้ rollback:

```sql
DROP TABLE IF EXISTS vocabulary;
```

## 7. Makefile สำหรับ Migration

โปรเจกต์นี้ตั้งค่า migration command ไว้ใน `Makefile`:

```makefile
GO=go

.PHONY: migrate-up migrate-down

migrate-up:
	$(GO) run . migrate

migrate-down:
	@test -n "$(VERSION)" || (echo "VERSION is required. Usage: make migrate-down VERSION=2" && exit 1)
	$(GO) run . migrate down $(VERSION)
```

## 8. Run Migration ครั้งแรก

ถ้า database ยังว่างเปล่า และยังไม่มี table `schema_migrations` ให้รัน:

```bash
make migrate-up
```

`golang-migrate` จะสร้าง table `schema_migrations` ให้อัตโนมัติ และรัน migration version 1

หลังรันสำเร็จ database จะมี table:

```text
schema_migrations
vocabulary
```

ตรวจสอบ table:

```bash
psql -U test -d db_vocabulary -c "\dt"
```

## 9. Rollback Migration

Rollback กลับไปยัง version ที่ต้องการ:

```bash
make migrate-down VERSION=0
```

กรณีนี้หมายถึง rollback จาก version 1 กลับไป version 0 และจะลบ table `vocabulary`

ถ้ามี migration version 2 แล้ว และต้องการ rollback ให้ database อยู่ที่ version 1:

```bash
make migrate-down VERSION=1
```

ถ้าไม่ส่ง `VERSION` จะได้ error:

```bash
VERSION is required. Usage: make migrate-down VERSION=2
```

## 10. Run Project

รัน API server:

```bash
go run .
```

จากค่า `.env` ปัจจุบัน server จะรันที่ port `8000`

ทดสอบ health check:

```bash
curl http://localhost:8000/health
```

ผลลัพธ์ที่ควรได้:

```json
{"status":"ok"}
```

## 11. สรุปคำสั่งที่ใช้บ่อย

ติดตั้ง dependency:

```bash
go mod tidy
```

รัน migration:

```bash
make migrate-up
```

rollback migration:

```bash
make migrate-down VERSION=0
```

run project:

```bash
go run .
```

build project:

```bash
go build ./...
```
