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

func longestNonDecreasing(s []byte) int {
	n := len(s)
	prefix4 := make([]int, n+1)
	suffix7 := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix4[i+1] = prefix4[i]
		if s[i] == '4' {
			prefix4[i+1]++
		}
	}
	for i := n - 1; i >= 0; i-- {
		suffix7[i] = suffix7[i+1]
		if s[i] == '7' {
			suffix7[i]++
		}
	}
	best := 0
	for i := 0; i <= n; i++ {
		val := prefix4[i] + suffix7[i]
		if val > best {
			best = val
		}
	}
	return best
}

func solveCase(s string, queries []string) string {
	arr := []byte(s)
	var out strings.Builder
	for _, q := range queries {
		parts := strings.Fields(q)
		if parts[0] == "switch" {
			l := toInt(parts[1]) - 1
			r := toInt(parts[2]) - 1
			for i := l; i <= r; i++ {
				if arr[i] == '4' {
					arr[i] = '7'
				} else {
					arr[i] = '4'
				}
			}
		} else {
			out.WriteString(fmt.Sprintf("%d\n", longestNonDecreasing(arr)))
		}
	}
	return strings.TrimRight(out.String(), "\n")
}

func toInt(s string) int { var x int; fmt.Sscan(s, &x); return x }

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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	m := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	arr := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			arr[i] = '4'
		} else {
			arr[i] = '7'
		}
	}
	sb.WriteString(string(arr))
	sb.WriteByte('\n')
	queries := make([]string, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = fmt.Sprintf("switch %d %d", l, r)
		} else {
			queries[i] = "count"
		}
		sb.WriteString(queries[i])
		sb.WriteByte('\n')
	}
	exp := solveCase(string(arr), queries)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
