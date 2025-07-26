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

func expectedA(input string) (string, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return "", fmt.Errorf("invalid input")
	}
	var n int
	fmt.Sscan(sc.Text(), &n)
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		if !sc.Scan() {
			return "", fmt.Errorf("invalid edge u")
		}
		var u int
		fmt.Sscan(sc.Text(), &u)
		if !sc.Scan() {
			return "", fmt.Errorf("invalid edge v")
		}
		var v int
		fmt.Sscan(sc.Text(), &v)
		deg[u]++
		deg[v]++
	}
	for i := 1; i <= n; i++ {
		if deg[i] == 2 {
			return "NO", nil
		}
	}
	return "YES", nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(18) + 2 // 2..19
	edges := make([][2]int, n-1)
	if rng.Intn(2) == 0 {
		// path
		for i := 2; i <= n; i++ {
			edges[i-2] = [2]int{i - 1, i}
		}
	} else {
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-2] = [2]int{p, i}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	exp, _ := expectedA(input)
	return input, exp
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
