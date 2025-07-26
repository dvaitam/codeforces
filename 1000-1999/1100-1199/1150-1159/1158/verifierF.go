package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD = 998244353

func solveF(n, c int, a []int) []int {
	maxR := n/c + 2
	dp := make([][]int, maxR)
	for j := range dp {
		dp[j] = make([]int, c+2)
	}
	dp[0][1] = 1
	for _, x := range a {
		newdp := make([][]int, maxR)
		for j := range newdp {
			newdp[j] = make([]int, c+2)
		}
		for j := 0; j < maxR; j++ {
			for t := 1; t <= c; t++ {
				v := dp[j][t]
				if v == 0 {
					continue
				}
				newdp[j][t] = (newdp[j][t] + v) % MOD
				if x == t {
					if t < c {
						newdp[j][t+1] = (newdp[j][t+1] + v) % MOD
					} else if j+1 < maxR {
						newdp[j+1][1] = (newdp[j+1][1] + v) % MOD
					}
				} else {
					newdp[j][t] = (newdp[j][t] + v) % MOD
				}
			}
		}
		dp = newdp
	}
	res := make([]int, n+1)
	for j := 0; j < maxR; j++ {
		for t := 1; t <= c; t++ {
			v := dp[j][t]
			if v != 0 && j <= n {
				res[j] = (res[j] + v) % MOD
			}
		}
	}
	res[0] = (res[0] - 1 + MOD) % MOD
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 1
		c := rng.Intn(3) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(c) + 1
		}
		input := fmt.Sprintf("%d %d\n", n, c)
		for i, v := range a {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		exp := solveF(n, c, a)
		var expStr strings.Builder
		for i, v := range exp {
			if i > 0 {
				expStr.WriteByte(' ')
			}
			expStr.WriteString(strconv.Itoa(v))
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expStr.String() {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", t+1, expStr.String(), got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
