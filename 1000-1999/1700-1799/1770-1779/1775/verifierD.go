package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Node struct {
	w, d, prev int
}

func expected(n int, a []int, s, t int) (int, []int) {
	mx := 0
	for _, v := range a {
		if v > mx {
			mx = v
		}
	}
	minp := make([]int, mx+1)
	primes := []int{}
	for i := 2; i <= mx; i++ {
		if minp[i] == 0 {
			minp[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			prod := p * i
			if prod > mx {
				break
			}
			minp[prod] = p
			if p == minp[i] {
				break
			}
		}
	}
	edg := make([][]int, mx+1)
	for idx, val := range a {
		tmp := val
		for tmp > 1 {
			p := minp[tmp]
			edg[p] = append(edg[p], idx)
			for tmp%p == 0 {
				tmp /= p
			}
		}
	}
	total := n + mx + 1
	dist := make([]int, total)
	nxt := make([]int, total)
	for i := range dist {
		dist[i] = -1
		nxt[i] = -1
	}
	queue := []Node{{t, 0, -1}}
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		w, d, pr := cur.w, cur.d, cur.prev
		if dist[w] != -1 {
			continue
		}
		dist[w] = d
		nxt[w] = pr
		if w < n {
			tmp := a[w]
			for tmp > 1 {
				p := minp[tmp]
				nb := n + p
				if dist[nb] == -1 {
					queue = append(queue, Node{nb, d + 1, w})
				}
				for tmp%p == 0 {
					tmp /= p
				}
			}
		} else {
			p := w - n
			for _, idx := range edg[p] {
				if dist[idx] == -1 {
					queue = append(queue, Node{idx, d + 1, w})
				}
			}
		}
	}
	if dist[s] == -1 {
		return -1, nil
	}
	length := dist[s]/2 + 1
	path := make([]int, 0, length)
	path = append(path, s)
	for i := nxt[s]; i != -1; i = nxt[i] {
		if i < n {
			path = append(path, i)
		}
	}
	res := make([]int, len(path))
	for i, v := range path {
		res[i] = v + 1
	}
	return length, res
}

func generateTests() []struct {
	n int
	a []int
	s int
	t int
} {
	rand.Seed(3)
	t := 100
	res := make([]struct {
		n int
		a []int
		s int
		t int
	}, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(4) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(18) + 2
		}
		s := rand.Intn(n) + 1
		tpos := rand.Intn(n) + 1
		res[i] = struct {
			n int
			a []int
			s int
			t int
		}{n, arr, s, tpos}
	}
	return res
}

func verifyCase(bin string, tc struct {
	n int
	a []int
	s int
	t int
}) error {
	expLen, expPath := expected(tc.n, tc.a, tc.s-1, tc.t-1)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", tc.a[i]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.s, tc.t))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}
	outStr := strings.TrimSpace(string(out))
	if expLen == -1 {
		if outStr != "-1" {
			return fmt.Errorf("expected -1 got %s", outStr)
		}
		return nil
	}
	scanner := bufio.NewScanner(strings.NewReader(outStr))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	gotLen, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err
	}
	if gotLen != expLen {
		return fmt.Errorf("expected len %d got %d", expLen, gotLen)
	}
	path := make([]int, 0, expLen)
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}
		path = append(path, v)
	}
	if len(path) != expLen {
		return fmt.Errorf("expected %d nodes got %d", expLen, len(path))
	}
	for i, v := range path {
		if expPath[i] != v {
			return fmt.Errorf("path mismatch expected %v got %v", expPath, path)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
