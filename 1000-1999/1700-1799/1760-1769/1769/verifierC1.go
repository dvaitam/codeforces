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

func solve(a []int) int {
	c := make([]int, 102)
	for _, v := range a {
		c[v]++
	}
	dp := [2]int{0, -1 << 60}
	best := 0
	for day := 1; day <= 101; day++ {
		nd := [2]int{-1 << 60, -1 << 60}
		for carry := 0; carry <= 1; carry++ {
			cur := dp[carry]
			if cur < 0 {
				continue
			}
			avail := c[day]
			if carry == 1 {
				newCarry := 0
				if avail > 0 {
					newCarry = 1
				}
				if cur+1 > nd[newCarry] {
					nd[newCarry] = cur + 1
				}
				if cur+1 > best {
					best = cur + 1
				}
				if avail > 0 {
					newCarry = 0
					if avail-1 > 0 {
						newCarry = 1
					}
					if cur+1 > nd[newCarry] {
						nd[newCarry] = cur + 1
					}
					if cur+1 > best {
						best = cur + 1
					}
				}
				newCarry = 0
				if avail > 0 {
					newCarry = 1
				}
				if 0 > nd[newCarry] {
					nd[newCarry] = 0
				}
			} else {
				if avail > 0 {
					newCarry := 0
					if avail-1 > 0 {
						newCarry = 1
					}
					if cur+1 > nd[newCarry] {
						nd[newCarry] = cur + 1
					}
					if cur+1 > best {
						best = cur + 1
					}
				}
				newCarry := 0
				if avail > 0 {
					newCarry = 1
				}
				if 0 > nd[newCarry] {
					nd[newCarry] = 0
				}
			}
		}
		dp = nd
	}
	return best
}

func runCase(bin string, a []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := solve(a)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(50) + 1
	a := make([]int, n)
	cur := rng.Intn(5) + 1
	for i := 0; i < n; i++ {
		cur += rng.Intn(3)
		a[i] = cur
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
