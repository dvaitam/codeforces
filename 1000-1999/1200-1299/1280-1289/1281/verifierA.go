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
	tmp, err := os.CreateTemp("", "refA-")
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
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genTest(rng *rand.Rand) string {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintln(&sb, t)
	suffixes := []string{"po", "desu", "masu", "mnida"}
	for i := 0; i < t; i++ {
		suf := suffixes[rng.Intn(len(suffixes))]
		l := rng.Intn(20)
		sb.WriteString(randWord(rng, l))
		sb.WriteString(suf)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	refBin, err := compileRef("1281A.go")
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genTest(rng)
		expectOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(expectOut)
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
