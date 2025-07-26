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

func reachable(limit, a, b int) int {
	vis := make([]bool, limit+1)
	q := []int{0}
	vis[0] = true
	for front := 0; front < len(q); front++ {
		u := q[front]
		v := u + a
		if v <= limit && !vis[v] {
			vis[v] = true
			q = append(q, v)
		}
		v = u - b
		if v >= 0 && !vis[v] {
			vis[v] = true
			q = append(q, v)
		}
	}
	cnt := 0
	for _, ok := range vis {
		if ok {
			cnt++
		}
	}
	return cnt
}

func expected(m, a, b int) int {
	sum := 0
	for i := 0; i <= m; i++ {
		sum += reachable(i, a, b)
	}
	return sum
}

func genCase(rng *rand.Rand) (string, int) {
	m := rng.Intn(20) + 1
	a := rng.Intn(5) + 1
	b := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", m, a, b)
	return sb.String(), expected(m, a, b)
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
