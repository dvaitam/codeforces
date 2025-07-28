package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return ""
	}
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		seen := make(map[string]int)
		found := false
		for i := 0; i < n-1 && !found; i++ {
			pair := s[i : i+2]
			if first, ok := seen[pair]; ok {
				if i-first >= 2 {
					found = true
				}
			} else {
				seen[pair] = i
			}
		}
		if found {
			fmt.Fprintln(&out, "YES")
		} else {
			fmt.Fprintln(&out, "NO")
		}
	}
	return strings.TrimSpace(out.String())
}

func genTestB(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 2
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		fmt.Fprintf(&buf, "%d %s\n", n, sb.String())
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestB(rng)
		expect := solveB(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
