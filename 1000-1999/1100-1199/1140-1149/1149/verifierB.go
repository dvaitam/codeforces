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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "oracle1149B")
	cmd := exec.Command("go", "build", "-o", path, "1149B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	q := rng.Intn(6) + 1
	letters := []byte("abc")
	var sb strings.Builder
	s := make([]byte, n)
	for i := range s {
		s[i] = letters[rng.Intn(len(letters))]
	}
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	sb.Write(s)
	sb.WriteByte('\n')
	lengths := [3]int{}
	for t := 0; t < q; t++ {
		if rng.Intn(2) == 0 || (lengths[0] == 0 && lengths[1] == 0 && lengths[2] == 0) {
			idx := rng.Intn(3)
			ch := letters[rng.Intn(len(letters))]
			fmt.Fprintf(&sb, "+ %d %c\n", idx+1, ch)
			lengths[idx]++
		} else {
			idx := rng.Intn(3)
			for lengths[idx] == 0 {
				idx = rng.Intn(3)
			}
			fmt.Fprintf(&sb, "- %d\n", idx+1)
			lengths[idx]--
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		input := genCase(rng)
		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", tc, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n\ngot:\n%s\n", tc, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
