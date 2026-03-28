# gofenc

A simple file encryption tool written in Go.

## About

`gofenc` is a CLI tool for encrypting files and folders using a **vault** concept. A vault is a folder on disk where files are encrypted individually. Unlike tools such as Cryptomator, it does not require FUSE or a virtual disk — files are encrypted and decrypted explicitly via command.

## Features

- Encryption of individual files and entire folders
- Support for two encryption algorithms: **AES-256-GCM** or **ChaCha20-Poly1305**
- Key derivation using **Argon2id** (a more modern alternative to scrypt)
- Authentication via **password** or **keyfile**
- Optional filename encryption
- Cross-platform — Windows, macOS, Linux

## Usage

### Creating a vault

```bash
# Vault with password, AES-256-GCM, filename encryption enabled
gofenc init ./myvault --cipher aes-gcm --encrypt-names --auth password

# Vault with keyfile, ChaCha20
gofenc init ./myvault --cipher chacha20 --auth keyfile
```

### Working with files

```bash
# Add a file to the vault
gofenc add ./myvault photo.jpg

# Remove a file from the vault
gofenc remove ./myvault photo.jpg

# List vault contents
gofenc list ./myvault
```

### Extracting files

```bash
# Decrypt and extract a single file
gofenc extract ./myvault photo.jpg ./output

# Decrypt and extract all files
gofenc extract-all ./myvault ./output
```

## Project structure

```
gofenc/
├── main.go
├── go.mod
├── go.sum
├── cmd/
│   ├── root.go
│   ├── init.go
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── extract.go
│   └── helpers.go
├── vault/
│   ├── vault.go
│   ├── init.go
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── lock.go
│   ├── extract.go
│   └── unlock.go
└── crypto/
    ├── kdf.go
    ├── masterkey.go
    ├── encrypt.go
    ├── decrypt.go
    └── filename.go
```

## Vault structure on disk

```
myvault/
├── vault.json        — metadata (algorithm, KDF parameters, encrypted master key)
└── files/
    ├── a1b2c3d4.enc  — encrypted file content
    └── e5f6g7h8.enc
```

### vault.json

```json
{
  "version": 1,
  "cipher": "aes-256-gcm",
  "kdf": "argon2id",
  "kdf_params": {
    "time": 3,
    "memory": 65536,
    "threads": 4,
    "salt": "base64encodedSalt..."
  },
  "encrypt_filenames": true,
  "auth": "password",
  "encrypted_master_key": "base64encodedEncryptedKey...",
  "master_key_nonce": "base64encodedNonce..."
}
```

### .enc file format

```
┌─────────────────────────────────────┐
│  HEADER (plaintext)                 │
│  - magic bytes: "GOFENC" (6 bytes)  │
│  - version: uint8 (1 byte)          │
├─────────────────────────────────────┤
│  per chunk (64 KB):                 │
│  - nonce: 12 bytes (AES-GCM)        │
│           or 24 bytes (ChaCha20)    │
│  - chunk length: uint32 (4 bytes)   │
│  - encrypted data + auth tag        │
└─────────────────────────────────────┘
```

## Cryptographic design

| Purpose | Algorithm |
|---|---|
| Content encryption | AES-256-GCM or ChaCha20-Poly1305 |
| Key derivation from password | Argon2id |
| Master key encryption | same algorithm as content |
| Integrity | AEAD auth tag (GCM / Poly1305) |

### Why Argon2id instead of scrypt?

Argon2id won the [Password Hashing Competition](https://www.password-hashing.net) in 2015 and is more resistant to side-channel attacks than scrypt, which is used by Cryptomator. The `id` variant combines the GPU resistance of Argon2d with the side-channel resistance of Argon2i, making it the recommended choice per RFC 9106.

### Key separation

The password never directly encrypts data. The flow is:

```
password → Argon2id → wrapping key → decrypts master key → encrypts files
```

Benefit: changing the password only requires re-encrypting the master key, not the entire vault.

## Dependencies

```
golang.org/x/crypto       — Argon2id, ChaCha20-Poly1305
github.com/spf13/cobra    — CLI framework
github.com/google/uuid    — UUID-based names for .enc files
golang.org/x/term         — secure password input (no echo)
```

## Installation

```bash
# Clone the repository
git clone https://github.com/njten/gofenc
cd gofenc

# Download dependencies
go mod download

# Build the binary
go build -o gofenc .
```

## Threat model

| Threat | Protected? | Notes |
|---|---|---|
| Disk / file theft | ✅ | AES-GCM / ChaCha20 |
| Weak password | ⚠️ | Argon2id slows down brute-force |
| File tampering | ✅ | AEAD auth tag |
| Metadata leak (filenames) | ✅ / ❌ | optional filename encryption |
| In-memory attack at runtime | ❌ | out of scope |
| Quantum computer | ❌ | out of scope |

## License

[MIT](LICENSE)