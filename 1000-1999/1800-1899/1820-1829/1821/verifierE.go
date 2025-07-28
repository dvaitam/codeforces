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

func isRegular(s string) bool {
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func minCostSeq(s string, memo map[string]int) int {
	if s == "" {
		return 0
	}
	if v, ok := memo[s]; ok {
		return v
	}
	n := len(s)
	best := math.MaxInt32
	for i := 0; i < n-1; i++ {
		if s[i] == '(' && s[i+1] == ')' {
			t := s[:i] + s[i+2:]
			c := n - (i + 1) + minCostSeq(t, memo)
			if c < best {
				best = c
			}
		}
	}
	memo[s] = best
	return best
}

func reachable(s string, k int) map[string]struct{} {
	type state struct {
		str   string
		moves int
	}
	q := []state{{s, 0}}
	vis := map[string]int{s: 0}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.moves == k {
			continue
		}
		n := len(cur.str)
		for i := 0; i < n; i++ {
			for j := 0; j <= n; j++ {
				if j == i || j == i+1 {
					continue
				}
				t := cur.str[:i] + cur.str[i+1:]
				jj := j
				if j > i {
					jj = j - 1
				}
				t = t[:jj] + string(cur.str[i]) + t[jj:]
				if v, ok := vis[t]; !ok || v > cur.moves+1 {
					vis[t] = cur.moves + 1
					q = append(q, state{t, cur.moves + 1})
				}
			}
		}
	}
	res := make(map[string]struct{})
	for str := range vis {
		if isRegular(str) {
			res[str] = struct{}{}
		}
	}
	return res
}

func solveE(k int, s string) int {
	best := math.MaxInt32
	for str := range reachable(s, k) {
		c := minCostSeq(str, map[string]int{})
		if c < best {
			best = c
		}
	}
	return best
}

func genRegular(rng *rand.Rand, n int) string {
	open := n / 2
	close := open
	bal := 0
	var sb strings.Builder
	for open+close > 0 {
		if open > 0 && (bal == 0 || rng.Intn(open+close) < open) {
			sb.WriteByte('(')
			open--
			bal++
		} else {
			sb.WriteByte(')')
			close--
			bal--
		}
	}
	return sb.String()
}

func genCaseE(rng *rand.Rand) (int, string) {
	n := rng.Intn(3)*2 + 2
	s := genRegular(rng, n)
	k := rng.Intn(3)
	return k, s
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		k, s := genCaseE(rng)
		input := fmt.Sprintf("1\n%d\n%s\n", k, s)
		expect := fmt.Sprintf("%d", solveE(k, s))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
