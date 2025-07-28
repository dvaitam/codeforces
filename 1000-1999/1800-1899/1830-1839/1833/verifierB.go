package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Test struct {
	input    string
	expected string
}

func solveCase(n, k int, a, b []int) string {
	type pair struct{ val, idx int }
	arr := make([]pair, n)
	for i := 0; i < n; i++ {
		arr[i] = pair{a[i], i}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
	sort.Ints(b)
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[arr[i].idx] = b[i]
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	k := rng.Intn(100)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(201) - 100
		diff := rng.Intn(2*k+1) - k
		b[i] = a[i] + diff
	}
	// shuffle b
	rng.Shuffle(n, func(i, j int) { b[i], b[j] = b[j], b[i] })
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", a[i])
	}
	input += "\n"
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", b[i])
	}
	input += "\n"
	return Test{input: input, expected: solveCase(n, k, a, append([]int(nil), b...))}
}

func runCase(bin string, t Test) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(t.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(t.expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		tcase := generateCase(rng)
		if err := runCase(bin, tcase); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tcase.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
