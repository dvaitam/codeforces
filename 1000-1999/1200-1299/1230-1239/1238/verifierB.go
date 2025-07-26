package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n int
	r int
	x []int
}

func expected(n, r int, x []int) string {
	sort.Slice(x, func(i, j int) bool { return x[i] > x[j] })
	uniq := []int{x[0]}
	for i := 1; i < n; i++ {
		if x[i] != x[i-1] {
			uniq = append(uniq, x[i])
		}
	}
	shots := 0
	shift := 0
	for _, pos := range uniq {
		if pos-shift <= 0 {
			break
		}
		shots++
		shift += r
	}
	return fmt.Sprintf("%d", shots)
}

func runCandidate(bin, input string) (string, error) {
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		r := rng.Intn(10) + 1
		x := make([]int, n)
		for j := 0; j < n; j++ {
			x[j] = rng.Intn(50) + 1
		}
		tests = append(tests, testCase{n: n, r: r, x: x})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d\n", t.n, t.r)
		for _, v := range t.x {
			input += fmt.Sprintf("%d ", v)
		}
		input = strings.TrimSpace(input) + "\n"
		want := expected(t.n, t.r, append([]int(nil), t.x...))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
