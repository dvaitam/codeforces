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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildOfficial() (string, error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("1907B_official_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, "1907B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build official solution: %v %s", err, string(out))
	}
	return tmp, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	official, err := buildOfficial()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(official)

	rand.Seed(1)
	tests := make([]string, 100)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZbB"
	for i := 0; i < 100; i++ {
		l := rand.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			sb.WriteByte(letters[rand.Intn(len(letters))])
		}
		tests[i] = sb.String()
	}

	input := fmt.Sprintf("%d\n", len(tests))
	for _, tcase := range tests {
		input += tcase + "\n"
	}

	exp, err := runBinary(official, input)
	if err != nil {
		fmt.Printf("official solution failed: %v\n", err)
		os.Exit(1)
	}
	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Printf("binary failed: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(out) {
		fmt.Printf("mismatch\nexpected:\n%s\nactual:\n%s\n", exp, out)
		os.Exit(1)
	}
	fmt.Println("All 100 test cases passed.")
}
