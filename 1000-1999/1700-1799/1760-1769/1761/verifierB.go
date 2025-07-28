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

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solve(n int, a []int) int {
	if n == 1 {
		return 1
	}
	alt := a[0] != a[1]
	if alt {
		for i := 2; i < n; i++ {
			if a[i] != a[i%2] {
				alt = false
				break
			}
		}
	}
	if alt {
		return (n + 3) / 2
	}
	return n
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	a := make([]int, n)
	if n == 1 {
		a[0] = rng.Intn(n) + 1
	} else {
		a[0] = rng.Intn(n) + 1
		for i := 1; i < n; i++ {
			v := rng.Intn(n) + 1
			for v == a[i-1] || (i == n-1 && v == a[0]) {
				v = rng.Intn(n) + 1
			}
			a[i] = v
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", solve(n, a))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
