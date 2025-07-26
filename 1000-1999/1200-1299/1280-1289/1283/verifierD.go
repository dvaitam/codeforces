package main

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n, m int
	xs   []int
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rnd.Intn(5) + 1
		m := rnd.Intn(5) + 1
		xs := make([]int, n)
		used := map[int]bool{}
		for j := 0; j < n; j++ {
			for {
				v := rnd.Intn(51) - 25
				if !used[v] {
					used[v] = true
					xs[j] = v
					break
				}
			}
		}
		tests[i] = testCase{n, m, xs}
	}
	return tests
}

type node struct{ pos, dist int }

func bfs(tc testCase) (int64, map[int]int) {
	q := list.New()
	dist := map[int]int{}
	for _, x := range tc.xs {
		dist[x] = 0
		q.PushBack(node{x, 0})
	}
	ansDist := make(map[int]int)
	var sum int64
	for q.Len() > 0 && len(ansDist) < tc.m {
		e := q.Front()
		q.Remove(e)
		cur := e.Value.(node)
		for _, d := range []int{-1, 1} {
			np := cur.pos + d
			if _, ok := dist[np]; ok {
				continue
			}
			dist[np] = cur.dist + 1
			q.PushBack(node{np, cur.dist + 1})
			if len(ansDist) < tc.m {
				ansDist[np] = cur.dist + 1
				sum += int64(cur.dist + 1)
			}
			if len(ansDist) == tc.m {
				break
			}
		}
	}
	return sum, ansDist
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), fmt.Errorf("exec error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		for i, x := range tc.xs {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", x)
		}
		input += "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) < 2 {
			fmt.Printf("test %d: insufficient output\n", idx+1)
			os.Exit(1)
		}
		var gotSum int64
		if _, err := fmt.Sscan(lines[0], &gotSum); err != nil {
			fmt.Printf("test %d: cannot parse sum\n", idx+1)
			os.Exit(1)
		}
		fields := strings.Fields(lines[1])
		if len(fields) != tc.m {
			fmt.Printf("test %d: expected %d positions got %d\n", idx+1, tc.m, len(fields))
			os.Exit(1)
		}
		posOut := make([]int, tc.m)
		seen := map[int]bool{}
		for i, f := range fields {
			if _, err := fmt.Sscan(f, &posOut[i]); err != nil {
				fmt.Printf("test %d: invalid position\n", idx+1)
				os.Exit(1)
			}
			if seen[posOut[i]] {
				fmt.Printf("test %d: duplicate position\n", idx+1)
				os.Exit(1)
			}
			seen[posOut[i]] = true
			for _, t := range tc.xs {
				if posOut[i] == t {
					fmt.Printf("test %d: position overlaps tree\n", idx+1)
					os.Exit(1)
				}
			}
		}
		expSum, distMap := bfs(tc)
		if gotSum != expSum {
			fmt.Printf("test %d: expected sum %d got %d\n", idx+1, expSum, gotSum)
			os.Exit(1)
		}
		var calcSum int64
		for _, p := range posOut {
			d, ok := distMap[p]
			if !ok {
				d = bfsDistance(tc.xs, p)
			}
			calcSum += int64(d)
		}
		if calcSum != expSum {
			fmt.Printf("test %d: positions produce sum %d expected %d\n", idx+1, calcSum, expSum)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func bfsDistance(src []int, target int) int {
	visited := map[int]bool{}
	q := list.New()
	for _, x := range src {
		visited[x] = true
		if x == target {
			return 0
		}
		q.PushBack(node{x, 0})
	}
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		cur := e.Value.(node)
		for _, d := range []int{-1, 1} {
			np := cur.pos + d
			if visited[np] {
				continue
			}
			if np == target {
				return cur.dist + 1
			}
			visited[np] = true
			q.PushBack(node{np, cur.dist + 1})
		}
	}
	return -1
}
