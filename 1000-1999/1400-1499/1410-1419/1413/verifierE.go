package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveE(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	nextInt64 := func() int64 {
		if !in.Scan() {
			return 0
		}
		v, _ := strconv.ParseInt(in.Text(), 10, 64)
		return v
	}
	t := nextInt64()
	var sb strings.Builder
	for ; t > 0; t-- {
		a := nextInt64()
		b := nextInt64()
		c := nextInt64()
		d := nextInt64()
		if a > b*c {
			sb.WriteString("-1")
			if t > 1 {
				sb.WriteByte('\n')
			}
			continue
		}
		w := c / d
		u := c - w*d
		compute := func(n, r int64) int64 {
			tval := n*d + r
			m := (tval - c) / d
			if tval < c {
				m = -1
			}
			if m > n {
				m = n
			}
			full := m + 1
			if full < 0 {
				full = 0
			}
			K := n - m
			heal := full*c + (K-1)*K/2*d + r*K
			heal *= b
			dmg := (n + 1) * a
			return dmg - heal
		}
		best := a
		n1 := a / (b * d)
		candidates := []int64{0, n1 - 1, n1, w}
		for _, n := range candidates {
			if n < 0 || n > w {
				continue
			}
			s := compute(n, 0)
			if s > best {
				best = s
			}
		}
		if u > 0 {
			s := compute(w, u)
			if s > best {
				best = s
			}
		}
		sb.WriteString(strconv.FormatInt(best, 10))
		if t > 1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(5))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		a := int64(rng.Intn(10) + 1)
		b := int64(rng.Intn(10) + 1)
		c := int64(rng.Intn(10) + 1)
		d := int64(rng.Intn(10) + 1)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, c, d))
		input := sb.String()
		expect := solveE(input)
		tests[i] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			return
		}
		if out != tc.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
