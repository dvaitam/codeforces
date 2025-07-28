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

// solveF is a placeholder that matches the simple stub in 1627F.go.
func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		for i := 0; i < n; i++ {
			var r1, c1, r2, c2 int
			fmt.Fscan(in, &r1, &c1, &r2, &c2)
			_ = r1
			_ = c1
			_ = r2
			_ = c2
		}
		out.WriteString("0\n")
	}
	return strings.TrimSpace(out.String())
}

func runProg(bin, input string) (string, error) {
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

func generateTests() []string {
	rng := rand.New(rand.NewSource(6))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		k := 2 * (rng.Intn(3) + 1)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			r1 := rng.Intn(k) + 1
			c1 := rng.Intn(k) + 1
			dir := rng.Intn(4)
			r2, c2 := r1, c1
			switch dir {
			case 0:
				if r1 > 1 {
					r2 = r1 - 1
				} else {
					r2 = r1 + 1
				}
			case 1:
				if r1 < k {
					r2 = r1 + 1
				} else {
					r2 = r1 - 1
				}
			case 2:
				if c1 > 1 {
					c2 = c1 - 1
				} else {
					c2 = c1 + 1
				}
			case 3:
				if c1 < k {
					c2 = c1 + 1
				} else {
					c2 = c1 - 1
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r1, c1, r2, c2))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveF(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
