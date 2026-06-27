package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

var version = "0.1.0"

func main() {
	var (
		inputFile  = flag.String("i", "", "Input file path")
		inputAlt   = flag.String("input-file", "", "Input file path (long form)")
		outputFile = flag.String("o", "", "Output file path")
		outputAlt  = flag.String("output-file", "", "Output file path (long form)")
		password   = flag.String("password", "", "Password (not recommended - use environment variable instead)")
		salt       = flag.String("salt", "", "Salt (optional, base64 encoded)")
		encoding   = flag.String("encoding", "base32", "Filename encoding: base32 or base64")
		action     = flag.String("action", "decrypt", "Action: encrypt or decrypt")
		showVer    = flag.Bool("version", false, "Show version")
		showHelp   = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *showVer {
		fmt.Printf("rclone-encrypt version %s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		fmt.Printf("Usage: rclone-encrypt [options]\n\n")
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  # Encrypt a file (prompts for password)\n")
		fmt.Printf("  rclone-encrypt -i myfile.txt -action encrypt\n\n")
		fmt.Printf("  # Decrypt a file\n")
		fmt.Printf("  rclone-encrypt -i myfile.encrypted -action decrypt -o myfile.txt\n\n")
		fmt.Printf("  # Use password from environment variable (more secure)\n")
		fmt.Printf("  RCLONE_ENCRYPT_PASSWORD=mypass rclone-encrypt -i file.encrypted\n")
		os.Exit(0)
	}

	if err := run(*inputFile, *inputAlt, *outputFile, *outputAlt, *password, *salt, *encoding, *action); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(inputFile, inputAlt, outputFile, outputAlt, password, saltStr, encoding, action string) error {
	// Handle long-form alternatives
	if inputFile == "" && inputAlt != "" {
		inputFile = inputAlt
	}
	if outputFile == "" && outputAlt != "" {
		outputFile = outputAlt
	}

	if inputFile == "" {
		return fmt.Errorf("input file required (-i or --input-file)")
	}

	// Validate action
	if action != "encrypt" && action != "decrypt" {
		return fmt.Errorf("action must be 'encrypt' or 'decrypt', got %q", action)
	}

	// Read input file
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Get password - try environment variable first, then command line, then prompt
	if password == "" {
		password = os.Getenv("RCLONE_ENCRYPT_PASSWORD")
	}

	if password == "" {
		pwd, err := promptPassword("Enter password: ")
		if err != nil {
			return err
		}
		if pwd == "" {
			return fmt.Errorf("password cannot be empty")
		}
		if len(pwd) < 12 {
			fmt.Fprintf(os.Stderr, "WARNING: Password is weak (less than 12 characters). Consider using a stronger password.\n")
		}
		password = pwd
	} else if os.Getenv("RCLONE_ENCRYPT_PASSWORD") == "" {
		// Only warn if using --password flag, not if using env var
		fmt.Fprintf(os.Stderr, "WARNING: Using --password flag is insecure!\n")
		fmt.Fprintf(os.Stderr, "  - Your password may appear in shell history\n")
		fmt.Fprintf(os.Stderr, "  - Use environment variables instead: RCLONE_ENCRYPT_PASSWORD=mypass\n")
		fmt.Fprintf(os.Stderr, "  - Or provide password interactively (just press Enter when prompted)\n\n")
	} else if password == "" && os.Getenv("RCLONE_ENCRYPT_PASSWORD") == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Get salt if needed
	var saltBytes []byte
	if saltStr != "" {
		fmt.Fprintf(os.Stderr, "Note: Using custom salt\n")
		saltBytes, err = base64.StdEncoding.DecodeString(saltStr)
		if err != nil {
			return fmt.Errorf("invalid salt encoding: %w", err)
		}
	}

	var output []byte
	if action == "encrypt" {
		output, err = encryptData(inputData, password, saltBytes)
	} else {
		output, err = decryptData(inputData, password)
	}

	if err != nil {
		return fmt.Errorf("failed to %s: %w", action, err)
	}

	// Write output
	if outputFile == "" {
		outputFile = inputFile + ".out"
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully %sed file: %s\n", action, outputFile)
	return nil
}

func promptPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Fprintln(os.Stderr)
	return string(bytePassword), nil
}

