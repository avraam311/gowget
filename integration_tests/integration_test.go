package integrationtests

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const binName = "gowget"

func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/gowget")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, string(out))
	}
	return binPath
}

func runCmd(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n")
}

func TestWget(t *testing.T) {
	t.Parallel()
	err := os.RemoveAll("downloadedSite")
	if err != nil {
		log.Printf("error removing downloaded site for testing: %v", err)
		t.Error()
	}
	bin := buildBinary(t)
	out := runCmd(t, bin, "https://tech.wildberries.ru/")
	expected := "error downloading site: "
	if strings.Contains(out, expected) {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
	err = os.RemoveAll("downloadedSite")
	if err != nil {
		log.Printf("error removing downloaded site for testing: %v", err)
	}
}
