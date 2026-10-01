# Game Guess Word System Design

เอกสารนี้อธิบายการออกแบบระบบเกมทายคำศัพท์ภาษาอังกฤษสำหรับ backend โดยอ้างอิง scenario ที่ต้องการและโครงสร้างโปรเจกต์ปัจจุบัน

## Objective

ระบบต้องรองรับให้ client login ได้หลายรูปแบบ จากนั้นดึงคำศัพท์สำหรับทายคำจาก API `/vocabulary/guess-word` โดยเรียงคำศัพท์จากระดับง่ายไปยาก และไม่ส่งคำศัพท์ที่ user เคยทายถูกแล้วกลับไปซ้ำ เมื่อ user ส่งคำตอบ ระบบต้องบันทึกผลการเล่นและคิดคะแนนตามระดับของคำศัพท์

## Authentication

Client สามารถ login ได้ 4 รูปแบบ

1. Guest
2. Email และ password
3. Google
4. Facebook

หลัง login สำเร็จ backend ต้อง return token กลับไปให้ client เพื่อใช้เรียก API ที่ต้องยืนยันตัวตน เช่น `/vocabulary/guess-word`

### Existing Auth Endpoints

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/users/guest` | สร้าง guest user |
| `POST` | `/auth/register` | สมัครด้วย email/password |
| `POST` | `/auth/login` | login ด้วย email/password |
| `GET` | `/auth/google/login` | เริ่ม Google OAuth |
| `GET` | `/auth/google/callback` | รับ callback จาก Google |
| `GET` | `/auth/facebook/login` | เริ่ม Facebook OAuth |
| `GET` | `/auth/facebook/callback` | รับ callback จาก Facebook |

## Main User Flow

1. Client login ด้วย guest, email/password, Google หรือ Facebook
2. Backend ตรวจสอบตัวตนและออก token
3. Client เรียก `GET /vocabulary/guess-word` พร้อม token
4. Backend หา vocabulary ที่ user ยังไม่เคยทายถูก โดยเรียง `vocabulary.level ASC`
5. Backend return response ด้วย model `GameVocabularyResponse`
6. Client แสดงคำใบ้และช่องให้ user ทายคำศัพท์ภาษาอังกฤษ
7. Client submit คำตอบ
8. Backend ตรวจคำตอบ
9. ถ้าถูกต้อง ให้บันทึก `is_correct = true`, เพิ่มคะแนนตาม level และ return `200`
10. ถ้าไม่ถูกต้อง ให้บันทึกหรือ update `is_correct = false` และให้ user ทายซ้ำได้จนกว่าจะถูก
11. ถ้า user เคยเล่น vocabulary นี้แล้ว และมี record ถูกต้องแล้ว ระบบต้องไม่ record ซ้ำ

## Guess Word API

### Endpoint

```http
GET /vocabulary/guess-word
Authorization: Bearer <token>
```

### Business Rules

- ต้องเป็น authenticated user เท่านั้น
- ดึง user id จาก token
- Query vocabulary โดยเรียง `level ASC`
- ถ้า level เท่ากัน ให้เรียง `id ASC` เพื่อให้ผลลัพธ์ deterministic
- ต้องไม่ return vocabulary ที่ user เคยทายถูกแล้ว
- คำว่า "เคยทายแล้ว" ในบริบทของการเลือกคำถัดไปควรหมายถึงเคยมี `game_answers.is_correct = true`
- ถ้า user เคยตอบผิด vocabulary นั้น แต่ยังไม่ถูก สามารถ return คำเดิมให้เล่นต่อได้
- Response ต้องไม่ส่ง `eng` หรือคำตอบจริงกลับไป client

### Response Model

ใช้ model `GameVocabularyResponse` จาก `models/game.go`

```go
type GameVocabularyResponse struct {
	ID            uint   `json:"id"`
	LengthOfWord  int    `json:"length_of_word"`
	PartsOfSpeech string `json:"parts_of_speech"`
	Thai          string `json:"thai"`
	Meaning       string `json:"meaning"`
	Synonyms      string `json:"synonyms"`
}
```

### Success Response

```json
{
  "vocabulary": {
    "id": 1,
    "length_of_word": 5,
    "parts_of_speech": "noun",
    "thai": "บ้าน",
    "meaning": "a building for people to live in",
    "synonyms": "residence, dwelling"
  }
}
```

### No Vocabulary Response

```json
{
  "message": "no vocabulary available",
  "vocabulary": null
}
```

## Submit Answer API

### Recommended Endpoint

```http
POST /vocabulary/guess-word/answer
Authorization: Bearer <token>
Content-Type: application/json
```

### Request Body

```json
{
  "id": 1,
  "word": "house"
}
```

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | number | yes | `vocabulary.id` ที่ user กำลังทาย |
| `word` | string | yes | คำตอบภาษาอังกฤษจาก user |

> หมายเหตุ: โปรเจกต์ปัจจุบันมี `SubmitGameAnswerRequest` อยู่แล้ว โดย field คือ `id` และ `word` สามารถ reuse request model นี้ได้

## Answer Checking Rules

ระบบต้องตรวจคำตอบโดย normalize ก่อนเปรียบเทียบ

- trim whitespace หน้าและหลัง
- เปรียบเทียบแบบ case-insensitive
- รับเฉพาะตัวอักษรภาษาอังกฤษตาม validation ปัจจุบัน

ตัวอย่าง

| Correct Answer | User Answer | Result |
| --- | --- | --- |
| `House` | `house` | correct |
| `house` | ` HOUSE ` | correct |
| `house` | `home` | incorrect |
| `house` | `บ้าน` | invalid |

## Scoring Rules

ถ้า user ตอบถูก ให้เพิ่มคะแนนตาม `vocabulary.level`

| Vocabulary Level | Score |
| --- | ---: |
| `A1` | 1 |
| `A2` | 2 |
| `B1` | 3 |
| `B2` | 4 |
| `C1` | 5 |
| `C2` | 5 |

ถ้า level ไม่อยู่ในรายการที่กำหนด ควร return error หรือ fallback เป็น 0 ตาม policy ที่ทีมเลือก แนะนำให้ return error ฝั่ง service เพื่อไม่ให้ข้อมูล vocabulary ผิด silently

## Record Answer Rules

### ตอบถูก

เมื่อ user submit คำตอบถูก

- return HTTP `200`
- record หรือ update `game_answers.is_correct = true`
- record `user_answer` ล่าสุด
- record `answered_at` ล่าสุด
- เพิ่ม `game_sessions.correct_answers`
- เพิ่ม `game_sessions.score` ตาม level
- vocabulary นี้จะไม่ถูก return จาก `/vocabulary/guess-word` อีกสำหรับ user คนนี้

### ตอบผิด

เมื่อ user submit คำตอบผิด

- record หรือ update `game_answers.is_correct = false`
- record `user_answer` ล่าสุด
- record `answered_at` ล่าสุด
- ไม่เพิ่ม score
- ไม่เพิ่ม correct answer
- user ยังสามารถ submit คำตอบ vocabulary เดิมซ้ำได้
- ถ้า user ตอบผิดหลายครั้ง ให้ update record เดิมเรื่อย ๆ จนกว่า `is_correct = true`

### เคยเล่น vocabulary นี้แล้ว

ถ้า user เคยมี record สำหรับ vocabulary เดียวกันแล้ว

- ถ้า record เดิม `is_correct = false` ให้ update record เดิม
- ถ้า record เดิม `is_correct = true` ห้าม create record ซ้ำ
- ถ้า record เดิม `is_correct = true` และ user submit ซ้ำ ควร return `200` พร้อมสถานะเดิม หรือ return `409 Conflict` แล้วแต่ policy ที่ทีมต้องการ
- แนะนำให้ใช้ `200` แบบ idempotent เพื่อให้ client retry ได้ง่าย

## Data Model

### Current Tables

ระบบปัจจุบันมี table หลักดังนี้

- `users`
- `user_auth_identities`
- `vocabulary`
- `game_sessions`
- `game_answers`

### Recommended Constraint

เพื่อป้องกัน duplicate answer record ต่อ user และ vocabulary ควรทำให้ค้นหา answer เดิมได้จาก `user_id + vocabulary_id`

เนื่องจาก `game_answers` ปัจจุบันผูกกับ `game_session_id` แต่ไม่มี `user_id` โดยตรง จึงมี 2 ทางเลือก

### Option A: Add `user_id` to `game_answers`

เพิ่ม column `user_id` ใน `game_answers`

```sql
ALTER TABLE game_answers
ADD COLUMN user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_game_answers_user_vocabulary_unique
ON game_answers(user_id, vocabulary_id);
```

ข้อดี

- เช็ค duplicate ง่าย
- query `/vocabulary/guess-word` ง่ายและเร็ว
- ตรงกับ rule ที่ว่า user เล่นคำนี้แล้วไม่ควร record ซ้ำ

### Option B: Keep Current Schema and Query Through Session

ใช้ `game_sessions.user_id` เป็นตัวเชื่อม

```sql
CREATE UNIQUE INDEX idx_game_answers_session_vocabulary_unique
ON game_answers(game_session_id, vocabulary_id);
```

ข้อจำกัด

- ป้องกัน duplicate ได้เฉพาะใน session เดียวกัน
- ถ้า user เริ่ม session ใหม่ ยังอาจมี duplicate vocabulary เดิมได้
- ต้อง query ผ่าน `game_sessions` ทุกครั้ง

แนะนำให้ใช้ Option A ถ้า business rule คือ "user หนึ่งคนเล่น vocabulary หนึ่งคำได้สำเร็จเพียงครั้งเดียวตลอดระบบ"

## Repository Design

ควรมี repository methods สำหรับ game answer ดังนี้

```go
FindUserAnswerByVocabularyID(ctx context.Context, userID uint, vocabularyID uint) (*models.GameAnswer, error)
CreateAnswer(ctx context.Context, answer *models.GameAnswer) error
UpdateAnswer(ctx context.Context, answer *models.GameAnswer) error
FindNextVocabularyForUser(ctx context.Context, userID uint) (*models.GameVocabularyResponse, error)
```

Query สำหรับหา vocabulary ถัดไป

```sql
SELECT
    v.id,
    CHAR_LENGTH(v.eng) AS length_of_word,
    v.parts_of_speech,
    v.thai,
    v.meaning,
    v.synonyms
