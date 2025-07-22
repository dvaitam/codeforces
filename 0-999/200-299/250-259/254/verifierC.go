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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(s, t string) (int, string) {
	n := len(s)
	Cs := make([]int, 26)
	Ct := make([]int, 26)
	for i := 0; i < n; i++ {
		Cs[s[i]-'A']++
		Ct[t[i]-'A']++
	}
	best := 0
	for c := 0; c < 26; c++ {
		best += min(Cs[c], Ct[c])
	}
	remCt := make([]int, 26)
	remCs := make([]int, 26)
	copy(remCt, Ct)
	copy(remCs, Cs)
	matched := 0
	u := make([]byte, n)
	for i := 0; i < n; i++ {
		si := int(s[i] - 'A')
		for c := 0; c < 26; c++ {
			if remCt[c] == 0 {
				continue
			}
			inc := 0
			if c == si {
				inc = 1
			}
			sum := 0
			for x := 0; x < 26; x++ {
				rc := remCt[x]
				rcs := remCs[x]
				if x == c {
					rc--
				}
				if x == si {
					rcs--
				}
				if rc < 0 {
					rc = 0
				}
				if rcs < 0 {
					rcs = 0
				}
				sum += min(rc, rcs)
			}
			if matched+inc+sum >= best {
				u[i] = byte('A' + c)
				if c == si {
					matched++
				}
				remCt[c]--
				remCs[si]--
				break
			}
		}
	}
	z := n - best
	return z, string(u)
}

func generateCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(10) + 1
	b1 := make([]byte, n)
	b2 := make([]byte, n)
	for i := 0; i < n; i++ {
		b1[i] = byte('A' + rng.Intn(26))
		b2[i] = byte('A' + rng.Intn(26))
	}
	s := string(b1)
	t := string(b2)
	input := s + "\n" + t + "\n"
	return input, s, t
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, s, t := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var z int
		var u string
		if _, err := fmt.Fscan(strings.NewReader(out), &z, &u); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		expZ, expU := solve(s, t)
		if z != expZ || u != expU {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %s got %d %s\ninput:\n%s", i+1, expZ, expU, z, u, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
