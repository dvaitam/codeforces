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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func simulate(start, T, a, b int, ops []int) int {
	c := start
	d := T - c
	for _, v := range ops {
		if v > 0 {
			move := v
			if move > c {
				move = c
			}
			free := b - d
			if move > free {
				move = free
			}
			c -= move
			d += move
		} else {
			move := -v
			if move > d {
				move = d
			}
			space := a - c
			if move > space {
				move = space
			}
			c += move
			d -= move
		}
	}
	return c
}

type testCaseE struct {
	n   int
	a   int
	b   int
	ops []int
}

func generateCaseE(rng *rand.Rand) (string, testCaseE) {
	n := rng.Intn(4) + 1
	a := rng.Intn(5) + 1
	b := rng.Intn(5) + 1
	ops := make([]int, n)
	for i := 0; i < n; i++ {
		v := rng.Intn(5) + 1
		if rng.Intn(2) == 0 {
			v = -v
		}
		ops[i] = v
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
	for i, v := range ops {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseE{n, a, b, ops}
}

func solveCaseE(n, a, b int, ops []int) [][]int {
	adds := make([]int, n)
	prefix := 0
	minPref := 0
	maxPref := 0
	for i := 0; i < n; i++ {
		adds[i] = -ops[i]
		prefix += adds[i]
		if prefix < minPref {
			minPref = prefix
		}
		if prefix > maxPref {
			maxPref = prefix
		}
	}
	total := prefix
	res := make([][]int, a+1)
	for i := 0; i <= a; i++ {
		res[i] = make([]int, b+1)
	}
	for T := 0; T <= a+b; T++ {
		L := max(0, T-b)
		U := min(a, T)
		if L > U {
			continue
		}
		resL := simulate(L, T, a, b, ops)
		resU := simulate(U, T, a, b, ops)
		thrLow := L - minPref
		thrHigh := U - maxPref
		for c := L; c <= U; c++ {
			var val int
			if c <= thrLow {
				val = resL
			} else if c >= thrHigh {
				val = resU
			} else {
				val = c + total
				if val < L {
					val = L
				} else if val > U {
					val = U
				}
			}
			d := T - c
			if d >= 0 && d <= b {
				res[c][d] = val
			}
		}
	}
	return res
}

func expectedE(tc testCaseE) string {
	res := solveCaseE(tc.n, tc.a, tc.b, tc.ops)
	var sb strings.Builder
	for i := 0; i <= tc.a; i++ {
		for j := 0; j <= tc.b; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", res[i][j]))
		}
		if i < tc.a {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseE(rng)
		expect := expectedE(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
