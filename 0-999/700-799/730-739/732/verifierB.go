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

func expectedExtra(n, k int, a []int) int {
	if n == 0 {
		return 0
	}
	b := make([]int, n)
	b[0] = a[0]
	extra := 0
	for i := 1; i < n; i++ {
		need := k - b[i-1]
		if need < a[i] {
			need = a[i]
		}
		if need < 0 {
			need = 0
		}
		b[i] = need
		extra += b[i] - a[i]
	}
	return extra
}

func buildInput(n, k int, a []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string, n int) (int, []int, error) {
	fields := strings.Fields(out)
	if len(fields) < 1+n {
		return 0, nil, fmt.Errorf("not enough numbers")
	}
	var extra int
	if _, err := fmt.Sscan(fields[0], &extra); err != nil {
		return 0, nil, fmt.Errorf("bad extra: %v", err)
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Sscan(fields[1+i], &b[i]); err != nil {
			return 0, nil, fmt.Errorf("bad schedule: %v", err)
		}
	}
	return extra, b, nil
}

func runCase(bin string, n, k int, a []int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	input := buildInput(n, k, a)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	extra, b, err := parseOutput(out.String(), n)
	if err != nil {
		return err
	}
	// verify constraints
	if extra != expectedExtra(n, k, a) {
		return fmt.Errorf("expected %d got %d", expectedExtra(n, k, a), extra)
	}
	sum := 0
	for i := 0; i < n; i++ {
		if b[i] < a[i] {
			return fmt.Errorf("day %d less than input", i)
		}
		if i > 0 && b[i]+b[i-1] < k {
			return fmt.Errorf("constraint failed day %d", i)
		}
		sum += b[i] - a[i]
	}
	if sum != extra {
		return fmt.Errorf("extra mismatch")
	}
	return nil
}

func randomCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(10) + 1
	k := rng.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10)
	}
	return n, k, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	type tc struct {
		n, k int
		a    []int
	}
	cases := []tc{
		{3, 5, []int{1, 2, 3}},
		{1, 1, []int{0}},
		{2, 3, []int{1, 1}},
	}
	for len(cases) < 105 {
		n, k, a := randomCase(rng)
		cases = append(cases, tc{n, k, a})
	}
	for i, c := range cases {
		if err := runCase(bin, c.n, c.k, c.a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(c.n, c.k, c.a))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
