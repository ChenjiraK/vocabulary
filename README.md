# Vocabulary — Vocab Challenge

แอปฝึกคำศัพท์ภาษาอังกฤษผ่านเกมทายคำจากความหมายภาษาอังกฤษ พร้อมคำใบ้ภาษาไทย บันทึกคำตอบและคะแนนของผู้เล่น โดยแยก Backend และ Frontend ออกจากกัน

## Project structure

```text
vocabulary/
├── be/                 # Go API server
│   ├── migrations/     # SQL migrations
│   ├── docs/           # เอกสาร Backend และ API
│   └── .env.example    # ตัวอย่างการตั้งค่า Backend
├── fe/                 # Next.js frontend
│   ├── app/            # หน้าจอ, components และ API services
│   ├── .nvmrc          # Node.js version
│   └── .env.example    # ตัวอย่างการตั้งค่า Frontend
└── README.md
```

## Backend

### Technology and versions

| Technology             | Version                                        | ใช้สำหรับ                  |
| ---------------------- | ---------------------------------------------- | -------------------------- |
| Go                     | `1.25.0` ตาม `be/go.mod`                       | Backend runtime            |
| Gin                    | `1.12.0`                                       | HTTP framework และ routing |
| GORM                   | `1.31.2`                                       | ORM                        |
| GORM PostgreSQL driver | `1.6.2`                                        | เชื่อมต่อ PostgreSQL       |
| golang-migrate         | `4.19.1`                                       | จัดการ database migrations |
| PostgreSQL             | ยังไม่ได้กำหนด server version ไว้ใน repository | ฐานข้อมูลหลัก              |

เช็คเวอร์ชัน PostgreSQL `psql --version` จะแสดงเพียงเวอร์ชันของ client ตรวจ server ได้ด้วยคำสั่งนี้หลังตั้งค่าฐานข้อมูลแล้ว:

```bash
psql -h localhost -U test -d db_vocabulary -c "SHOW server_version;"
```

### Install Backend

เตรียม Go, PostgreSQL server ที่เปิดใช้งานแล้ว และ PostgreSQL client tools (`psql`, `createuser`, `createdb`) หากต้องการใช้คำสั่ง `make` ให้ติดตั้ง Make ด้วย

จาก root ของ repository:

```bash
cd be
go version
go mod download
cp .env.example .env
```

สำหรับการติดตั้งครั้งแรก ให้สร้าง `.env` จากตัวอย่าง หากมีไฟล์อยู่แล้วให้แก้ค่าของไฟล์เดิม ค่าเริ่มต้นสำหรับ local development คือ:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=test
DB_PASSWORD=password
DB_NAME=db_vocabulary
DB_SSLMODE=disable
APP_PORT=8000
AVATAR_UPLOAD_DIR=./uploads/avatars
PUBLIC_ASSET_BASE_URL=http://localhost:8000
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

| Variable                 | รายละเอียด                                                                                          |
| ------------------------ | --------------------------------------------------------------------------------------------------- |
| `DB_HOST`, `DB_PORT`     | ที่อยู่และ port ของ PostgreSQL                                                                      |
| `DB_USER`, `DB_PASSWORD` | บัญชีที่ใช้เชื่อมต่อฐานข้อมูล                                                                       |
| `DB_NAME`                | ชื่อฐานข้อมูล ต้องสร้างก่อนรัน migration                                                            |
| `DB_SSLMODE`             | รูปแบบ SSL สำหรับเชื่อมต่อฐานข้อมูล ตัวอย่าง local ใช้ `disable`                                    |
| `APP_PORT`               | Port ของ API ตัวอย่างใช้ `8000`; หากไม่ตั้งค่าจะใช้ `8080`                                          |
| `AVATAR_UPLOAD_DIR`      | Directory สำหรับเก็บรูปโปรไฟล์                                                                      |
| `PUBLIC_ASSET_BASE_URL`  | Base URL ที่ Frontend ใช้เปิดรูปโปรไฟล์                                                             |
| `CORS_ALLOWED_ORIGINS`   | Origin ของ Frontend คั่นหลายค่าด้วย comma ได้                                                       |
| `JWT_SECRET`             | Secret สำหรับลงชื่อ token สามารถเพิ่มใน `.env` ได้; หากไม่ตั้งค่าจะใช้ค่าเริ่มต้นสำหรับ development |

สร้าง database user และ database หากยังไม่มี โดยใช้บัญชีผู้ดูแล PostgreSQL ของเครื่อง ตัวอย่างนี้ใช้ `postgres`:

```bash
createuser -h localhost -U postgres --pwprompt test
createdb -h localhost -U postgres -O test db_vocabulary
```

ตั้งรหัสผ่านของ `test` ให้ตรงกับ `DB_PASSWORD` ใน `be/.env` หรือปรับค่า `DB_*` ให้ตรงกับบัญชีและฐานข้อมูลที่มีอยู่แล้ว

### Run migrations

