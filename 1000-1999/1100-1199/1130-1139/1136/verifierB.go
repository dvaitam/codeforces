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

func expected(n, k int) string {
	if k-1 < n-k {
		return fmt.Sprintf("%d", 3*n+k-1)
	}
	return fmt.Sprintf("%d", 3*n+n-k)
}

func runCase(bin, input, want string) error {
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
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %s got %s", want, got)
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

	type test struct{ n, k int }
	tests := []test{
		{n: 2, k: 2},
		{n: 5, k: 3},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(5000-2+1) + 2
		k := rng.Intn(n) + 1
		tests = append(tests, test{n: n, k: k})
	}
	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		want := expected(tc.n, tc.k)
		if err := runCase(bin, input, want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
