package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type Test struct {
	n     int
	deg   []int
	s     []int
	edges [][2]int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(10) + 1
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{i, p})
	}
	deg := make([]int, n)
	s := make([]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		deg[a]++
		deg[b]++
		s[a] ^= b
		s[b] ^= a
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", deg[i], s[i]))
	}
	return Test{n: n, deg: deg, s: s, edges: edges, input: sb.String()}
}

func checkOutput(t Test, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("invalid m: %v", err)
	}
	edges := make([][2]int, 0, m)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			return fmt.Errorf("invalid edge line")
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid numbers")
		}
		edges = append(edges, [2]int{a, b})
	}
	if len(edges) != m {
		return fmt.Errorf("edge count mismatch")
	}
	deg := make([]int, t.n)
	s := make([]int, t.n)
	for _, e := range edges {
		a, b := e[0], e[1]
		if a < 0 || a >= t.n || b < 0 || b >= t.n || a == b {
			return fmt.Errorf("invalid edge values")
		}
		deg[a]++
		deg[b]++
		s[a] ^= b
		s[b] ^= a
	}
	for i := 0; i < t.n; i++ {
		if deg[i] != t.deg[i] || s[i] != t.s[i] {
			return fmt.Errorf("degree/xor mismatch")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(t, out); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, t.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
