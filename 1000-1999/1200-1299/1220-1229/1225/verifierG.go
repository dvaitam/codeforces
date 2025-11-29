package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `4 3 24 26 8 7
2 4 20 30
4 7 23 30 24 11
5 6 14 12 22 8 3
5 5 21 9 21 1 28
3 2 3 6 5
5 4 24 25 16 3 18
4 6 19 12 15 9
3 2 26 24 25
2 7 19 25
3 3 7 20 14
4 7 12 16 19 8
5 2 30 7 30 24 6
5 2 25 23 28 9 11
2 6 6 2
4 4 10 14 22 15
4 3 20 14 8 24
4 2 7 23 15 24
2 4 14 19
5 2 19 24 9 21 9
5 4 7 5 17 26 24
3 2 29 9 1
2 7 9 20
5 7 15 4 20 2 21
5 7 14 5 17 29 25
4 5 30 18 24 11
2 7 3 22
4 2 27 21 25 26
3 7 18 28 30
5 7 18 2 18 15 17
2 6 26 21
3 6 6 26 16
4 7 2 14 11 28
3 2 4 12 2
3 5 24 24 10
3 2 11 20 23
3 2 10 8 27
2 5 22 20
4 2 6 1 14 8
5 2 24 24 30 24 4
4 2 3 2 20 14
2 2 10 9
4 3 20 6 27 4
5 7 6 18 16 11 18
4 7 18 17 15 11
3 7 17 10 30
2 2 7 27
2 4 9 21
5 7 12 9 12 12 29
5 5 21 23 10 3 10
4 7 15 2 11 23
4 2 23 7 12 24
3 5 12 1 18
4 2 13 11 15 8
4 7 25 6 13 17
5 4 1 19 25 22 4
2 3 7 1
5 6 4 11 7 26 21
4 6 22 1 20 19
2 2 30 16
5 2 18 9 10 21 24
3 2 8 1 1
3 3 6 27 20
2 2 12 1
2 4 11 21
4 3 23 9 29 18
5 3 20 14 14 24 14
5 2 25 10 12 30 26
4 4 15 19 8 21
2 4 21 21
4 5 16 9 19 25
5 2 28 9 19 18 7
3 3 9 1 16
2 3 15 19
2 3 8 15
5 2 19 14 5 3 6
2 7 23 24
3 3 20 12 26
2 6 15 10
2 7 25 15
2 3 8 11
2 3 7 25
5 6 17 8 14 27 14
3 3 11 9 19
5 4 15 12 15 23 18
5 2 2 20 20 11 23
2 4 1 22
4 6 6 24 8 2
2 2 25 4
5 3 13 17 12 25 21
3 3 28 4 13
2 2 7 5
3 2 13 7 16
3 6 14 2 14
2 6 18 12
5 2 16 21 9 28 26
5 2 24 23 10 14 13
5 2 25 24 20 17 6
4 6 16 3 3 9
2 4 27 10
4 4 25 28 25 15`

type parent struct {
	prevKey int
	idx     int
}

type testCase struct {
	n int
	k int
	a []int
}

func reduce(x, k int) int {
	for x%k == 0 {
		x /= k
	}
	return x
}

// Embedded reference logic from 1225G.go.
func solve(n, k int, a []int) string {
	sum := 0
	for _, v := range a {
		sum += v
	}
	maxVal := sum
	fullMask := 1<<n - 1
	visited := make(map[int]parent)
	type queueEntry struct {
		mask int
		val  int
	}
	queue := make([]queueEntry, 0)
	key := func(mask, val int) int { return mask*(maxVal+1) + val }

	for i := 0; i < n; i++ {
		m := 1 << i
		v := a[i]
		k0 := key(m, v)
		visited[k0] = parent{prevKey: -1, idx: i}
		queue = append(queue, queueEntry{mask: m, val: v})
	}

	finalKey := -1
	for front := 0; front < len(queue); front++ {
		cur := queue[front]
		if cur.mask == fullMask && cur.val == 1 {
			finalKey = key(cur.mask, cur.val)
			break
		}
		for j := 0; j < n; j++ {
			if cur.mask&(1<<j) != 0 {
				continue
			}
			newMask := cur.mask | (1 << j)
			newVal := reduce(cur.val+a[j], k)
			k2 := key(newMask, newVal)
			if _, ok := visited[k2]; !ok {
				visited[k2] = parent{prevKey: key(cur.mask, cur.val), idx: j}
				queue = append(queue, queueEntry{mask: newMask, val: newVal})
			}
		}
	}

	if finalKey == -1 {
		return "NO"
	}

	order := make([]int, 0, n)
	curKey := finalKey
	for {
		info := visited[curKey]
		order = append(order, info.idx)
		if info.prevKey == -1 {
			break
		}
		curKey = info.prevKey
	}
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	res := a[order[0]]
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 1; i < len(order); i++ {
		x := res
		y := a[order[i]]
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		res = reduce(x+y, k)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err1 := strconv.Atoi(parts[0])
		k, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid n or k", idx+1)
		}
		if len(parts) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(parts)-2)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid element", idx+1)
			}
			a[i] = v
		}
		res = append(res, testCase{n: n, k: k, a: a})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		for j, v := range tc.a {
			if j > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"

		expected := solve(tc.n, tc.k, tc.a)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", i+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
