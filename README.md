# gofenc

A simple file encryption tool written in Go.

## About

`gofenc` is a CLI tool for encrypting files and folders using a **vault** concept. A vault is a folder on the disk where files are encrypted individually. Unlike tools such as Cryptomator, it does not require FUSE or a virtual disk — files are encrypted and decrypted explicitly via command. Filenames are always encrypted.

## Features

- Encryption of individual files and entire folders
- Support for two encryption algorithms: **AES-256-GCM** or **ChaCha20-Poly1305**
- Key derivation using **Argon2id** (a more modern alternative to scrypt)
- Authentication via **password** or **keyfile** (auto-generated)
- Filename encryption — original filenames are never stored in plaintext
- Lock/unlock mechanism — locking re-wraps the master key with a fresh nonce and hides the encrypted files directory, making the vault unusable without the correct secret
- Cross-platform — Windows, macOS, Linux

## Installation MacOS / Linux

```bash
# Clone the repository
git clone https://github.com/njten/gofenc
cd gofenc

# Download dependencies
go mod download

# Build the binary
go build -o gofenc .

# Run
./gofenc --help
```

## Installation Windows
```powershell
# Clone the repository
git clone https://github.com/njten/gofenc
cd gofenc

# Download dependencies
go mod download

# Build the binary
go build -o gofenc.exe .

# Run
.\gofenc.exe --help
```
> **Note:** On Windows, replace `./` with `.\` in all commands.

## Usage on MacOS

### Creating a vault

```bash
# Vault with password and AES-256-GCM
gofenc init ./myvault --cipher aes-gcm --auth password

# Vault with auto-generated keyfile and ChaCha20
gofenc init ./myvault --cipher chacha20 --auth keyfile
# Output: Keyfile generated: ./myvault.key — keep it safe!
```

### Working with files

```bash
# Add a file to the vault
gofenc add ./myvault photo.jpg

# Add an entire directory
gofenc add ./myvault ./documents

# Remove a file from the vault (by index or original filename)
gofenc remove ./myvault 1
gofenc remove ./myvault photo.jpg

# List vault contents
gofenc list ./myvault
```

### Extracting files

```bash
# Decrypt and extract a single file by index (see: gofenc list)
gofenc extract ./myvault 1 ./output

# Decrypt and extract all files
gofenc extract-all ./myvault ./output
```

### Locking and unlocking

```bash
# Lock the vault — re-wraps master key, hides files, disables all operations
gofenc lock ./myvault

# Unlock the vault — verifies secret, restores files, re-enables all operations
gofenc unlock ./myvault
```

### Deleting a vault

```bash
# Delete the vault permanently (asks for confirmation)
gofenc delete ./myvault

# Skip confirmation
gofenc delete ./myvault --force
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
│   ├── lock.go
│   ├── unlock.go
│   ├── delete.go
│   └── helpers.go
├── vault/
│   ├── vault.go
│   ├── init.go
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── lock.go
│   ├── unlock.go
│   ├── extract.go
│   └── delete.go
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
├── .locked           — present when vault is locked
├── files/            — encrypted files (present when unlocked)
│   ├── a1b2c3d4.enc
│   └── e5f6g7h8.enc
└── .files/           — hidden files directory (present when locked)
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
  "auth": "password",
  "encrypted_master_key": "base64encodedEncryptedKey...",
  "master_key_nonce": "base64encodedNonce...",
  "files": [
    {
      "original_name": "base64encryptedFilename...",
      "encrypted_name": "uuid.enc"
    }
  ]
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
| Filename encryption | AES-256-GCM |
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

### Lock/unlock mechanism

Locking combines two security measures:

1. **Cryptographic lock** — the master key is re-wrapped with a fresh nonce and saved back to `vault.json`. The plaintext master key is wiped from memory. Simply deleting `.locked` does not restore access.
2. **Physical lock** — the `files/` directory is renamed to `.files/`. Even if `.locked` is removed, the encrypted files are no longer at the expected path and all vault operations will fail until a proper unlock is performed.

Unlocking verifies the user's secret, decrypts the master key, renames `.files/` back to `files/` and removes `.locked`.

## Dependencies

```
golang.org/x/crypto       — Argon2id, ChaCha20-Poly1305
github.com/spf13/cobra    — CLI framework
github.com/google/uuid    — UUID-based names for .enc files
golang.org/x/term         — secure password input (no echo)
```

## Threat model

| Threat | Protected? | Notes |
|---|---|---|
| Disk / file theft | ✅ | AES-GCM / ChaCha20 |
| Weak password | ⚠️ | Argon2id slows down brute-force |
| File tampering | ✅ | AEAD auth tag |
| Metadata leak (filenames) | ✅ | filenames are always encrypted |
| Bypassing lock by deleting .locked | ✅ | files/ directory is also hidden |
| In-memory attack at runtime | ❌ | out of scope |
| Quantum computer | ❌ | out of scope |

## License

[MIT](LICENSE)