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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	n   int
	arr []int
}

func genCase(rng *rand.Rand) Case {
	n := rng.Intn(8) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(3) + 1
	}
	return Case{n: n, arr: arr}
}

func isPossible(c Case) bool {
	first := c.arr[0]
	for _, v := range c.arr {
		if v != first {
			return true
		}
	}
	return false
}

func checkOutput(c Case, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("empty output")
	}
	ans := strings.ToUpper(strings.TrimSpace(scanner.Text()))
	possible := isPossible(c)
	if !possible {
		if ans != "NO" {
			return fmt.Errorf("expected NO")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("expected YES")
	}
	edges := make([][2]int, 0, c.n-1)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			return fmt.Errorf("invalid edge line: %q", scanner.Text())
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid integers")
		}
		if u < 1 || u > c.n || v < 1 || v > c.n || u == v {
			return fmt.Errorf("edge out of range")
		}
		if c.arr[u-1] == c.arr[v-1] {
			return fmt.Errorf("edge connects same gang")
		}
		edges = append(edges, [2]int{u - 1, v - 1})
	}
	if len(edges) != c.n-1 {
		return fmt.Errorf("expected %d edges got %d", c.n-1, len(edges))
	}
	// check connectivity via DSU
	parent := make([]int, c.n)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		fa, fb := find(a), find(b)
		if fa != fb {
			parent[fa] = fb
		}
	}
	for _, e := range edges {
		union(e[0], e[1])
	}
	root := find(0)
	for i := 1; i < c.n; i++ {
		if find(i) != root {
			return fmt.Errorf("edges not connected")
		}
	}
	return nil
}

func runCase(bin string, c Case) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for i, v := range c.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	return checkOutput(c, out)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := genCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
