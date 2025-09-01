package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	q := rng.Intn(4) + 1
	v := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		v[i] = rng.Intn(21) - 10
		c[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		a := rng.Intn(21) - 10
		b := rng.Intn(21) - 10
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	return sb.String()
}

func expectedOutput(input string) string {
	fields := strings.Fields(input)
	p := 0
	n := atoi(fields[p])
	p++
	q := atoi(fields[p])
	p++
	v := make([]int64, n)
	for i := 0; i < n; i++ {
		v[i] = int64(atoi(fields[p]))
		p++
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = atoi(fields[p])
		p++
	}
	const NEG int64 = -1 << 60
	var sb strings.Builder
	for qi := 0; qi < q; qi++ {
		a := int64(atoi(fields[p]))
		p++
		b := int64(atoi(fields[p]))
		p++
		dp := make([]int64, n+1)
		for i := range dp {
			dp[i] = NEG
		}
		best1Val, best1Col := NEG, 0
		best2Val, best2Col := NEG, 0
		ans := int64(0)
		for i := 0; i < n; i++ {
			col := c[i]
			vi := v[i]
			cand := vi * b
			if dp[col] != NEG {
				t := dp[col] + vi*a
				if t > cand {
					cand = t
				}
			}
			other := best1Val
			if best1Col == col {
				other = best2Val
			}
			if other != NEG {
				t := other + vi*b
				if t > cand {
					cand = t
				}
			}
			if cand > dp[col] {
				dp[col] = cand
			}
			if dp[col] > ans {
				ans = dp[col]
			}
			val := dp[col]
			if best1Col == col {
				if val > best1Val {
					best1Val = val
				}
			} else {
				if val > best1Val {
					best2Val, best2Col = best1Val, best1Col
					best1Val, best1Col = val, col
				} else if best2Col == col {
					if val > best2Val {
						best2Val = val
					}
				} else if val > best2Val {
					best2Val, best2Col = val, col
				}
			}
		}
		if ans < 0 {
			ans = 0
		}
		sb.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(sb.String())
}

func atoi(s string) int {
	sign := 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i = 1
	}
	v := 0
	for ; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v * sign
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := generateCase(rng)
		expect := expectedOutput(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
