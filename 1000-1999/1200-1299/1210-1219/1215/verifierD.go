package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(n int, s string) string {
	half := n / 2
	diff := 0
	q1, q2 := 0, 0
	for i := 0; i < half; i++ {
		if s[i] == '?' {
			q1++
		} else {
			diff += int(s[i] - '0')
		}
	}
	for i := half; i < n; i++ {
		if s[i] == '?' {
			q2++
		} else {
			diff -= int(s[i] - '0')
		}
	}
	if diff+(q1-q2)/2*9 == 0 {
		return "Bicarp"
	}
	return "Monocarp"
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(45)
	const T = 100
	digits := "0123456789?"
	for i := 0; i < T; i++ {
		half := rand.Intn(5) + 1
		n := half * 2
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = digits[rand.Intn(len(digits))]
		}
		// ensure even number of '?'
		cnt := 0
		for _, c := range b {
			if c == '?' {
				cnt++
			}
		}
		if cnt%2 == 1 {
			b[0] = '0'
		}
		s := string(b)
		input := fmt.Sprintf("%d\n%s\n", n, s)
		expect := expected(n, s)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput:%s\n", i+1, err, input)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expect {
			fmt.Printf("test %d failed: expected %s got %s\ninput:%s\n", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
