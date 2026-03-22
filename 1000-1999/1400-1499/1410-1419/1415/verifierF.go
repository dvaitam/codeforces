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

func absI64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func maxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Embedded solver for 1415F
func solve1415F(input string) string {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)

	t := make([]int64, n+2)
	x := make([]int64, n+2)

	t[0] = 0
	x[0] = 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &t[i], &x[i])
	}
	const INF int64 = 2e18
	t[n+1] = INF
	x[n+1] = 0

	canSeq := make([]bool, n+2)
	for k := 1; k < n; k++ {
		canSeq[k] = (t[k]+absI64(x[k]-x[k+1]) <= t[k+1])
	}

	pos := make([]bool, n+2)
	dp := make([]int64, n+2)
	for i := 0; i <= n+1; i++ {
		dp[i] = INF
	}
	pos[0] = true

	ans := false

	for i := 0; i < n; i++ {
		if pos[i] {
			if t[i]+absI64(x[i]-x[i+1]) <= t[i+1] {
				if i+1 == n {
					ans = true
				} else {
					pos[i+1] = true
					val := t[i] + absI64(x[i]-x[i+1]) + absI64(x[i+1]-x[i+2])
					if val < dp[i+1] {
						dp[i+1] = val
					}
				}
			}

			isValid := true
			for j := i + 2; j <= n; j++ {
				if j-2 >= i+1 {
					if !canSeq[j-2] {
						isValid = false
					}
				}
				if !isValid {
					break
				}

				T := t[i] + absI64(x[i]-x[j]) + absI64(x[j]-x[i+1])
				if T <= t[i+1] {
					if j == n {
						ans = true
					} else {
						val := t[j-1] + absI64(x[j-1]-x[j+1])
						if val < dp[j] {
							dp[j] = val
						}
					}
				}
			}
		}

		if dp[i] != INF {
			if dp[i] <= t[i+1] {
				if i+1 == n {
					ans = true
				} else {
					pos[i+1] = true
					val := maxI64(t[i], dp[i]) + absI64(x[i+1]-x[i+2])
					if val < dp[i+1] {
						dp[i+1] = val
					}
				}
			}

			isValid := true
			for j := i + 2; j <= n; j++ {
				if j-2 >= i+1 {
					if !canSeq[j-2] {
						isValid = false
					}
				}
				if !isValid {
					break
				}

				T := maxI64(t[i], dp[i]+absI64(x[i+1]-x[j])) + absI64(x[j]-x[i+1])
				if T <= t[i+1] {
					if j == n {
						ans = true
					} else {
						val := t[j-1] + absI64(x[j-1]-x[j+1])
						if val < dp[j] {
							dp[j] = val
						}
					}
				}
			}
		}
	}

	if ans {
		return "YES"
	}
	return "NO"
}

func run(bin string, input string) (string, error) {
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

func uniqueRandInts(rng *rand.Rand, n int, low, high int64) []int64 {
	m := make(map[int64]struct{})
	res := make([]int64, 0, n)
	for len(res) < n {
		v := rng.Int63n(high-low+1) + low
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func runCase(bin string, input string) error {
	expect := solve1415F(input)

	got, err := run(bin, input)
	if err != nil {
		return err
	}
	// Compare case-insensitively for YES/NO
	if !strings.EqualFold(expect, got) {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// edge case
	{
		var sb strings.Builder
		sb.WriteString("1\n1 0\n")
		if err := runCase(bin, sb.String()); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	for total < 100 {
		n := rng.Intn(6) + 1
		tVals := make([]int64, n+1)
		xVals := make([]int64, n+1)
		curT := int64(0)
		for i := 1; i <= n; i++ {
			curT += int64(rng.Intn(5) + 1)
			tVals[i] = curT
		}
		coords := uniqueRandInts(rng, n, -10, 10)
		for i := 1; i <= n; i++ {
			xVals[i] = coords[i-1]
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tVals[i], xVals[i]))
		}
		if err := runCase(bin, sb.String()); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
