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

func bigAdd(a, b string) string {
	i := len(a) - 1
	j := len(b) - 1
	carry := 0
	var res []byte
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0')
			j--
		}
		res = append(res, byte('0'+sum%10))
		carry = sum / 10
	}
	for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
		res[l], res[r] = res[r], res[l]
	}
	return string(res)
}

func lessNum(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func solve(l int, s string) string {
	mid := l / 2
	left := -1
	for i := mid; i >= 1; i-- {
		if s[i] != '0' {
			left = i
			break
		}
	}
	right := -1
	for i := mid + 1; i < l; i++ {
		if s[i] != '0' {
			right = i
			break
		}
	}
	var best string
	if left != -1 {
		best = bigAdd(s[:left], s[left:])
	}
	if right != -1 {
		cand := bigAdd(s[:right], s[right:])
		if best == "" || lessNum(cand, best) {
			best = cand
		}
	}
	return best
}

func runCase(bin string, l int, s string) error {
	input := fmt.Sprintf("%d\n%s\n", l, s)
	expect := solve(l, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
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
		l := rng.Intn(20) + 2
		b := make([]byte, l)
		b[0] = byte('1' + rng.Intn(9))
		for j := 1; j < l; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		s := string(b)
		if err := runCase(bin, l, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
