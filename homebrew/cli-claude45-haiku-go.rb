class CliClaude45HaikuGo < Formula
  desc "A CLI tool that encrypts and decrypts files using rclone encryption defaults"
  homepage "https://github.com/llm-supermarket-org/cli-claude45-haiku-go"
  url "https://github.com/llm-supermarket-org/cli-claude45-haiku-go/archive/v0.1.0.tar.gz"
  sha256 "abc123placeholder"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "cli-encrypt"
    bin.install "cli-encrypt"
  end

  test do
    # Create a test file
    test_file = testpath/"test.txt"
    test_file.write("Hello, World!")

    # Encrypt the file with a password
    system "#{bin}/cli-encrypt", "-i", test_file, "-action", "encrypt", "-password", "test123"
    assert_path_exists testpath/"test.txt.out"

    # Decrypt the file
    system "#{bin}/cli-encrypt", "-i", testpath/"test.txt.out", "-action", "decrypt", "-password", "test123", "-o", testpath/"test.decrypted.txt"
    assert_path_exists testpath/"test.decrypted.txt"
    assert_equal test_file.read, (testpath/"test.decrypted.txt").read
  end
end
