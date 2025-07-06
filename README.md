# SecureNest
**SecureNest** は、指定したフォルダ内のすべてのファイルを安全に暗号化・復号するCLIツールです。  
パスワードベースのArgon2id + ChaCha20-Poly1305で高いセキュリティを実現します。

## ✨ 特徴
- フォルダ内を再帰的に処理
- 暗号化時は`.vault`拡張子を付与
- 復号時は`.vault`拡張子を除去
- 処理後は元ファイルを自動削除
- Go製の軽量単一バイナリ

## 📦 フォルダ構成
```
SecureNest/
├── cmd/SecureNest/main.go
├── internal/vault/argon2.go
├── internal/vault/encrypt.go
├── internal/vault/decrypt.go
├── internal/vault/crypto.go
├── internal/vault/util.go
├── go.mod
└── README.md
```

## 🛠️ ビルド
Go 1.20以上が必要です。

```bash
go build -o SecureNest ./cmd/SecureNest
```

## 🚀 使い方
### 暗号化
指定したフォルダを再帰的に暗号化し、.vault拡張子を付与します。
```bash
./SecureNest -mode encrypt -dir /path/to/your/folder -time 3 -memory 65536 -threads 4
```
パスワードはプロンプトで入力します。

### 復号
.vault拡張子のファイルを復号し、元のファイル名で復元します。
```bash
./SecureNest -mode decrypt -dir /path/to/your/folder
```

### オプション
| オプション      | 説明                                   |
| ---------- | ------------------------------------ |
| `-mode`    | `encrypt` または `decrypt` を指定          |
| `-dir`     | 対象フォルダ（再帰処理されます）                     |
| `-time`    | Argon2のTimeパラメータ（暗号化時の計算量調整）         |
| `-memory`  | Argon2のMemoryパラメータ（KB単位、暗号化時のメモリ量調整） |
| `-threads` | Argon2の並列度（暗号化時のスレッド数）               |


## 🔒 暗号化方式
- ChaCha20-Poly1305 (AEAD)
- パスワードハッシュは Argon2id を使用（Salt:16byte, Key:32byte）

## ⚠️ 注意
- パスワードを忘れると復号できません。
- 復号や暗号化に失敗するとファイルが失われる可能性があります。事前にバックアップを推奨します。