FROM vocabulary AS v
WHERE NOT EXISTS (
    SELECT 1
    FROM game_answers AS ga
    WHERE ga.vocabulary_id = v.id
      AND ga.user_id = :user_id
      AND ga.is_correct = TRUE
)
ORDER BY v.level ASC, v.id ASC
LIMIT 1;
```

ถ้ายังไม่เพิ่ม `user_id` ใน `game_answers` ให้ join ผ่าน `game_sessions`

```sql
SELECT
    v.id,
    CHAR_LENGTH(v.eng) AS length_of_word,
    v.parts_of_speech,
    v.thai,
    v.meaning,
    v.synonyms
FROM vocabulary AS v
WHERE NOT EXISTS (
    SELECT 1
    FROM game_answers AS ga
    JOIN game_sessions AS gs ON gs.id = ga.game_session_id
    WHERE ga.vocabulary_id = v.id
      AND gs.user_id = :user_id
      AND ga.is_correct = TRUE
)
ORDER BY v.level ASC, v.id ASC
LIMIT 1;
```

## Service Design

### Get Vocabulary

```text
VocabularyForGame(userID)
  1. ตรวจว่า user มีอยู่จริง
  2. query vocabulary ที่ user ยังไม่เคยตอบถูก
  3. return GameVocabularyResponse
