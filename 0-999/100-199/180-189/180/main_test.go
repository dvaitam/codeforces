package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type fragment struct {
	file int
	idx  int
}

var exePath = "./sol_bin"

func TestMain(m *testing.M) {
	cmd := exec.Command("g++", "-std=c++17", "sol.cpp", "-o", exePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		os.Exit(1)
	}
	code := m.Run()
	os.Remove(exePath)
	os.Exit(code)
}

func runSolution(input string) (string, error) {
	cmd := exec.Command(exePath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func verify(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	lens := make([]int, m+1)
	clusters := make([]fragment, n+1)
	used := make([]bool, n+1)
	for i := 1; i <= m; i++ {
		var ni int
		if _, err := fmt.Fscan(in, &ni); err != nil {
			return fmt.Errorf("input read ni: %v", err)
		}
		lens[i] = ni
		for j := 1; j <= ni; j++ {
			var x int
			fmt.Fscan(in, &x)
			if x < 1 || x > n {
				return fmt.Errorf("invalid cluster number %d", x)
			}
			if used[x] {
				return fmt.Errorf("cluster %d appears multiple times", x)
			}
			used[x] = true
			clusters[x] = fragment{file: i, idx: j}
		}
	}
	total := 0
	for i := 1; i <= m; i++ {
		total += lens[i]
	}

	out := bufio.NewReader(strings.NewReader(output))
	var k int
	if _, err := fmt.Fscan(out, &k); err != nil {
		return fmt.Errorf("output parse k: %v", err)
	}
	if k < 0 || k > 2*n {
		return fmt.Errorf("invalid k %d", k)
	}
	for step := 0; step < k; step++ {
		var a, b int
		if _, err := fmt.Fscan(out, &a, &b); err != nil {
			return fmt.Errorf("read op %d: %v", step+1, err)
		}
		if a < 1 || a > n || b < 1 || b > n || a == b {
			return fmt.Errorf("invalid op %d: %d %d", step+1, a, b)
		}
		clusters[b] = clusters[a]
	}

	usedFile := make([]bool, m+1)
	pos := 1
	for pos <= total {
		frag := clusters[pos]
		if frag.file == 0 || frag.idx != 1 {
			return fmt.Errorf("cluster %d does not start a file", pos)
		}
		fid := frag.file
		if fid < 1 || fid > m {
			return fmt.Errorf("unknown file id %d", fid)
		}
		if usedFile[fid] {
			return fmt.Errorf("file %d appears multiple times", fid)
		}
		usedFile[fid] = true
		for j := 1; j <= lens[fid]; j++ {
			if pos > n {
				return fmt.Errorf("ran out of clusters")
			}
			frag = clusters[pos]
			if frag.file != fid || frag.idx != j {
				return fmt.Errorf("file %d fragment %d missing at cluster %d", fid, j, pos)
			}
			pos++
		}
	}
	for i := 1; i <= m; i++ {
		if !usedFile[i] {
			return fmt.Errorf("file %d not found in final arrangement", i)
		}
	}
	if pos != total+1 {
		return fmt.Errorf("extra data in first %d clusters", total)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(199) + 2 // 2..200
	m := rng.Intn(min(n-1, 20)) + 1
	lens := make([]int, m)
	remaining := n - 1
	for i := 0; i < m; i++ {
		maxLen := min(5, remaining-(m-i-1))
		if maxLen < 1 {
			maxLen = 1
		}
		l := rng.Intn(maxLen) + 1
		lens[i] = l
		remaining -= l
	}
	perm := rng.Perm(n)
	idx := 0
	lines := make([]string, m)
	for i := 0; i < m; i++ {
		l := lens[i]
		arr := perm[idx : idx+l]
		idx += l
		nums := make([]string, l)
		for j, v := range arr {
			nums[j] = fmt.Sprintf("%d", v+1)
		}
		lines[i] = fmt.Sprintf("%d %s", l, strings.Join(nums, " "))
	}
	return fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(lines, "\n"))
}

func TestRandomCases(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	longClusters := make([]string, 199)
	for i := 0; i < 199; i++ {
		longClusters[i] = fmt.Sprintf("%d", i+1)
	}
	cases := []string{
		"2 1\n1 1\n",
		"200 1\n199 " + strings.Join(longClusters, " ") + "\n",
	}

	for _, tc := range cases {
		out, err := runSolution(tc)
		if err != nil {
			t.Fatalf("runtime error on predefined case: %v", err)
		}
		if err := verify(tc, out); err != nil {
			t.Fatalf("verification failed on predefined case: %v", err)
		}
	}

	for i := 0; i < 120; i++ {
		tc := generateCase(rng)
		out, err := runSolution(tc)
		if err != nil {
			t.Fatalf("runtime error on case %d: %v", i, err)
		}
		if err := verify(tc, out); err != nil {
			t.Fatalf("case %d verification failed: %v\ninput:\n%s\noutput:\n%s", i, err, tc, out)
		}
	}
}
