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

type testCaseG struct {
	n   int
	arr []int
}

func generateTestsG(num int) []testCaseG {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseG, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(7) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(n) + 1
		}
		tests[i] = testCaseG{n: n, arr: arr}
	}
	return tests
}

func bfs(arr []int) int {
	n := len(arr)
	type state string
	toKey := func(a []int) state {
		b := make([]byte, 0, n*4)
		for i, v := range a {
			if i > 0 {
				b = append(b, ',')
			}
			b = append(b, fmt.Sprint(v)...)
		}
		return state(b)
	}
	start := make([]int, n)
	copy(start, arr)
	q := [][]int{start}
	vis := map[state]struct{}{toKey(start): {}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for i := 0; i < n; i++ {
			j := cur[i] - 1
			if j < 0 || j >= n {
				continue
			}
			nxt := make([]int, n)
			copy(nxt, cur)
			nxt[i], nxt[j] = nxt[j], nxt[i]
			key := toKey(nxt)
			if _, ok := vis[key]; !ok {
				vis[key] = struct{}{}
				q = append(q, nxt)
			}
		}
	}
	return len(vis)
}

func solveG(tc testCaseG) string {
	if tc.n > 9 {
		return "0"
	}
	ans := bfs(tc.arr)
	const mod = 1000000007
	return fmt.Sprint(ans % mod)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsG(100)
	var input bytes.Buffer
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := solveG(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
