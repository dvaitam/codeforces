package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildOfficial() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	bin := filepath.Join(dir, "officialE.bin")
	cmd := exec.Command("go", "build", "-o", bin, filepath.Join(dir, "1061E.go"))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func genCase(rng *rand.Rand) (int, int, int, []int, [][2]int, [][2]int, [][2]int, [][2]int) {
	n := rng.Intn(6) + 2
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(100) + 1
	}
	t1 := genTree(rng, n)
	t2 := genTree(rng, n)
	q1 := rng.Intn(n) + 1
	q2 := rng.Intn(n) + 1
	dem1 := make([][2]int, q1)
	used := make(map[int]bool)
	used[x] = true
	dem1[0] = [2]int{x, rng.Intn(n) + 1}
	for i := 1; i < q1; i++ {
		k := rng.Intn(n) + 1
		for used[k] {
			k = rng.Intn(n) + 1
		}
		used[k] = true
		dem1[i] = [2]int{k, rng.Intn(n) + 1}
	}
	dem2 := make([][2]int, q2)
	used = map[int]bool{y: true}
	dem2[0] = [2]int{y, rng.Intn(n) + 1}
	for i := 1; i < q2; i++ {
		k := rng.Intn(n) + 1
		for used[k] {
			k = rng.Intn(n) + 1
		}
		used[k] = true
		dem2[i] = [2]int{k, rng.Intn(n) + 1}
	}
	return n, x, y, a, t1, t2, dem1, dem2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	off, err := buildOfficial()
	if err != nil {
		fmt.Println("failed to build official solution:", err)
		return
	}
	defer os.Remove(off)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, x, y, a, t1, t2, d1, d2 := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, e := range t1 {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, e := range t2 {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(d1)))
		for _, d := range d1 {
			sb.WriteString(fmt.Sprintf("%d %d\n", d[0], d[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(d2)))
		for _, d := range d2 {
			sb.WriteString(fmt.Sprintf("%d %d\n", d[0], d[1]))
		}
		input := sb.String()
		exp, err1 := runBinary(off, input)
		out, err2 := runBinary(candidate, input)
		if err1 != nil || err2 != nil {
			fmt.Printf("Runtime error on test %d\n", i+1)
			if err1 != nil {
				fmt.Println("official:", err1)
			}
			if err2 != nil {
				fmt.Println("candidate:", err2)
			}
			fmt.Println("input:\n" + input)
			return
		}
		if exp != out {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, input, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}
