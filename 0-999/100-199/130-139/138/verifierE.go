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

type Constr struct {
	c    byte
	l, r int
}

func expectedCount(s string, k, L, R int, cons []Constr) int64 {
	n := len(s)
	pref := make([][26]int, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i]
		pref[i+1][s[i]-'a']++
	}
	var total int64
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			cnt := 0
			for _, c := range cons {
				val := pref[j+1][c.c-'a'] - pref[i][c.c-'a']
				if val >= c.l && val <= c.r {
					cnt++
				}
			}
			if cnt >= L && cnt <= R {
				total++
			}
		}
	}
	return total
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(6) + 1
	s := randomString(rng, n)
	k := rng.Intn(3) + 1
	L := rng.Intn(k + 1)
	R := L + rng.Intn(k-L+1)
	cons := make([]Constr, k)
	var sb strings.Builder
	fmt.Fprintln(&sb, s)
	fmt.Fprintf(&sb, "%d %d %d\n", k, L, R)
	for i := 0; i < k; i++ {
		c := byte('a' + rng.Intn(3))
		l := rng.Intn(n + 1)
		r := l + rng.Intn(n-l+1)
		cons[i] = Constr{c, l, r}
		fmt.Fprintf(&sb, "%c %d %d\n", c, l, r)
	}
	exp := expectedCount(s, k, L, R, cons)
	return sb.String(), exp
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output %q\n", t, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", t, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
