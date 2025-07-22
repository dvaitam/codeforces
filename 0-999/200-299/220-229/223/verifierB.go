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

func solve(s, t string) string {
	n := len(s)
	m := len(t)
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1]
		if pref[i] < m && s[i-1] == t[pref[i]] {
			pref[i]++
		}
	}
	suff := make([]int, n+2)
	for i := n; i >= 1; i-- {
		suff[i] = suff[i+1]
		if suff[i] < m && s[i-1] == t[m-1-suff[i]] {
			suff[i]++
		}
	}
	for i := 1; i <= n; i++ {
		A := pref[i-1]
		B := suff[i+1]
		covered := false
		for j := 0; j < m; j++ {
			if t[j] == s[i-1] && j <= A && m-1-j <= B {
				covered = true
				break
			}
		}
		if !covered {
			return "No"
		}
	}
	return "Yes"
}

func genCase(rng *rand.Rand) (string, string) {
	letters := []byte{'a', 'b', 'c'}
	n := rng.Intn(8) + 1
	m := rng.Intn(8) + 1
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	tb := make([]byte, m)
	for i := range tb {
		tb[i] = letters[rng.Intn(len(letters))]
	}
	return string(sb), string(tb)
}

func runCase(bin, s, t, expected string) error {
	input := fmt.Sprintf("%s\n%s\n", s, t)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, t := genCase(rng)
		exp := solve(s, t)
		if err := runCase(bin, s, t, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, s, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
