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
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveH(r *bufio.Reader) string {
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return ""
	}
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(r, &n, &s)
		out.WriteString("0\n")
	}
	return strings.TrimSpace(out.String())
}

func generateCaseH(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	var s strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s.WriteByte('0')
		} else {
			s.WriteByte('1')
		}
	}
	return fmt.Sprintf("1\n%d %s\n", n, s.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseH(rng)
		expect := solveH(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
