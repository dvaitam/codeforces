package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type caseE struct {
	n int
	s string
}

func parseCases(path string) ([]caseE, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	cases := []caseE{}
	for {
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		if !sc.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		s := strings.TrimSpace(sc.Text())
		cases = append(cases, caseE{n: n, s: s})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func bfs(start int, s string, pos map[byte][]int) []int {
	n := len(s)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	used := make(map[byte]bool)
	for head := 0; head < len(q); head++ {
		u := q[head]
		d := dist[u]
		if u > 0 && dist[u-1] == -1 {
			dist[u-1] = d + 1
			q = append(q, u-1)
		}
		if u+1 < n && dist[u+1] == -1 {
			dist[u+1] = d + 1
			q = append(q, u+1)
		}
		c := s[u]
		if !used[c] {
			used[c] = true
			for _, v := range pos[c] {
				if dist[v] == -1 {
					dist[v] = d + 1
					q = append(q, v)
				}
			}
		}
	}
	return dist
}

func solve(tc caseE) (int, int64) {
	pos := make(map[byte][]int)
	for i := 0; i < tc.n; i++ {
		ch := tc.s[i]
		pos[ch] = append(pos[ch], i)
	}
	best := -1
	var cnt int64
	for i := 0; i < tc.n; i++ {
		dist := bfs(i, tc.s, pos)
		for j := i + 1; j < tc.n; j++ {
			d := dist[j]
			if d > best {
				best = d
				cnt = 1
			} else if d == best {
				cnt++
			}
		}
	}
	return best, cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesE.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		d, c := solve(tc)
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "case %d: incomplete output\n", idx+1)
			os.Exit(1)
		}
		gotD, _ := strconv.Atoi(fields[0])
		gotC, _ := strconv.ParseInt(fields[1], 10, 64)
		if gotD != d || gotC != c {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %d %d\n", idx+1, d, c, gotD, gotC)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
