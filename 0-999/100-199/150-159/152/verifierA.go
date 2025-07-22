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

func expected(input string) (string, error) {
	in := strings.Fields(input)
	if len(in) < 2 {
		return "", fmt.Errorf("invalid input")
	}
	idx := 0
	n := toInt(in[idx])
	idx++
	m := toInt(in[idx])
	idx++
	grades := make([]string, n)
	for i := 0; i < n; i++ {
		if idx >= len(in) {
			return "", fmt.Errorf("invalid input grades")
		}
		grades[i] = in[idx]
		idx++
	}
	maxGrade := make([]byte, m)
	for j := 0; j < m; j++ {
		maxGrade[j] = '0'
		for i := 0; i < n; i++ {
			if grades[i][j] > maxGrade[j] {
				maxGrade[j] = grades[i][j]
			}
		}
	}
	count := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grades[i][j] == maxGrade[j] {
				count++
				break
			}
		}
	}
	return fmt.Sprintf("%d", count), nil
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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	m := rng.Intn(8) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sb.WriteByte(byte('1' + rng.Intn(9)))
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		exp, err := expected(strings.ReplaceAll(tc, "\n", " "))
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			os.Exit(1)
		}
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
