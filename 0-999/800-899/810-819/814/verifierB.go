package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func validate(n int, a, b []int, output string) error {
	fields := strings.Fields(output)
	if len(fields) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	p := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		p[i] = v
	}

	freq := make([]int, n+1)
	for _, v := range p {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		freq[v]++
	}
	for v := 1; v <= n; v++ {
		if freq[v] != 1 {
			return fmt.Errorf("value %d appears %d times (not a permutation)", v, freq[v])
		}
	}

	da, db := 0, 0
	for i := 0; i < n; i++ {
		if p[i] != a[i] {
			da++
		}
		if p[i] != b[i] {
			db++
		}
	}
	if da != 1 {
		return fmt.Errorf("differs from a in %d positions (need exactly 1)", da)
	}
	if db != 1 {
		return fmt.Errorf("differs from b in %d positions (need exactly 1)", db)
	}
	return nil
}

// genCase generates a valid (a, b) pair derived from a random permutation p,
// by replacing exactly one position in each to a different value.
func genCase(rng *rand.Rand, n int) (a, b []int) {
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	a = make([]int, n)
	b = make([]int, n)
	copy(a, p)
	copy(b, p)

	// corrupt a at position i with a value != p[i]
	i := rng.Intn(n)
	ai := rng.Intn(n-1) + 1
	if ai >= p[i] {
		ai++
	}
	a[i] = ai

	// corrupt b at position j with a value != p[j]
	j := rng.Intn(n)
	bj := rng.Intn(n-1) + 1
	if bj >= p[j] {
		bj++
	}
	b[j] = bj

	// If a == b (only possible when i==j and ai==bj), move b's corruption
	// to position k != j so that b still differs from p in exactly one place.
	equal := true
	for k := 0; k < n; k++ {
		if a[k] != b[k] {
			equal = false
			break
		}
	}
	if equal {
		b[j] = p[j] // restore
		k := (j + 1) % n
		bk := rng.Intn(n-1) + 1
		if bk >= p[k] {
			bk++
		}
		b[k] = bk
	}
	return a, b
}

func buildInput(n int, a, b []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	run := func(idx int, n int, a, b []int) {
		input := buildInput(n, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var outBuf, errBuf strings.Builder
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(outBuf.String())
		if err := validate(n, a, b, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sgot: %s\nreason: %v\n", idx, input, got, err)
			os.Exit(1)
		}
	}

	idx := 0
	for idx < 200 {
		idx++
		n := rng.Intn(50) + 2
		a, b := genCase(rng, n)
		run(idx, n, a, b)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
