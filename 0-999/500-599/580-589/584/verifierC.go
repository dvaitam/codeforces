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

func note(x, y byte) byte {
	for z := byte('a'); z <= 'z'; z++ {
		if z != x && z != y {
			return z
		}
	}
	return 'a'
}

func solve(s1, s2 string, t int) (string, bool) {
	n := len(s1)
	eq := 0
	for i := 0; i < n; i++ {
		if s1[i] == s2[i] {
			eq++
		}
	}
	dif := n - eq
	sdMin := t - dif
	if sdMin < 0 {
		sdMin = 0
	}
	half := (dif + 1) / 2
	sdMax := t - half
	if sdMax > eq {
		sdMax = eq
	}
	if sdMin > sdMax {
		return "", false
	}
	sd := sdMin
	tPrime := t - sd
	db := 2*tPrime - dif

	ans := make([]byte, n)
	c := 0
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			if c < db {
				ans[i] = note(s1[i], s2[i])
			} else {
				if c%2 == 1 {
					ans[i] = s1[i]
				} else {
					ans[i] = s2[i]
				}
			}
			c++
		}
	}
	c = 0
	for i := 0; i < n; i++ {
		if s1[i] == s2[i] {
			if c < sd {
				ans[i] = byte('a' + (s1[i]-'a'+1)%26)
			} else {
				ans[i] = s1[i]
			}
			c++
		}
	}
	return string(ans), true
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func diff(a, b string) int {
	cnt := 0
	for i := range a {
		if a[i] != b[i] {
			cnt++
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type tc struct {
		n      int
		t      int
		s1, s2 string
	}
	cases := make([]tc, 0, 100)
	// some edge cases
	cases = append(cases, tc{1, 0, "a", "a"})
	cases = append(cases, tc{1, 1, "a", "b"})
	cases = append(cases, tc{2, 1, "aa", "aa"})

	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		t := rng.Intn(n + 1)
		var b1, b2 strings.Builder
		for i := 0; i < n; i++ {
			b1.WriteByte(byte('a' + rng.Intn(26)))
			b2.WriteByte(byte('a' + rng.Intn(26)))
		}
		cases = append(cases, tc{n, t, b1.String(), b2.String()})
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n%s\n%s\n", tc.n, tc.t, tc.s1, tc.s2)
		expect, ok := solve(tc.s1, tc.s2, tc.t)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if !ok {
			if out != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:%s", idx+1, out, input)
				os.Exit(1)
			}
			continue
		}
		if out == "-1" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected string got -1\ninput:%s", idx+1, input)
			os.Exit(1)
		}
		if len(out) != tc.n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected length %d got %d\n", idx+1, tc.n, len(out))
			os.Exit(1)
		}
		if diff(tc.s1, out) != tc.t || diff(tc.s2, out) != tc.t {
			fmt.Fprintf(os.Stderr, "case %d failed: output doesn't differ by t\ninput:%soutput:%s\n", idx+1, input, out)
			os.Exit(1)
		}
		if expect != out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
