package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const inf = int(1e9)

var valid [128][128][128]bool

func init() {
	for p := 0; p < 128; p++ {
		for c := 0; c < 128; c++ {
			for n := 0; n < 128; n++ {
				good := true
				for j := 1; j <= 5; j++ {
					if ((c>>j)&1) == 1 &&
						((p>>(j-1))&1) == 1 &&
						((p>>(j+1))&1) == 1 &&
						((n>>(j-1))&1) == 1 &&
						((n>>(j+1))&1) == 1 {
						good = false
						break
					}
				}
				valid[p][c][n] = good
			}
		}
	}
}

func expected(grid []string) string {
	var orig [7]int
	for i := 0; i < 7; i++ {
		mask := 0
		for j := 0; j < 7; j++ {
			if grid[i][j] == 'B' {
				mask |= 1 << j
			}
		}
		orig[i] = mask
	}
	var cost [7][128]int
	for i := 0; i < 7; i++ {
		for m := 0; m < 128; m++ {
			cost[i][m] = bits.OnesCount(uint(m ^ orig[i]))
		}
	}

	var dp [128][128]int
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			dp[i][j] = inf
		}
	}
	for m0 := 0; m0 < 128; m0++ {
		for m1 := 0; m1 < 128; m1++ {
			dp[m0][m1] = cost[0][m0] + cost[1][m1]
		}
	}
	for row := 1; row <= 5; row++ {
		var newdp [128][128]int
		for i := 0; i < 128; i++ {
			for j := 0; j < 128; j++ {
				newdp[i][j] = inf
			}
		}
		for prev := 0; prev < 128; prev++ {
			for cur := 0; cur < 128; cur++ {
				val := dp[prev][cur]
				if val >= inf {
					continue
				}
				for next := 0; next < 128; next++ {
					if !valid[prev][cur][next] {
						continue
					}
					v := val + cost[row+1][next]
					if v < newdp[cur][next] {
						newdp[cur][next] = v
					}
				}
			}
		}
		dp = newdp
	}
	res := inf
	for prev := 0; prev < 128; prev++ {
		for cur := 0; cur < 128; cur++ {
			if dp[prev][cur] < res {
				res = dp[prev][cur]
			}
		}
	}
	return fmt.Sprint(res)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func randGrid(rng *rand.Rand) []string {
	grid := make([]string, 7)
	for i := 0; i < 7; i++ {
		row := make([]byte, 7)
		for j := 0; j < 7; j++ {
			if rng.Intn(2) == 0 {
				row[j] = 'W'
			} else {
				row[j] = 'B'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	var tests [][]string
	for i := 0; i < 100; i++ {
		tests = append(tests, randGrid(rng))
	}
	for idx, grid := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		for _, row := range grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(grid)
		if got != exp {
			fmt.Printf("test %d failed: expected=%s got=%s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
