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

const modD int64 = 998244353

type op struct {
	typ byte
	val int64
}

type testCase struct {
	input    string
	expected int64
}

func solve(ops []op) int64 {
	n := len(ops)
	ans := int64(0)
	for i := 0; i < n; i++ {
		if ops[i].typ != '+' {
			continue
		}
		x := ops[i].val
		dp := make([][2]int64, n+1)
		dp[0][0] = 1
		for j := 0; j < n; j++ {
			ndp := make([][2]int64, n+1)
			if j == i {
				for k := 0; k <= n; k++ {
					ndp[k][1] = (ndp[k][1] + dp[k][0] + dp[k][1]) % modD
				}
				dp = ndp
				continue
			}
			if ops[j].typ == '+' {
				val := ops[j].val
				if (j < i && val <= x) || (j > i && val < x) {
					for k := 0; k <= n; k++ {
						for b := 0; b < 2; b++ {
							v := dp[k][b]
							if v == 0 {
								continue
							}
							ndp[k][b] = (ndp[k][b] + v) % modD
							if k+1 <= n {
								ndp[k+1][b] = (ndp[k+1][b] + v) % modD
							}
						}
					}
				} else {
					for k := 0; k <= n; k++ {
						for b := 0; b < 2; b++ {
							v := dp[k][b]
							if v == 0 {
								continue
							}
							ndp[k][b] = (ndp[k][b] + 2*v) % modD
						}
					}
				}
			} else { // '-'
				for k := 0; k <= n; k++ {
					for b := 0; b < 2; b++ {
						v := dp[k][b]
						if v == 0 {
							continue
						}
						ndp[k][b] = (ndp[k][b] + v) % modD
						if k > 0 {
							ndp[k-1][b] = (ndp[k-1][b] + v) % modD
						} else {
							ndp[k][0] = (ndp[k][0] + v) % modD
						}
					}
				}
			}
			dp = ndp
		}
		sum := int64(0)
		for k := 0; k <= n; k++ {
			sum = (sum + dp[k][1]) % modD
		}
		ans = (ans + sum*x) % modD
	}
	return ans % modD
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	ops := make([]op, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			val := int64(rng.Intn(100) + 1)
			ops[i] = op{typ: '+', val: val}
			sb.WriteString(fmt.Sprintf("+ %d\n", val))
		} else {
			ops[i] = op{typ: '-'}
			sb.WriteString("-\n")
		}
	}
	return testCase{input: sb.String(), expected: solve(ops)}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
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
	if got%modD != tc.expected%modD {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
