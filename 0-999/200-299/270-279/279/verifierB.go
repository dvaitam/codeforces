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

func solveB(n int, t int64, a []int64) string {
	var sum int64
	ans := 0
	l := 0
	for r := 0; r < n; r++ {
		sum += a[r]
		for sum > t {
			sum -= a[l]
			l++
		}
		if r-l+1 > ans {
			ans = r - l + 1
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	t := int64(rng.Intn(500) + 1)
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(100) + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, t)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveB(n, t, a)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
