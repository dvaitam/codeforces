package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000007

func expected(input string) string {
	fields := strings.Fields(input)
	idx := 0
	n := toInt(fields[idx])
	idx++
	m := toInt(fields[idx])
	idx++
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fields[idx]
		idx++
	}
	res := 1
	for col := 0; col < m; col++ {
		seen := make(map[byte]bool)
		for i := 0; i < n; i++ {
			seen[names[i][col]] = true
		}
		res = res * len(seen) % mod
	}
	return fmt.Sprintf("%d", res)
}

func toInt(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func randomName(m int, rng *rand.Rand) string {
	b := make([]byte, m)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	m := rng.Intn(6) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%s\n", randomName(m, rng))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		exp := expected(strings.ReplaceAll(tc, "\n", " "))
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
