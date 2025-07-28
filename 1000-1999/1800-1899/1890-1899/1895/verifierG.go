package main

import (
	"bytes"
	"fmt"
	"math"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveG(n int, s string, r, b []int64) int64 {
	base := int64(0)
	for i := 0; i < n; i++ {
		base += b[i]
	}
	onesTotal := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			onesTotal++
		}
	}
	inf := int64(math.MinInt64 / 4)
	dp := make([]int64, onesTotal+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	onesSeen := 0
	for idx := 0; idx < n; idx++ {
		diff := r[idx] - b[idx]
		if s[idx] == '1' {
			for k := onesSeen; k >= 0; k-- {
				if dp[k] == inf {
					continue
				}
				cand := dp[k] + diff
				if cand > dp[k+1] {
					dp[k+1] = cand
				}
			}
			onesSeen++
		} else {
			for k := 0; k <= onesSeen; k++ {
				if dp[k] == inf {
					continue
				}
				cand := dp[k] + diff - int64(k)
				if cand > dp[k] {
					dp[k] = cand
				}
			}
		}
	}
	ans := int64(math.MinInt64)
	for k := 0; k <= onesSeen; k++ {
		if dp[k] > ans {
			ans = dp[k]
		}
	}
	return base + ans
}

func genCase(rng *rand.Rand) (int, string, []int64, []int64) {
	n := rng.Intn(5) + 1
	b := make([]int64, n)
	r := make([]int64, n)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		b[i] = int64(rng.Intn(10))
	}
	for i := 0; i < n; i++ {
		r[i] = int64(rng.Intn(10))
	}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return n, sb.String(), r, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, s, rArr, bArr := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		sb.WriteString(s)
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rArr[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(bArr[j]))
		}
		sb.WriteByte('\n')
		expect := fmt.Sprint(solveG(n, s, rArr, bArr))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
