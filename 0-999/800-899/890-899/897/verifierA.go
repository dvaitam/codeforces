package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(reader *bufio.Reader) string {
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	var s string
	fmt.Fscan(reader, &s)
	b := []byte(s)
	for i := 0; i < m; i++ {
		var l, r int
		var c1, c2 string
		fmt.Fscan(reader, &l, &r, &c1, &c2)
		for j := l - 1; j < r; j++ {
			if string(b[j]) == c1 {
				b[j] = c2[0]
			}
		}
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	m := rng.Intn(100) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	s := sb.String()
	input := fmt.Sprintf("%d %d\n%s\n", n, m, s)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		c1 := string(byte('a' + rng.Intn(26)))
		c2 := string(byte('a' + rng.Intn(26)))
		input += fmt.Sprintf("%d %d %s %s\n", l, r, c1, c2)
	}
	expect := solve(bufio.NewReader(strings.NewReader(input)))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
