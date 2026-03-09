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

const MOD = 1000000007

// oracle: take top m XOR values among all pairs, sum mod MOD.
func oracle(a []int, m int) int64 {
	var xors []int
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			xors = append(xors, a[i]^a[j])
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(xors)))
	var ans int64
	for i := 0; i < m && i < len(xors); i++ {
		ans = (ans + int64(xors[i])) % MOD
	}
	return ans
}

func genCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(7) + 2 // n in [2,8]
	maxPairs := n * (n - 1) / 2
	m := rng.Intn(maxPairs + 1)
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(1 << 20)
	}
	return n, m, a
}

func buildInput(n, m int, a []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 200; i++ {
		n, m, a := genCase(rng)
		input := buildInput(n, m, a)
		expect := oracle(a, m)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		var gotVal int64
		if _, err := fmt.Sscan(got, &gotVal); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output %q\ninput:\n%s", i, got, input)
			os.Exit(1)
		}
		if gotVal != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%s", i, expect, gotVal, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 200 tests passed")
}
