# データ設計

## ER図
```mermaid
erDiagram
    User ||--o{ Avatar : "has"
    User ||--o{ UserInfo : "has"
    User ||--o{ Mission : "owns"
    User ||--o{ UserAvatarRelation : "interacts_with"
    User ||--o{ MissionUnlock : "unlocked_by"
    User ||--o{ Matching : "matches_as_user1"
    User ||--o{ Matching : "matches_as_user2"
    Avatar ||--o{ UserAvatarRelation : "receives_interaction"
    UserInfo ||--|| Mission : "unlocked_by"
    Mission ||--o{ MissionUnlock : "has"

    User {
        string id PK "ULID"
        string firebase_uid UK "Firebase UID"
        string display_name "表示名"
        int gender "性別(0:男性,1:女性,2:その他)"
        timestamp birth_date "生年月日"
        string bio "自己紹介"
        string profile_image_url "プロフィール画像URL"
        boolean is_onboarding_completed "オンボーディング完了フラグ"
        timestamp created_at
        timestamp updated_at
    }

    Avatar {
        string id PK "ULID"
        string user_id FK "所有者のユーザID"
        string avatar_icon_url "アバターアイコンURL"
        string prompt "設定プロンプト"
        jsonb personality_traits "性格特性"
        timestamp created_at
        timestamp updated_at
    }

    UserAvatarRelation {
        string id PK "ULID"
        string user_id FK "ユーザID"
        string avatar_id FK "アバターID"
        int matching_point "マッチングポイント"
        timestamp created_at
        timestamp updated_at
    }

    Matching {
        string id PK "ULID"
        string user1_id FK "ユーザ1のID(user1_id < user2_id)"
        string user2_id FK "ユーザ2のID"
        timestamp created_at
        timestamp updated_at
    }

    UserInfo {
        string id PK "ULID"
        string user_id FK "ユーザID"
        string info_type "情報タイプ(text/image)"
        string key "キー(項目名/画像タイトル)"
        string value "値(テキスト内容/画像URL)"
        boolean is_mission_reward "ミッション報酬フラグ"
        timestamp created_at
        timestamp updated_at
    }

    Mission {
        string id PK "ULID"
        string mission_owner_user_id FK "ミッション所有者ID"
        string user_info_id FK "解禁対象のUserInfoのID"
        int threshold_point_condition "閾値ポイント条件"
        string unlock_condition "解禁条件"
        timestamp created_at
        timestamp updated_at
    }

    MissionUnlock {
        string id PK "ULID"
        string mission_id FK "ミッションID"
        string unlocked_user_id FK "解禁したユーザID"
        timestamp created_at
        timestamp updated_at
    }
```

## Firestore

### /notifications/${user_id}/notification/{notification_id}
通知情報を格納するコレクション

**パス変数:**
- `user_id`: 通知対象のユーザID (ULID)
- `notification_id`: 通知のID (ULID)

**ドキュメント構造:**
```json
{
  "id": "string (ULID)",
  "user_id": "string (ULID)",
  "title": "string (通知タイトル)",
  "message": "string (通知メッセージ)",
  "created_at": "timestamp"
}
```

### /onboarding_chats/${user_id}/chat/{chat_id}
オンボーディング時のAIキャラクタとの対話ログ

**パス変数:**
- `user_id`: ユーザID (ULID)
- `chat_id`: チャットメッセージのID (ULID)

**ドキュメント構造:**
```json
{
  "id": "string (ULID)",
  "sender_type": "string (user/system/avatar_ai)",
  "message": "string (メッセージ内容)",
  "created_at": "timestamp"
}
```

### /user_avatar_chats/${user_id}/${avatar_id}/{chat_id}
ユーザと異性ユーザのアバターAIとのチャットログ

**パス変数:**
- `user_id`: チャットを行うユーザID (ULID)
- `avatar_id`: 対話相手のアバターID (ULID)
- `chat_id`: チャットメッセージのID (ULID)

**ドキュメント構造:**
```json
{
  "id": "string (ULID)",
  "sender_type": "string (user/system/avatar_ai)",
  "message": "string (メッセージ内容)",
  "avatar": "Avatar (アバター情報のスナップショット)",
  "created_at": "timestamp"
}
```

### /user_chats/${user1_id}/${user2_id}/{chat_id}
マッチング成立後のユーザ同士のチャットログ

**パス変数:**
- `user1_id`: ユーザ1のID (ULID, user1_id < user2_id の制約あり)
- `user2_id`: ユーザ2のID (ULID)
- `chat_id`: チャットメッセージのID (ULID)

**制約:**
- `user1_id < user2_id` となるようにパスを構築すること（ULID辞書順）

**ドキュメント構造:**
```json
{
  "id": "string (ULID)",
  "sender_type": "string (user/system)",
  "message": "string (メッセージ内容)",
  "user1": "User (ユーザ1の情報スナップショット)",
  "user2": "User (ユーザ2の情報スナップショット)",
  "created_at": "timestamp"
}
```
