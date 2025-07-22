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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func compute(a, b, w, x, c int64) int64 {
	if c <= a {
		return 0
	}
	D := c - a
	g := gcd(w, x)
	L := w / g
	wrap := make([]int64, L)
	var S int64
	for i := int64(0); i < L; i++ {
		bi := (b - (i*x)%w + w) % w
		if bi < x {
			wrap[i] = 1
			S++
		}
	}
	noWrap := L - S
	k := D / noWrap
	rem := D - k*noWrap
	if rem == 0 {
		rem = noWrap
		k--
	}
	prefix := make([]int64, L+1)
	for i := int64(0); i < L; i++ {
		prefix[i+1] = prefix[i] + wrap[i]
	}
	var r int64
	for i := int64(1); i <= L; i++ {
		if i-prefix[i] >= rem {
			r = i
			break
		}
	}
	return k*L + r
}

func genCase(r *rand.Rand) Test {
	w := int64(r.Intn(10) + 2)
	x := int64(r.Intn(int(w-1)) + 1)
	a := int64(r.Intn(20) + 1)
	b := int64(r.Intn(int(w)))
	c := a + int64(r.Intn(20))
	input := fmt.Sprintf("%d %d %d %d %d\n", a, b, w, x, c)
	out := fmt.Sprintf("%d", compute(a, b, w, x, c))
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 25; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
