package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(n int, arr []int64) string {
	last := n - 1
	for last >= 0 && arr[last] == 0 {
		last--
	}
	var ans int64
	for i := 0; i < last; i++ {
		ans += arr[i]
		if arr[i] == 0 {
			ans++
		}
	}
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(18) + 2
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(10)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	exp := expected(n, append([]int64(nil), a...))
	return sb.String(), exp
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1705))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
