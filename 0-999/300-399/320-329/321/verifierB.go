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

type JiroCard struct {
	atk bool
	s   int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expectedB(jiro []JiroCard, ciel []int) int {
	n := len(jiro)
	m := len(ciel)
	memo := make(map[[2]int]int)
	var dfs func(jmask, cmask int) int
	dfs = func(jmask, cmask int) int {
		key := [2]int{jmask, cmask}
		if v, ok := memo[key]; ok {
			return v
		}
		best := 0
		for i := 0; i < m; i++ {
			if (cmask>>i)&1 != 0 {
				continue
			}
			if jmask == 0 {
				val := ciel[i] + dfs(jmask, cmask|(1<<i))
				if val > best {
					best = val
				}
				continue
			}
			for j := 0; j < n; j++ {
				if (jmask>>j)&1 == 0 {
					continue
				}
				c := ciel[i]
				jc := jiro[j]
				if jc.atk {
					if c >= jc.s {
						val := c - jc.s + dfs(jmask&^(1<<j), cmask|(1<<i))
						if val > best {
							best = val
						}
					}
				} else {
					if c > jc.s {
						val := dfs(jmask&^(1<<j), cmask|(1<<i))
						if val > best {
							best = val
						}
					}
				}
			}
		}
		memo[key] = best
		return best
	}
	fullMask := (1 << n) - 1
	return dfs(fullMask, 0)
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	jiro := make([]JiroCard, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		typ := rng.Intn(2)
		s := rng.Intn(11)
		if typ == 0 {
			sb.WriteString(fmt.Sprintf("ATK %d\n", s))
			jiro[i] = JiroCard{atk: true, s: s}
		} else {
			sb.WriteString(fmt.Sprintf("DEF %d\n", s))
			jiro[i] = JiroCard{atk: false, s: s}
		}
	}
	ciel := make([]int, m)
	for i := 0; i < m; i++ {
		c := rng.Intn(11)
		sb.WriteString(fmt.Sprintf("%d\n", c))
		ciel[i] = c
	}
	input := sb.String()
	expect := expectedB(jiro, ciel)
	return input, expect
}

func runCase(bin string, input string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resStr := strings.TrimSpace(out.String())
	var got int
	if _, err := fmt.Sscan(resStr, &got); err != nil {
		return fmt.Errorf("bad output %q", resStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
