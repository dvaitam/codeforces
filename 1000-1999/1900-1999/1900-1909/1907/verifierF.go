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
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("1907F_official_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, "1907F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build official solution: %v %s", err, string(out))
	}
	return tmp, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	official, err := buildOfficial()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(official)

	rand.Seed(5)
	t := 100
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rand.Intn(7) + 1
		input.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			v := rand.Intn(10)
			input.WriteString(fmt.Sprintf("%d ", v))
		}
		input.WriteString("\n")
	}

	exp, err := runBinary(official, input.String())
	if err != nil {
		fmt.Printf("official solution failed: %v\n", err)
		os.Exit(1)
	}
	out, err := runBinary(bin, input.String())
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
