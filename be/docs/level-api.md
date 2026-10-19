# GET /level

Requires `Authorization: Bearer <token>`. The token identifies the user whose
progress is returned. No query parameters or request body are needed.

Returns HTTP 200 with all six levels:

```json
{
  "A1": "inactive",
  "A2": "inactive",
  "B1": "complete",
  "B2": "active",
  "C1": "inactive",
  "C2": "inactive"
}
```

- `inactive`: No vocabulary exists for the level. (ไม่มีคำศัพท์ใน level)
- `complete`: The user has answered every vocabulary in the level correctly.(ผู้ใช้ตอบถูกครบทุกคำใน level)
- `active`: Vocabulary exists, but at least one word has not been answered correctly
  by the user (including levels the user has never played). (มีคำศัพท์ แต่ยังตอบถูกไม่ครบ รวมถึงยังไม่ได้เริ่มเล่น)

Progress uses `game_answers.user_id`, `vocabulary_id`, and `is_correct = true`,
consistent with the guess-word vocabulary selection. Each word counts once even
if multiple correct answers exist across sessions. Other users' answers and
incorrect answers do not count. Adding an unanswered word to a completed level
makes that level active again.

HTTP 401 indicates missing or invalid authentication. HTTP 500 indicates a
database failure and returns `{"message":"cannot load vocabulary levels"}`.
