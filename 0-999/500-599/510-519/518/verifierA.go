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

func expected(s, t string) string {
	// compute next lexicographical string after s
	n := len(s)
	next := []byte(s)
	i := n - 1
	for ; i >= 0; i-- {
		if next[i] != 'z' {
			next[i]++
			break
		}
		next[i] = 'a'
	}
	if i < 0 {
		return "No such string"
	}
	res := string(next)
	if res < t {
		return res
	}
	return "No such string"
}

func generateCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(10) + 1
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = byte('a' + rng.Intn(26))
	}
	// generate t greater than s
	t := make([]byte, n)
	copy(t, s)
	pos := rng.Intn(n)
	// ensure some position increases
	if t[pos] < 'z' {
		t[pos] = byte(int(t[pos]) + 1 + rng.Intn(int('z'-t[pos])))
	} else {
		j := pos
		for j >= 0 && t[j] == 'z' {
			j--
		}
		if j >= 0 {
			t[j] = byte(int(t[j]) + 1 + rng.Intn(int('z'-t[j])))
		} else {
			// all z, append 'a'
			t = append(t, 'a')
			n++
		}
	}
	for i := pos + 1; i < n && i < len(t); i++ {
		t[i] = byte('a' + rng.Intn(26))
	}
	return string(s), string(t), expected(string(s), string(t))
}

func runCase(bin, s, t, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	input := fmt.Sprintf("%s\n%s\n", s, t)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ s, t, exp string }{
		{"a", "c", expected("a", "c")},
		{"aaaa", "zzzz", expected("aaaa", "zzzz")},
		{"az", "ba", expected("az", "ba")},
	}
	for i := 0; i < 100; i++ {
		s, t, e := generateCase(rng)
		cases = append(cases, struct{ s, t, exp string }{s, t, e})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.s, tc.t, tc.exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, tc.s, tc.t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
