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

// expectedF currently matches the placeholder solution 1864F.go which
// outputs 0 for each query.
func expectedF(n int, arr []int, queries [][2]int) string {
	var sb strings.Builder
	for range queries {
		sb.WriteString("0\n")
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(10)
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	expect := expectedF(n, arr, queries)
	return sb.String(), expect
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
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
