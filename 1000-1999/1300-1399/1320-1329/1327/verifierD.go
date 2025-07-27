package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func divisors(n int) []int {
	ds := []int{}
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			ds = append(ds, i)
			if i != n/i {
				ds = append(ds, n/i)
			}
		}
	}
	sort.Ints(ds)
	return ds
}

func solve(n int, p []int, c []int) int {
	visited := make([]bool, n)
	ans := n

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		var cycle []int
		j := i
		for !visited[j] {
			visited[j] = true
			cycle = append(cycle, j)
			j = p[j]
		}
		l := len(cycle)
		ds := divisors(l)
	NextCycle:
		for _, d := range ds {
			if d >= ans {
				continue
			}
			for start := 0; start < d; start++ {
				color := c[cycle[start]]
				good := true
				for pos := start; pos < l; pos += d {
					if c[cycle[pos]] != color {
						good = false
						break
					}
				}
				if good {
					if d < ans {
						ans = d
					}
					break NextCycle
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(4)
	t := 100
	ns := make([]int, t)
	ps := make([][]int, t)
	cs := make([][]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		ns[i] = n
		perm := rand.Perm(n)
		p := make([]int, n)
		for j := 0; j < n; j++ {
			p[j] = perm[j]
		}
		colors := make([]int, n)
		for j := 0; j < n; j++ {
			colors[j] = rand.Intn(n) + 1
		}
		ps[i] = p
		cs[i] = colors
	}

	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := ns[i]
		input.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", ps[i][j]+1))
		}
		input.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", cs[i][j]))
		}
		input.WriteByte('\n')
	}
	in := input.String()

	var expectedOut strings.Builder
	for i := 0; i < t; i++ {
		res := solve(ns[i], ps[i], cs[i])
		expectedOut.WriteString(fmt.Sprintf("%d\n", res))
	}
	want := strings.TrimSpace(expectedOut.String())

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Runtime error: %v\n%s", err, out.String())
		os.Exit(1)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := strings.Split(want, "\n")
	if len(gotLines) != len(wantLines) {
		fmt.Println("Wrong answer: line count mismatch")
		os.Exit(1)
	}
	for i := range wantLines {
		if strings.TrimSpace(gotLines[i]) != strings.TrimSpace(wantLines[i]) {
			fmt.Printf("Wrong answer on test %d: expected %s got %s\n", i+1, wantLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
