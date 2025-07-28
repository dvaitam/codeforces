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

type Op struct {
	b int
	x int64
}

const INF int64 = 1e18

func solveCaseD(ops []Op, queries []int64) []int64 {
	n := len(ops)
	lens := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if ops[i-1].b == 1 {
			if lens[i-1] < INF {
				lens[i] = lens[i-1] + 1
				if lens[i] > INF {
					lens[i] = INF
				}
			} else {
				lens[i] = INF
			}
		} else {
			if lens[i-1] == 0 {
				lens[i] = 0
			} else if lens[i-1] >= INF/(ops[i-1].x+1) {
				lens[i] = INF
			} else {
				lens[i] = lens[i-1] * (ops[i-1].x + 1)
				if lens[i] > INF {
					lens[i] = INF
				}
			}
		}
	}
	res := make([]int64, len(queries))
	for qi, k := range queries {
		idx := n
		for {
			l, r := 1, idx
			pos := idx
			for l <= r {
				m := (l + r) / 2
				if lens[m] >= k {
					pos = m
					r = m - 1
				} else {
					l = m + 1
				}
			}
			op := ops[pos-1]
			if op.b == 1 {
				res[qi] = op.x
				break
			}
			k = (k-1)%lens[pos-1] + 1
			idx = pos - 1
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int64, []int64) {
	n := rng.Intn(4) + 1
	q := rng.Intn(4) + 1
	ops := make([]Op, n)
	ops[0] = Op{1, int64(rng.Intn(5) + 1)}
	for i := 1; i < n; i++ {
		b := rng.Intn(2) + 1
		if b == 1 {
			ops[i] = Op{1, int64(rng.Intn(5) + 1)}
		} else {
			ops[i] = Op{2, int64(rng.Intn(3) + 1)}
		}
	}
	lens := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if ops[i-1].b == 1 {
			lens[i] = lens[i-1] + 1
		} else {
			lens[i] = lens[i-1] * (ops[i-1].x + 1)
		}
		if lens[i] > INF {
			lens[i] = INF
		}
	}
	queries := make([]int64, q)
	for i := 0; i < q; i++ {
		if lens[n] == 0 {
			queries[i] = 1
		} else {
			queries[i] = int64(rng.Intn(int(lens[n])) + 1)
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for _, op := range ops {
		fmt.Fprintf(&sb, "%d %d\n", op.b, op.x)
	}
	for i := 0; i < q; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", queries[i])
	}
	sb.WriteByte('\n')
	expect := solveCaseD(ops, queries)
	return sb.String(), expect, queries
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect, _ := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", i+1, len(expect), len(fields), in)
			os.Exit(1)
		}
		for j, f := range fields {
			var v int64
			if _, err := fmt.Sscan(f, &v); err != nil || v != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d expected %d got %s\ninput:\n%s", i+1, j+1, expect[j], f, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
