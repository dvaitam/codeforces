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

func solve(n int, s, t string) int {
	type1 := []int{}
	type2 := []int{}
	for i := 0; i < n; i++ {
		if s[i] == 'a' && t[i] == 'b' {
			type1 = append(type1, i+1)
		} else if s[i] == 'b' && t[i] == 'a' {
			type2 = append(type2, i+1)
		}
	}
	s1 := len(type1)
	s2 := len(type2)
	if (s1%2)^(s2%2) == 1 {
		return -1
	}
	ops := s1/2 + s2/2
	if s1%2 == 1 {
		ops += 2
	}
	return ops
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(44)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(8) + 2
		a := make([]byte, n)
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				a[j] = 'a'
			} else {
				a[j] = 'b'
			}
			if rand.Intn(2) == 0 {
				b[j] = 'a'
			} else {
				b[j] = 'b'
			}
		}
		s := string(a)
		t := string(b)
		input := fmt.Sprintf("%d\n%s\n%s\n", n, s, t)

		expect := solve(n, s, t)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput: %s\n", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if expect == -1 {
			if len(fields) != 1 || fields[0] != "-1" {
				fmt.Printf("test %d failed: expected -1 got %s\ninput:%s\n", i+1, out, input)
				os.Exit(1)
			}
			continue
		}
		if len(fields) < 1 {
			fmt.Printf("test %d: empty output\n", i+1)
			os.Exit(1)
		}
		var k int
		if _, err := fmt.Sscanf(fields[0], "%d", &k); err != nil {
			fmt.Printf("test %d: parse k failed in '%s'\n", i+1, out)
			os.Exit(1)
		}
		if k != expect {
			fmt.Printf("test %d failed: expected k=%d got %d\ninput:%s\n", i+1, expect, k, input)
			os.Exit(1)
		}
		if len(fields) != 1+2*k {
			fmt.Printf("test %d failed: wrong number of indices\n", i+1)
			os.Exit(1)
		}
		sa := []byte(s)
		tb := []byte(t)
		idx := 1
		for j := 0; j < k; j++ {
			var x, y int
			fmt.Sscanf(fields[idx], "%d", &x)
			fmt.Sscanf(fields[idx+1], "%d", &y)
			idx += 2
			if x < 1 || x > n || y < 1 || y > n {
				fmt.Printf("test %d failed: index out of range\n", i+1)
				os.Exit(1)
			}
			sa[x-1], tb[y-1] = tb[y-1], sa[x-1]
		}
		if string(sa) != string(tb) {
			fmt.Printf("test %d failed: strings not equal after operations\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
