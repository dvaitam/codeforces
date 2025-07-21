package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TestCase struct {
	p1, p2, p3 int
	ans        int
}

func compute(p1, p2, p3 int) int {
	maxD := -1
	for d1 := 0; d1 <= 30; d1++ {
		for d2 := 0; d2 <= 30; d2++ {
			for d3 := 0; d3 <= 30; d3++ {
				r1 := p1 - d1 - d2
				r2 := p2 - d1 - d3
				r3 := p3 - d2 - d3
				if r1 < 0 || r2 < 0 || r3 < 0 {
					continue
				}
				if r1%2 != 0 || r2%2 != 0 || r3%2 != 0 {
					continue
				}
				d := d1 + d2 + d3
				if d > maxD {
					maxD = d
				}
			}
		}
	}
	return maxD
}

func genCases(n int) []TestCase {
	rand.Seed(time.Now().UnixNano())
	cases := make([]TestCase, n)
	for i := 0; i < n; i++ {
		p1 := rand.Intn(31)
		p2 := p1 + rand.Intn(31-p1)
		p3 := p2 + rand.Intn(31-p2)
		cases[i] = TestCase{p1, p2, p3, compute(p1, p2, p3)}
	}
	return cases
}

func buildInput(cs []TestCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintf(&sb, "%d %d %d\n", c.p1, c.p2, c.p3)
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases(100)
	input := buildInput(cases)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines, got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, res := range outputs {
		v, err := strconv.Atoi(res)
		if err != nil || v != cases[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cases[i].ans, res)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
