package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveC(n, m int, k int64, mat [][]int) string {
	D := n + m - 1
	type cell struct {
		p int
		d int
	}
	cells := make([]cell, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cells = append(cells, cell{p: mat[i][j], d: i + j})
		}
	}
	sort.Slice(cells, func(i, j int) bool { return cells[i].p < cells[j].p })
	seen := make([]bool, D)
	order := make([]int, 0, D)
	for _, c := range cells {
		if !seen[c.d] {
			seen[c.d] = true
			order = append(order, c.d)
		}
	}
	fixed := make([]int, D)
	var countWays func() int64
	countWays = func() int64 {
		dp := make([]int64, D+2)
		dp[0] = 1
		for pos := 0; pos < D; pos++ {
			ndp := make([]int64, D+2)
			if fixed[pos] == 1 {
				for bal := 0; bal <= D; bal++ {
					v := dp[bal]
					if v == 0 {
						continue
					}
					nb := bal + 1
					ndp[nb] += v
					if ndp[nb] > k {
						ndp[nb] = k
					}
				}
			} else if fixed[pos] == 2 {
				for bal := 1; bal <= D; bal++ {
					v := dp[bal]
					if v == 0 {
						continue
					}
					nb := bal - 1
					ndp[nb] += v
					if ndp[nb] > k {
						ndp[nb] = k
					}
				}
			} else {
				for bal := 0; bal <= D; bal++ {
					v := dp[bal]
					if v == 0 {
						continue
					}
					nb := bal + 1
					ndp[nb] += v
					if ndp[nb] > k {
						ndp[nb] = k
					}
					if bal > 0 {
						nb2 := bal - 1
						ndp[nb2] += v
						if ndp[nb2] > k {
							ndp[nb2] = k
						}
					}
				}
			}
			dp = ndp
		}
		return dp[0]
	}
	for _, d := range order {
		fixed[d] = 1
		cnt := countWays()
		if cnt < k {
			k -= cnt
			fixed[d] = 2
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if fixed[i+j] == 1 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte(')')
			}
		}
		if i+1 < n {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int, int, int64, [][]int) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	k := int64(rng.Intn(5) + 1)
	mat := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			v := rng.Intn(10)
			row[j] = v
		}
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", row[j]))
		}
		sb.WriteByte('\n')
		mat[i] = row
	}
	return sb.String(), n, m, k, mat
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, m, k, mat := genCase(rng)
		expect := strings.TrimSpace(solveC(n, m, k, mat))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
