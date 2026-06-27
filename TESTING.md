# Testing Guide for cli-claude45-haiku-go

## Unit Tests

All unit tests pass with:
```bash
go test -v ./...
```

This runs:
- TestEncryptDecryptSymmetric - Tests basic encrypt/decrypt round-trip
- TestEncryptDecryptWithoutSalt - Tests auto-generated salt
- TestDecryptWithWrongPassword - Tests security against wrong password
- TestEncryptDecryptBinaryData - Tests binary file handling
- TestEncryptDecryptLargeData - Tests 10MB+ files
- TestEncryptFilenameBase32 - Tests filename encryption with base32
- TestEncryptFilenameBase64 - Tests filename encryption with base64
- TestPKCS7Padding - Tests padding algorithm

## CLI Tests

### Test 1: Basic Encryption and Decryption

```bash
# Create test file
echo "Test content with Bip39 style words" > test_file.txt

# Encrypt with password prompt
./cli-encrypt -i test_file.txt -action encrypt -o test_file.encrypted
# (Enter password when prompted)

# Decrypt
./cli-encrypt -i test_file.encrypted -action decrypt -o test_file.decrypted

# Verify content matches
diff test_file.txt test_file.decrypted
```

### Test 2: Environment Variable Support

```bash
# Encrypt with environment variable password
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.txt -action encrypt -o test_file.encrypted

# Decrypt with environment variable password
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.encrypted -action decrypt -o test_file.decrypted

# Verify
diff test_file.txt test_file.decrypted
```

### Test 3: Custom Salt

```bash
# Generate a base64 salt
SALT=$(head -c 16 /dev/urandom | base64)

# Encrypt with custom salt
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.txt -action encrypt -salt "$SALT" -o test_file.encrypted

# Decrypt with same salt
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.encrypted -action decrypt -o test_file.decrypted

# Verify
diff test_file.txt test_file.decrypted
```

### Test 4: Base64 Filename Encoding

```bash
# Encrypt with base64 encoding
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.txt -action encrypt -encoding base64 -o test_file.encrypted

# Decrypt
RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i test_file.encrypted -action decrypt -o test_file.decrypted

# Verify
diff test_file.txt test_file.decrypted
```

### Test 5: Large File Handling

```bash
# Create 100MB test file
dd if=/dev/urandom of=large_test.bin bs=1M count=100

# Encrypt
time RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i large_test.bin -action encrypt -o large_test.encrypted

# Decrypt
time RCLONE_ENCRYPT_PASSWORD="Testpassword1" ./cli-encrypt -i large_test.encrypted -action decrypt -o large_test.decrypted

# Verify
md5sum large_test.bin large_test.decrypted
```

### Test 6: Security Warning for --password Flag

```bash
# This should show a warning
./cli-encrypt -i test_file.txt -action encrypt -password "mypass" -o test_file.encrypted 2>&1 | grep WARNING
```

## Release and Installation Testing

### Windows (Scoop)

When release is available:
```bash
# Install from release (once Scoop bucket is set up)
scoop bucket add cli-claude45-haiku-go https://github.com/llm-supermarket-org/cli-claude45-haiku-go.git
scoop install cli-claude45-haiku-go/cli-claude45-haiku-go

# Test installation
cli-encrypt -version

# Test functionality
echo "test" > test.txt
cli-encrypt -i test.txt -action encrypt -password "test123" -o test.encrypted
cli-encrypt -i test.encrypted -action decrypt -password "test123" -o test.decrypted
diff test.txt test.decrypted
```

### macOS/Linux (Homebrew)

When release is available:
```bash
# Install (once Homebrew tap is set up)
brew tap llm-supermarket-org/cli-claude45-haiku-go https://github.com/llm-supermarket-org/cli-claude45-haiku-go.git
brew install cli-claude45-haiku-go

# Test installation
cli-encrypt -version

# Test functionality
echo "test" > test.txt
cli-encrypt -i test.txt -action encrypt -password "test123" -o test.encrypted
cli-encrypt -i test.encrypted -action decrypt -password "test123" -o test.decrypted
diff test.txt test.decrypted
```

## Cross-Platform Testing

The release build creates binaries for:
- Windows (amd64, arm64)
- macOS (amd64, arm64)
- Linux (amd64, arm64, armv7)

Each should be tested on the target platform to ensure proper functionality.

## Test Files

Pre-encrypted test files are included in the repository:
- `kr9tu4e1da4u3nifdd99g9tf5o` - Base32 encoded filename
- `Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY` - Base64 encoded filename

Note: These are encrypted with rclone format and may not be compatible with this CLI's implementation unless the encryption format is fully rclone-compatible.
