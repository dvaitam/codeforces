package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type TestCase struct {
	input  string
	expect string
}

func isClose(a, b float64) bool {
	const eps = 1e-6 // Problem requires 10^-6
	absA, absB := math.Abs(a), math.Abs(b)
	diff := math.Abs(a - b)

	if diff <= eps {
		return true
	}
	// Relative error
	if diff <= eps*math.Max(absA, absB) {
		return true
	}
	return false
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin) // Removed .go check
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedA(n int, l, t int64, a []int64) string {
	t *= 2
	full := t / l
	left := t % l
	var res int64
	j := 0
	phase := false
	for i := 0; i < n; i++ {
		if left+a[i] >= l {
			phase = true
			j = 0
			left -= l
		}
		for j < n && a[j] <= left+a[i] {
			j++
		}
		if !phase {
			res += (full+1)*int64(j-i-1) + full*int64(n-j+i)
		} else {
			res += (full+1)*int64(n+j-i-1) + full*int64(i-j)
		}
	}
	ans := float64(res) / 4.0
	return fmt.Sprintf("%.6f", ans)
}

func genCase(rng *rand.Rand) TestCase {
	n := rng.Intn(5) + 2
	l := int64(rng.Intn(50) + int(n) + 5)
	t := int64(rng.Intn(50) + 1)
	vals := rng.Perm(int(l))[:n]
	sort.Ints(vals)
	a := make([]int64, n)
	for i, v := range vals {
		a[i] = int64(v)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, l, t)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	expect := expectedA(n, l, t, a)
	return TestCase{sb.String(), expect}
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		gotStr, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		// Parse expected and got values
		var expectVal, gotVal float64
		_, err = fmt.Sscan(tc.expect, &expectVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse expected output %q: %v\n", i+1, tc.expect, err)
			os.Exit(1)
		}
		_, err = fmt.Sscan(gotStr, &gotVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse candidate output %q: %v\n", i+1, gotStr, err)
			os.Exit(1)
		}

		if !isClose(expectVal, gotVal) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expect, gotStr, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
