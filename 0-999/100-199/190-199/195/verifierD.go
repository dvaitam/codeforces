package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveDFromInput(input string) (int, error) {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return 0, err
	}
	m := make(map[[2]int64]struct{})
	for i := 0; i < n; i++ {
		var k, b int64
		fmt.Fscan(r, &k, &b)
		if k == 0 {
			continue
		}
		num := -b
		den := k
		if den < 0 {
			num = -num
			den = -den
		}
		g := gcd(abs64(num), den)
		num /= g
		den /= g
		m[[2]int64{num, den}] = struct{}{}
	}
	return len(m), nil
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		k := rng.Intn(11) - 5
		b := rng.Intn(21) - 10
		fmt.Fprintf(&sb, "%d %d\n", k, b)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expectInt, err := solveDFromInput(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			os.Exit(1)
		}
		expect := fmt.Sprint(expectInt)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
