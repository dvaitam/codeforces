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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	instr := make([][]int, n)
	for i := 0; i < n; i++ {
		instr[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &instr[i][j])
		}
	}

	coreLocked := make([]bool, n)
	coreLockTime := make([]int, n)
	cellLocked := make([]bool, k+1)

	for cycle := 1; cycle <= m; cycle++ {
		writes := make(map[int][]int)
		for i := 0; i < n; i++ {
			if coreLocked[i] {
				continue
			}
			x := instr[i][cycle-1]
			if x == 0 {
				continue
			}
			if cellLocked[x] {
				coreLocked[i] = true
				coreLockTime[i] = cycle
			} else {
				writes[x] = append(writes[x], i)
			}
		}
		for cell, cores := range writes {
			if len(cores) >= 2 {
				cellLocked[cell] = true
				for _, ci := range cores {
					if !coreLocked[ci] {
						coreLocked[ci] = true
						coreLockTime[ci] = cycle
					}
				}
			}
		}
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d\n", coreLockTime[i])
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := rng.Intn(k + 1)
			fmt.Fprintf(&sb, "%d", x)
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin, input string) error {
	expect := solve(input)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"1 1 1\n0\n",
		"2 1 1\n1\n1\n",
		"2 2 2\n1 0\n0 1\n",
		"3 3 1\n1 1 1\n0 0 0\n1 1 1\n",
		"1 3 1\n1 1 1\n",
	}
	for i := 0; i < 100; i++ {
		tests = append(tests, generateCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
