package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	in  string
	out string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func computeAnswer(s1, s2 string) string {
	parts := strings.Split(s1, "|")
	if len(parts) != 2 {
		return "Impossible"
	}
	left := parts[0]
	right := parts[1]
	a, b, c := len(left), len(right), len(s2)
	if (a+b+c)%2 != 0 || abs(a-b) > c {
		return "Impossible"
	}
	diff := abs(a - b)
	idx := 0
	if a < b {
		left += s2[idx : idx+diff]
	} else {
		right += s2[idx : idx+diff]
	}
	idx += diff
	rem := s2[idx:]
	half := len(rem) / 2
	left += rem[:half]
	right += rem[half:]
	return left + "|" + right
}

func genCase(r *rand.Rand) Test {
	letters := r.Perm(26)
	l1 := r.Intn(6)
	l2 := r.Intn(6)
	rem := r.Intn(6)
	if l1+l2+rem > 26 {
		rem = 26 - l1 - l2
	}
	idx := 0
	lb := make([]byte, l1)
	for i := 0; i < l1; i++ {
		lb[i] = byte('A' + letters[idx])
		idx++
	}
	rb := make([]byte, l2)
	for i := 0; i < l2; i++ {
		rb[i] = byte('A' + letters[idx])
		idx++
	}
	mb := make([]byte, rem)
	for i := 0; i < rem; i++ {
		mb[i] = byte('A' + letters[idx])
		idx++
	}
	s1 := string(lb) + "|" + string(rb)
	s2 := string(mb)
	input := s1 + "\n" + s2 + "\n"
	out := computeAnswer(s1, s2)
	return Test{input, out}
}

func runCase(bin string, t Test) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(t.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(t.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 25; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
