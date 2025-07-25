package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func simulate(seq string) int {
	stack := []int{}
	for _, c := range seq {
		switch c {
		case '+', '*':
			if len(stack) < 2 {
				return -1
			}
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res int
			if c == '+' {
				res = a + b
			} else {
				res = a * b
			}
			if res > 1_000_000 {
				return -1
			}
			stack = append(stack, res)
		default:
			stack = append(stack, int(c-'0'))
		}
	}
	if len(stack) == 0 {
		return 0
	}
	return stack[len(stack)-1]
}

func generate(r *rand.Rand) string {
	for {
		length := r.Intn(19) + 1
		stackDepth := 0
		var sb strings.Builder
		for i := 0; i < length; i++ {
			if stackDepth < 2 {
				d := r.Intn(10)
				sb.WriteByte(byte('0' + d))
				stackDepth++
				continue
			}
			if r.Intn(2) == 0 {
				d := r.Intn(10)
				sb.WriteByte(byte('0' + d))
				stackDepth++
			} else {
				op := '+'
				if r.Intn(2) == 0 {
					op = '*'
				}
				sb.WriteByte(byte(op))
				stackDepth--
			}
		}
		seq := sb.String()
		if res := simulate(seq); res != -1 {
			return seq
		}
	}
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(47))
	for i := 1; i <= 100; i++ {
		seq := generate(r)
		expected := fmt.Sprintf("%d", simulate(seq))
		input := fmt.Sprintf("%s\n", seq)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("wrong answer on test %d: expected %s got %s (seq %s)\n", i, expected, strings.TrimSpace(out), seq)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
