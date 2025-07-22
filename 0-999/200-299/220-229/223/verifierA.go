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

func solve(s string) (int, string) {
	n := len(s)
	cnt := make([]int, n+1)
	for i := 1; i <= n; i++ {
		cnt[i] = cnt[i-1]
		if s[i-1] == '[' {
			cnt[i]++
		}
	}
	stack := []int{0}
	L, R, ans := 0, 0, 0
	for i := 1; i <= n; i++ {
		c := s[i-1]
		if c == '[' || c == '(' {
			stack = append(stack, i)
		} else {
			var top byte
			if len(stack) > 0 {
				pos := stack[len(stack)-1]
				if pos > 0 {
					top = s[pos-1]
				}
			}
			if (c == ']' && top != '[') || (c == ')' && top != '(') {
				stack = []int{i}
				continue
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				l := stack[len(stack)-1]
				r := i
				if ans < cnt[r]-cnt[l] {
					ans = cnt[r] - cnt[l]
					L = l
					R = r
				}
			}
		}
	}
	if L < R {
		return ans, s[L:R]
	}
	return 0, ""
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	chars := []byte{'(', ')', '[', ']'}
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b)
}

func runCase(bin, input string, ans int, seq string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 lines output, got %d", len(lines))
	}
	gotAns := strings.TrimSpace(lines[0])
	gotStr := strings.TrimSpace(lines[1])
	if gotAns != fmt.Sprint(ans) || gotStr != seq {
		return fmt.Errorf("expected %d %q got %s %q", ans, seq, gotAns, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		ans, seq := solve(s)
		if err := runCase(bin, s, ans, seq); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
