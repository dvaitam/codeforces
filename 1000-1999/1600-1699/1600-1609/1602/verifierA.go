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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(reader *bufio.Reader) string {
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		minIdx := 0
		for i := 1; i < len(s); i++ {
			if s[i] < s[minIdx] {
				minIdx = i
			}
		}
		a := string(s[minIdx])
		b := s[:minIdx] + s[minIdx+1:]
		fmt.Fprintf(&sb, "%s %s\n", a, b)
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for ; t > 0; t-- {
		n := rng.Intn(99) + 2 // length 2..100
		data := make([]byte, n)
		for i := 0; i < n; i++ {
			data[i] = byte(rng.Intn(26) + 'a')
		}
		fmt.Fprintf(&sb, "%s\n", string(data))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
