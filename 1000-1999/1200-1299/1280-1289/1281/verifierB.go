package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef(src string) (string, error) {
	tmp, err := os.CreateTemp("", "refB-")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func randWord(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func genTest(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintln(&sb, t)
	for i := 0; i < t; i++ {
		sLen := rng.Intn(6) + 2
		cLen := rng.Intn(6) + 1
		fmt.Fprintf(&sb, "%s %s\n", randWord(rng, sLen), randWord(rng, cLen))
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	refBin, err := compileRef("1281B.go")
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genTest(rng)
		expOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(expOut)
		gotOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotOut)
		if got != expected {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
