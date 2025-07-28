package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		res := make([]int, n*m)
		idx := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				d1 := i + j
				d2 := i + (m - 1 - j)
				d3 := (n - 1 - i) + j
				d4 := (n - 1 - i) + (m - 1 - j)
				maxd := d1
				if d2 > maxd {
					maxd = d2
				}
				if d3 > maxd {
					maxd = d3
				}
				if d4 > maxd {
					maxd = d4
				}
				res[idx] = maxd
				idx++
			}
		}
		sort.Ints(res)
		for i, v := range res {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
		out.WriteByte('\n')
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
	rng := rand.New(rand.NewSource(2))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveB(t)
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
