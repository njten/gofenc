# gofenc

Jednoduchý šifrovací nástroj napsaný v Go.

## O projektu

`gofenc` je CLI nástroj pro šifrování souborů a složek pomocí konceptu **vaultu** (trezoru). Vault je složka na disku, ve které jsou soubory šifrované individuálně. Na rozdíl od nástrojů jako Cryptomator nevyžaduje FUSE ani virtuální disk — soubory se šifrují a dešifrují explicitně příkazem.

## Funkce

- Šifrování jednotlivých souborů i celých složek
- Podpora dvou šifrovacích algoritmů: **AES-256-GCM** nebo **ChaCha20-Poly1305**
- Derivace klíče pomocí **Argon2id** (modernější alternativa k scrypt)
- Autentizace pomocí **hesla** nebo **keyfile**
- Volitelné šifrování názvů souborů
- Multiplatformní — Windows, macOS, Linux

## Použití

### Vytvoření vaultu

```bash
# Vault s heslem, AES-256-GCM, šifrování názvů souborů
gofenc vault init ./trezor --cipher aes-gcm --encrypt-names --auth password

# Vault s keyfile, ChaCha20
gofenc vault init ./trezor --cipher chacha20 --auth keyfile
```

### Práce se soubory

```bash
# Přidání souboru do vaultu
gofenc vault add ./trezor foto.jpg

# Odebrání souboru z vaultu
gofenc vault remove ./trezor foto.jpg

# Výpis obsahu vaultu
gofenc vault list ./trezor
```

### Zamčení / odemčení

```bash
# Zašifruje všechny soubory ve vaultu
gofenc vault lock ./trezor

# Dešifruje všechny soubory do výstupní složky
gofenc vault unlock ./trezor --out ./vystup
```

## Struktura projektu

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
│   ├── lock.go
│   └── unlock.go
├── vault/
│   ├── vault.go
│   ├── init.go
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── lock.go
│   └── unlock.go
└── crypto/
    ├── kdf.go
    ├── masterkey.go
    ├── encrypt.go
    ├── decrypt.go
    └── filename.go
```

## Struktura vaultu na disku

```
trezor/
├── vault.json        - metadata (algoritmus, KDF parametry, zašifrovaný master key)
└── files/
    ├── a1b2c3d4.enc  - zašifrovaný obsah souboru
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

### Formát .enc souboru

```
┌─────────────────────────────────────┐
│  HEADER (plaintext)                 │
│  - magic bytes: "GOFENC" (5 bytes)  │
│  - version: uint8 (1 byte)          │
│  - nonce: 12 bytes (AES-GCM)        │
│           nebo 24 bytes (ChaCha20)  │
├─────────────────────────────────────┤
│  ENCRYPTED PAYLOAD                  │
│  - původní název souboru            │
│  - obsah souboru (chunky po 64KB)   │
├─────────────────────────────────────┤
│  AUTH TAG (16 bytes)                │
│  - GCM nebo Poly1305 tag            │
└─────────────────────────────────────┘
```

## Kryptografický design

| Účel | Algoritmus |
|---|---|
| Šifrování obsahu | AES-256-GCM nebo ChaCha20-Poly1305 |
| Derivace klíče z hesla | Argon2id |
| Šifrování master keye | stejný algoritmus jako obsah |
| Integrita | AEAD auth tag (součást GCM / Poly1305) |

### Proč Argon2id místo scrypt?

Argon2id je vítěz [Password Hashing Competition](https://www.password-hashing.net) a je odolnější vůči side-channel útokům než scrypt, který používá například Cryptomator.

### Separace klíčů

Heslo nikdy přímo nešifruje data. Tok je následující:

```
heslo → Argon2id → wrapping key → dešifruje master key → šifruje soubory
```

Výhoda: změna hesla znamená pouze přešifrování master keye, ne celého vaultu.

## Závislosti

```
golang.org/x/crypto       - Argon2id, ChaCha20-Poly1305
github.com/spf13/cobra    - CLI framework
github.com/google/uuid    - UUID pro názvy .enc souborů
```

## Instalace

```bash
# Naklonování repozitáře
git clone https://github.com/njten/gofenc
cd gofenc

# Stažení závislostí
go mod download

# Sestavení binárky
go build -o gofenc .
```

## Threat model

| Hrozba | Chráněno? | Poznámka |
|---|---|---|
| Krádež disku / souboru | ✅ | AES-GCM / ChaCha20 |
| Slabé heslo | ⚠️ | Argon2id zpomaluje útok |
| Modifikace souboru | ✅ | AEAD auth tag |
| Metadata leak (názvy souborů) | ✅ / ❌ | volitelné šifrování názvů |
| Paměťový útok za běhu | ❌ | mimo scope |
| Kvantový počítač | ❌ | mimo scope |

## Licence
[MIT](LICENSE)