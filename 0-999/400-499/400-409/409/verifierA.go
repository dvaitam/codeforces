package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s1, s2 string) string {
	switch {
	case s1 > s2:
		return "TEAM 1 WINS"
	case s2 > s1:
		return "TEAM 2 WINS"
	default:
		return "TIE"
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	type test struct{ s1, s2 string }
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		l := 2 + i%19
		c1 := byte('a' + i%26)
		c2 := byte('a' + (i*7)%26)
		s1 := strings.Repeat(string(c1), l)
		s2 := strings.Repeat(string(c2), l)
		if i%3 == 0 {
			s2 = s1
		}
		tests = append(tests, test{s1, s2})
	}

	for i, t := range tests {
		inp := fmt.Sprintf("%s\n%s\n", t.s1, t.s2)
		want := expected(t.s1, t.s2)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