```

### Submit Answer

```text
SubmitGuessWordAnswer(userID, vocabularyID, word)
  1. validate word
  2. load user
  3. load vocabulary
  4. find existing answer by userID + vocabularyID
  5. ถ้า existing answer is_correct = true
       return success แบบไม่เพิ่มคะแนนซ้ำ
  6. normalize word และ vocabulary.eng
  7. ถ้าตอบถูก
       score = scoreByLevel(vocabulary.level)
       upsert answer is_correct = true
       update session score/correct_answers
       return 200
  8. ถ้าตอบผิด
       upsert answer is_correct = false
       ไม่เพิ่ม score
       return response ให้ client ทายต่อ
```

## Response Design

### Correct Answer Response

```json
{
  "is_correct": true,
  "score_added": 1,
  "message": "correct"
}
```

### Incorrect Answer Response

```json
{
  "is_correct": false,
  "score_added": 0,
  "message": "incorrect"
}
```

แนะนำให้ return `200` ได้ทั้งถูกและผิด เพราะ request ถูกต้องและระบบบันทึกผลสำเร็จ ส่วนความถูกผิดเป็น business result ใน response body แต่ถ้าต้องการยึดตามโค้ดปัจจุบันที่ตอบผิด return `400` ก็ทำได้ เพียงแต่ client ต้อง handle เพิ่ม

## Edge Cases

| Case | Expected Behavior |
| --- | --- |
| token ไม่ถูกต้อง | `401 Unauthorized` |
| user ไม่มีในระบบ | `404 Not Found` |
| vocabulary id ไม่มีในระบบ | `404 Not Found` |
| word ว่าง | `400 Bad Request` |
| word มีอักขระที่ไม่ใช่ภาษาอังกฤษ | `400 Bad Request` |
| user ตอบผิดซ้ำหลายครั้ง | update record เดิม |
| user ตอบถูกหลังจากเคยตอบผิด | update record เดิมเป็น `is_correct = true` และเพิ่ม score |
| user submit คำที่เคยตอบถูกแล้ว | ไม่ create ซ้ำและไม่เพิ่ม score ซ้ำ |
| vocabulary หมดแล้ว | return `vocabulary: null` |

## Implementation Notes For Current Codebase

สิ่งที่โค้ดปัจจุบันมีแล้ว

- `GET /vocabulary/guess-word`
- `GameVocabularyResponse`
- auth middleware สำหรับดึง user id จาก token
- query ที่ไม่ return vocabulary ที่เคยตอบถูกแล้ว
- login ด้วย guest, email/password, Google และ Facebook

สิ่งที่ต้องปรับเพิ่มเพื่อให้ตรง scenario

- เปลี่ยน scoring จากค่าคงที่ `10` เป็นคะแนนตาม `vocabulary.level`
- เปลี่ยนการ submit answer จาก create ใหม่ทุกครั้ง เป็น find existing answer แล้ว update
- เพิ่มการป้องกัน duplicate record เมื่อ user เคยตอบถูก vocabulary เดิมแล้ว
- เพิ่ม unique constraint หรือ index ที่เหมาะสม
- พิจารณาเพิ่ม endpoint submit answer ที่ไม่ต้องพึ่ง `game_session_id` ถ้าเกมนี้ต้องการเล่นแบบต่อคำตาม user โดยตรง

## Acceptance Criteria

- Login ได้ครบทั้ง guest, email/password, Google และ Facebook
- หลัง login สามารถเรียก `/vocabulary/guess-word` ด้วย token ได้
- API `/vocabulary/guess-word` return เฉพาะ field ของ `GameVocabularyResponse`
- API ไม่ส่งคำตอบจริง `eng` กลับไป frontend
- Vocabulary ถูกเลือกโดยเรียง `level ASC, id ASC`
- Vocabulary ที่ user เคยตอบถูกแล้วไม่ถูกส่งกลับไปอีก
- ตอบถูกแล้วได้คะแนนตาม level
- ตอบผิดแล้วบันทึก `is_correct = false`
- ตอบผิดซ้ำแล้ว update record เดิม
- ตอบถูกหลังเคยผิดแล้ว update record เดิมเป็น `is_correct = true`
- ตอบถูกซ้ำ vocabulary เดิมแล้วไม่เพิ่ม record และไม่เพิ่มคะแนนซ้ำ
