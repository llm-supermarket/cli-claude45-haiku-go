# cli-claude45-haiku-go

A small CLI tool that encrypts and decrypts files using the rclone encryption defaults.

Rclone uses a custom salt if no salt is provided, which this tool will use by default. A few similar tools:

- https://github.com/rclone/rclone
- https://github.com/mcolatosti/rclonedecrypt
- https://github.com/br0kenpixel/rclone-rcc
- @fyears/rclone-crypt

## Encryption Details

Rclone encryption uses:
- NaCl SecretBox (XSalsa20 + Poly1305) for the file contents
- AES256 for the filenames
- scrypt for key material derivation

## Installation

### Using Scoop (Windows)

```bash
scoop bucket add cli-claude45-haiku-go https://github.com/llm-supermarket-org/cli-claude45-haiku-go.git
scoop install cli-claude45-haiku-go/cli-claude45-haiku-go
```

### Using Homebrew (macOS / Linux)

```bash
brew tap llm-supermarket-org/cli-claude45-haiku-go https://github.com/llm-supermarket-org/cli-claude45-haiku-go.git
brew install cli-claude45-haiku-go
```

### From GitHub Releases

Download the latest release for your platform from [GitHub Releases](https://github.com/llm-supermarket-org/cli-claude45-haiku-go/releases).

### From Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/llm-supermarket-org/cli-claude45-haiku-go.git
cd cli-claude45-haiku-go
go build -o cli-encrypt
./cli-encrypt -help
```

## Usage

### Basic Encryption

Encrypt a file (you will be prompted for a password):

```bash
cli-encrypt -i myfile.txt -action encrypt -o myfile.encrypted
```

Or using long-form flags:

```bash
cli-encrypt --input-file myfile.txt --action encrypt --output-file myfile.encrypted
```

### Basic Decryption

Decrypt a file:

```bash
cli-encrypt -i myfile.encrypted -action decrypt -o myfile.txt
```

The action defaults to `decrypt`, so you can omit it:

```bash
cli-encrypt -i myfile.encrypted -o myfile.txt
```

### Using Environment Variables (More Secure)

Instead of using the `--password` flag (which appears in shell history), use an environment variable:

```bash
export RCLONE_ENCRYPT_PASSWORD="your-secure-password"
cli-encrypt -i myfile.txt -action encrypt -o myfile.encrypted
```

Or as a one-liner:

```bash
RCLONE_ENCRYPT_PASSWORD="your-secure-password" cli-encrypt -i myfile.txt -action encrypt
```

### Using Custom Salt

Provide a base64-encoded salt:

```bash
cli-encrypt -i myfile.txt -action encrypt -salt "dGVzdHNhbHQxMjM0NTY3OA==" -o myfile.encrypted
```

### Filename Encoding Options

Specify how to encode filenames (base32 or base64):

```bash
# Base32 (default)
cli-encrypt -i myfile.txt -action encrypt -encoding base32 -o myfile.encrypted

# Base64
cli-encrypt -i myfile.txt -action encrypt -encoding base64 -o myfile.encrypted
```

### Show Version

```bash
cli-encrypt -version
```

### Show Help

```bash
cli-encrypt -help
```

## Flags

- `-i, --input-file string` - Input file path (required)
- `-o, --output-file string` - Output file path (optional, defaults to input file with `.out` extension)
- `-password string` - Password (not recommended - use environment variable instead for security)
- `-salt string` - Salt (optional, base64 encoded)
- `-encoding string` - Filename encoding: base32 or base64 (default: base32)
- `-action string` - Action: encrypt or decrypt (default: decrypt)
- `-version` - Show version
- `-help` - Show help

## Examples

### Script: Encrypt All Files in Directory

```bash
#!/bin/bash
PASSWORD="your-secure-password"
TARGET_DIR="./documents"

for file in "$TARGET_DIR"/*; do
    if [ -f "$file" ]; then
        RCLONE_ENCRYPT_PASSWORD="$PASSWORD" cli-encrypt -i "$file" -action encrypt -o "$file.encrypted"
        echo "Encrypted: $file"
    fi
done
```

### Script: Decrypt All Files

```bash
#!/bin/bash
PASSWORD="your-secure-password"
TARGET_DIR="./documents"

for file in "$TARGET_DIR"/*.encrypted; do
    if [ -f "$file" ]; then
        OUTPUT="${file%.encrypted}"
        RCLONE_ENCRYPT_PASSWORD="$PASSWORD" cli-encrypt -i "$file" -action decrypt -o "$OUTPUT"
        echo "Decrypted: $file -> $OUTPUT"
    fi
done
```

### Using with Pipes (Partial Support)

While the tool primarily works with files, you can combine it with other utilities:

```bash
# Encrypt and pipe to a file
cli-encrypt -i plain.txt -action encrypt | tee encrypted.bin > /dev/null

# Note: Direct stdin/stdout is not yet supported; use file I/O instead
```

## Security Considerations

- **Never use the `--password` flag in production** - It appears in shell history and process listings
- Always use the `RCLONE_ENCRYPT_PASSWORD` environment variable for automation
- For interactive use, let the tool prompt you for the password (it won't appear in history)
- Ensure your passwords are strong and unique
- Keep your salts secure if using custom salts

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Compatibility

This tool is designed to be compatible with rclone's encryption format for encrypted file storage. It uses the same cryptographic algorithms and key derivation as rclone:

- File content: NaCl SecretBox (XSalsa20-Poly1305)
- Filenames: AES-256 CBC
- Key derivation: scrypt with parameters (2^17, 8, 1, 32)

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Issues

For bugs, feature requests, or questions, please open an issue on [GitHub Issues](https://github.com/llm-supermarket-org/cli-claude45-haiku-go/issues).
