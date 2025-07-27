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

func expectedB(aCost, bCost, p int64, s string) int {
	n := len(s)
	var cost int64
	var last byte
	ans := 1
	for i := n - 2; i >= 0; i-- {
		if s[i] != last {
			if s[i] == 'A' {
				cost += aCost
			} else {
				cost += bCost
			}
			last = s[i]
		}
		if cost > p {
			ans = i + 2
			break
		}
	}
	if cost <= p {
		ans = 1
	}
	return ans
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	letters := []byte{'A', 'B'}
	for i := 0; i < 100; i++ {
		n := rand.Intn(18) + 2
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[rand.Intn(2)])
		}
		s := sb.String()
		a := rand.Intn(10) + 1
		b := rand.Intn(10) + 1
		p := rand.Intn(20) + 1
		input := fmt.Sprintf("1\n%d %d %d\n%s\n", a, b, p, s)
		expect := expectedB(int64(a), int64(b), int64(p), s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