รันจาก directory `be/` หลังสร้าง database และตั้งค่า `.env` เรียบร้อยแล้ว:

```bash
make migrate-up
```

หากไม่ได้ติดตั้ง Make ใช้คำสั่ง Go ได้โดยตรง:

```bash
go run . migrate
```

คำสั่งนี้จะรัน migration ที่ยังไม่ได้ใช้ทั้งหมดจาก `be/migrations/` ปัจจุบันถึง version `3` และสร้างตาราง `vocabulary`, `users`, `user_auth_identities`, `game_sessions`, `game_answers` รวมถึง `schema_migrations` สำหรับเก็บสถานะ migration รันซ้ำได้เมื่อ schema เป็นปัจจุบันแล้ว

ตรวจสอบสถานะและตาราง:

```bash
psql -h localhost -U test -d db_vocabulary -c "SELECT version, dirty FROM schema_migrations;"
psql -h localhost -U test -d db_vocabulary -c '\dt'
```

หากต้องการ rollback จาก version `3` กลับไป version `2` ให้เลือกใช้คำสั่งใดคำสั่งหนึ่ง:

```bash
make migrate-down VERSION=2
# หรือ
go run . migrate down 2
```

`VERSION` คือเวอร์ชันปลายทาง ไม่ใช่จำนวนขั้นตอนที่จะย้อนกลับ การ rollback อาจลบ column, table และข้อมูลตามไฟล์ `.down.sql` ของแต่ละเวอร์ชัน

### Serve Backend

เปิด terminal สำหรับ Backend แล้วรันจาก directory `be/`:

```bash
go run .
```

