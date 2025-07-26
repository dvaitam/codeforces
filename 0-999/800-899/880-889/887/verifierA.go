package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(s string) string {
	idx := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			idx = i
			break
		}
	}
	if idx == -1 {
		return "no"
	}
	zero := 0
	for i := idx + 1; i < len(s); i++ {
		if s[i] == '0' {
			zero++
		}
	}
	if zero >= 6 {
		return "yes"
	}
	return "no"
}

func runCase(bin, s string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(s + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(strings.TrimSpace(s))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	cases := []string{
		"0",
		"1",
		"1000000",
		"1111111",
		"000000",
		"100000",
		"1100000",
		"1010101010",
	}
	for i := len(cases); i < 100; i++ {
		n := rand.Intn(100) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		cases = append(cases, sb.String())
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
