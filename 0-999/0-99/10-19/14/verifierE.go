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

func compute(n, t int) int64 {
	valleyTarget := t - 1
	dp := make([][][][][]int64, n+1)
	for i := range dp {
		dp[i] = make([][][][]int64, 5)
		for y := 0; y < 5; y++ {
			dp[i][y] = make([][][]int64, t+1)
			for p := 0; p <= t; p++ {
				dp[i][y][p] = make([][]int64, t+1)
				for v := 0; v <= t; v++ {
					dp[i][y][p][v] = make([]int64, 3)
				}
			}
		}
	}
	for y := 1; y <= 4; y++ {
		dp[1][y][0][0][0] = 1
	}
	for pos := 1; pos < n; pos++ {
		for lastY := 1; lastY <= 4; lastY++ {
			for p := 0; p <= t; p++ {
				for v := 0; v <= valleyTarget; v++ {
					for lastDir := 0; lastDir < 3; lastDir++ {
						cnt := dp[pos][lastY][p][v][lastDir]
						if cnt == 0 {
							continue
						}
						for y2 := 1; y2 <= 4; y2++ {
							if y2 == lastY {
								continue
							}
							var dir int
							if y2 > lastY {
								dir = 1
							} else {
								dir = 2
							}
							np, nv := p, v
							if lastDir == 1 && dir == 2 {
								np++
							} else if lastDir == 2 && dir == 1 {
								nv++
							}
							if np > t || nv > valleyTarget {
								continue
							}
							dp[pos+1][y2][np][nv][dir] += cnt
						}
					}
				}
			}
		}
	}
	var res int64
	for y := 1; y <= 4; y++ {
		for lastDir := 1; lastDir <= 2; lastDir++ {
			res += dp[n][y][t][valleyTarget][lastDir]
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(7) + 3
	t := rng.Intn(4) + 1
	if t > n {
		t = n
	}
	input := fmt.Sprintf("%d %d\n", n, t)
	return input, compute(n, t)
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