เมื่อใช้ `.env` ตามตัวอย่าง API จะอยู่ที่ [http://localhost:8000](http://localhost:8000) ตรวจสอบการทำงานจากอีก terminal:

```bash
curl http://localhost:8000/health
```

ผลลัพธ์ที่คาดหวัง:

```json
{ "status": "ok" }
```

การเปิด server จะไม่รัน migration ให้อัตโนมัติ ต้องรัน migration ก่อนใช้งาน API ที่อ่านหรือเขียนข้อมูล

## Frontend

### Technology and versions

| Technology        | Version                                                 | ใช้สำหรับ                        |
| ----------------- | ------------------------------------------------------- | -------------------------------- |
| Node.js           | `22.x` ตาม `fe/.nvmrc`; `package.json` กำหนด `>=22.0.0` | Frontend runtime และ build tools |
| React / React DOM | `19.3.0`                                                | UI library                       |
| Next.js           | `16.3.8`                                                | React framework ใช้ App Router   |
| Material UI       | `9.4.0`                                                 | UI components                    |
| Tailwind CSS      | `4.3.3`                                                 | Styling                          |
| TanStack Query    | `5.104.0`                                               | จัดการ server state              |
| TypeScript        | `6.0.3`                                                 | Type checking                    |

เวอร์ชัน package ในตารางอ้างอิงจาก `fe/package-lock.json` ซึ่งใช้กับ `npm ci` โดย `fe/package.json` กำหนดช่วงเวอร์ชัน React เป็น `^19.2.8` และ Next.js เป็น `^16.3.0`

### Install Frontend

เปิด terminal ใหม่ที่ root ของ repository หากใช้ nvm ให้เลือก Node.js ตาม `.nvmrc`:

```bash
cd fe
nvm install
nvm use
node --version
npm ci
cp .env.example .env.local
```

หากไม่ได้ใช้ nvm ให้ติดตั้ง Node.js `22.x` แล้วข้ามคำสั่ง `nvm install` และ `nvm use` หากมี `.env.local` อยู่แล้วให้แก้ไฟล์เดิม

กำหนดค่าใน `fe/.env.local` ดังนี้ โดยใส่ URL ให้ครบ เพราะ `.env.example` เว้นค่านี้ไว้:

```dotenv
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000
```

### Serve Frontend

รันจาก directory `fe/` และเปิด Backend ควบคู่กันไว้:

```bash
npm run dev
```

เปิดแอปที่ [http://localhost:3000](http://localhost:3000)

สำหรับ build และเปิดด้วย production mode:

```bash
npm run build
npm run start
```

หากเปลี่ยน port ของ Backend ให้แก้ `NEXT_PUBLIC_API_BASE_URL` แล้ว restart Frontend หรือ build ใหม่สำหรับ production หากเปลี่ยน origin ของ Frontend ให้แก้ `CORS_ALLOWED_ORIGINS` และ restart Backend ด้วย

## How to use this app

### 1. นำเข้าคำศัพท์ก่อนเริ่มเล่น

Migration สร้างเฉพาะโครงสร้างฐานข้อมูล ไม่มี seed คำศัพท์ให้ ต้องนำเข้าข้อมูลผ่าน API ก่อน ปัจจุบันยังไม่มีหน้าจอ admin สำหรับอัปโหลด

สร้างไฟล์ Excel นามสกุล `.xlsx` โดยใช้ sheet แรก แถวแรกต้องมีชื่อ column ตามลำดับนี้:

| eng   | parts_of_speech | thai      | meaning                        | synonyms | level |
| ----- | --------------- | --------- | ------------------------------ | -------- | ----- |
| happy | adjective       | มีความสุข | feeling or showing pleasure    | glad     | A1    |
| quick | adjective       | รวดเร็ว   | moving or doing something fast | fast     | A2    |

กรอกให้ครบทั้ง 6 column ทุกแถว เพราะแถวที่มีช่องว่างจะถูกข้าม ใช้ระดับ `A1`, `A2`, `B1`, `B2`, `C1` หรือ `C2` และคำศัพท์ภาษาอังกฤษที่มีเฉพาะตัวอักษร A–Z เพื่อให้เล่นผ่านหน้าจอเกมได้

อัปโหลดเมื่อ Backend เปิดอยู่ โดยเปลี่ยน path ให้เป็นไฟล์จริง:

```bash
curl -X POST http://localhost:8000/vocabularies/upload \
  -F "file=@/absolute/path/vocabulary.xlsx"
```

ตรวจผลลัพธ์ `inserted`, `updated`, `skipped`, `invalid` และ `errors` จาก API โดยคำศัพท์ที่มี `eng` และ `parts_of_speech` ซ้ำกับข้อมูลเดิมจะถูกข้าม

### 2. สมัครสมาชิกและเข้าสู่ระบบ

1. เปิด [http://localhost:3000](http://localhost:3000) เพื่อเข้าหน้า Welcome
2. สมัครบัญชีที่ [http://localhost:3000/register](http://localhost:3000/register) โดยกรอกชื่อที่แสดง, email, password และ confirm password สามารถแนบรูปโปรไฟล์ได้
3. ใช้รหัสผ่านอย่างน้อย 8 ตัวอักษร และไม่เกิน 72 bytes พร้อมยืนยันรหัสผ่านให้ตรงกัน
4. หลังสมัครสำเร็จ ไปที่ [http://localhost:3000/login](http://localhost:3000/login) และเข้าสู่ระบบด้วย email/password เพื่อให้ Frontend เก็บ token สำหรับเรียก API
5. เมื่อ login สำเร็จ ระบบจะเปิดหน้าเกม หรือกลับไปหน้า Welcome แล้วกด **Start** ได้

ปัจจุบันปุ่ม Google/Facebook ในหน้า login ถูกปิดใช้งาน ส่วน Backend มี OAuth endpoints แล้ว การเปิดใช้ต้องตั้งค่า `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` หรือ `FACEBOOK_CLIENT_ID`, `FACEBOOK_CLIENT_SECRET`, `FACEBOOK_REDIRECT_URL` และเชื่อม flow กลับมายัง Frontend เพิ่มเติม โดย `OAUTH_STATE` เป็นค่าที่ Backend รองรับสำหรับตรวจ state

### 3. เล่นเกมทายคำ

1. อ่านความหมายภาษาอังกฤษ, ชนิดของคำ และ synonyms ที่แสดง
2. คลิกช่องคำตอบแล้วพิมพ์คำศัพท์ภาษาอังกฤษให้ครบจำนวนช่อง ระบบจะส่งตรวจคำตอบอัตโนมัติเมื่อพิมพ์ครบ
3. กด **HINT** เพื่อดูอักษรตัวแรก, ตัวสุดท้าย และความหมายภาษาไทย
4. ถ้าตอบผิด ให้แก้คำตอบแล้วลองอีกครั้ง หรือกด Enter เพื่อส่งคำตอบที่กรอกครบแล้ว
5. ถ้าตอบถูก กด **NEXT WORD** เพื่อไปคำถัดไป คำที่ผู้เล่นตอบถูกแล้วจะไม่ถูกเลือกซ้ำ
6. กดปุ่มย้อนกลับด้านบนเพื่อกลับหน้า Welcome

หากพบข้อความ **No vocabulary available.** อาจยังไม่มีคำศัพท์ในฐานข้อมูล หรือผู้เล่นตอบถูกครบทุกคำแล้ว ให้นำเข้าคำศัพท์เพิ่มเพื่อเล่นต่อ

### Current feature status

- หน้า [Vocabulary List](http://localhost:3000/vocabulary-list) ใช้ค้นหาและกรองข้อมูลตัวอย่างจาก `fe/app/data/vocabulary.ts` ยังไม่ได้ดึงรายการคำศัพท์จาก Backend
- ตัวจับเวลาในหน้าเกมยังถูกปิดไว้ แม้หน้า Welcome จะแสดงกติกา 1 นาทีต่อข้อ
- แบบฝึกหัด 20 ข้อสำหรับพิมพ์และหน้าจอ admin ที่อธิบายไว้ใน `document.md` ยังไม่ได้ทำเป็นหน้าจอใช้งาน

## Additional documentation

- [Backend setup](be/docs/setup-golang-gin-gorm-migration.md)
