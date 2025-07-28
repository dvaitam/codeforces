package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func inOrder(n int, left, right []int) []int {
	order := make([]int, 0, n)
	stack := []int{}
	curr := 1
	for curr != 0 || len(stack) > 0 {
		for curr != 0 {
			stack = append(stack, curr)
			curr = left[curr]
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, curr)
		curr = right[curr]
	}
	return order
}

func expected(n, k int, s string, left, right, parent []int) string {
	order := inOrder(n, left, right)
	chars := []byte(s)
	want := make([]bool, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && chars[order[j]-1] == chars[order[i]-1] {
			j++
		}
		next := byte('{')
		if j < n {
			next = chars[order[j]-1]
		}
		if chars[order[i]-1] < next {
			for t := i; t < j; t++ {
				want[t] = true
			}
		}
		i = j
	}
	dup := make([]bool, n+1)
	dsu := make([]int, n+1)
	for i := 0; i <= n; i++ {
		dsu[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if dsu[x] != x {
			dsu[x] = find(dsu[x])
		}
		return dsu[x]
	}
	remaining := k
	for idx := 0; idx < n && remaining > 0; idx++ {
		u := order[idx]
		if dup[u] {
			continue
		}
		if want[idx] {
			path := make([]int, 0)
			v := find(u)
			for v != 0 && !dup[v] {
				path = append(path, v)
				v = find(parent[v])
			}
			if len(path) <= remaining {
				for _, x := range path {
					dup[x] = true
					dsu[x] = find(parent[x])
				}
				remaining -= len(path)
			}
		}
	}
	var res []byte
	stack := []int{}
	curr := 1
	for curr != 0 || len(stack) > 0 {
		for curr != 0 {
			stack = append(stack, curr)
			curr = left[curr]
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ch := chars[curr-1]
		res = append(res, ch)
		if dup[curr] {
			res = append(res, ch)
		}
		curr = right[curr]
	}
	return string(res)
}

func generateTree(n int, rng *rand.Rand) (left, right, parent []int) {
	left = make([]int, n+1)
	right = make([]int, n+1)
	parent = make([]int, n+1)
	for child := 2; child <= n; child++ {
		for {
			p := rng.Intn(child-1) + 1
			if rng.Intn(2) == 0 {
				if left[p] == 0 {
					left[p] = child
					parent[child] = p
					break
				}
			} else {
				if right[p] == 0 {
					right[p] = child
					parent[child] = p
					break
				}
			}
		}
	}
	return
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 1
		k := rng.Intn(n) + 1
		letters := make([]byte, n)
		for i := 0; i < n; i++ {
			letters[i] = byte('a' + rng.Intn(3))
		}
		left, right, parent := generateTree(n, rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		sb.WriteString(fmt.Sprintf("%s\n", string(letters)))
		for i := 1; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", left[i], right[i]))
		}
		input := sb.String()
		exp := expected(n, k, string(letters), left, right, parent)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
