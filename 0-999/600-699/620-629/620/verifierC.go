package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func maxSegments(arr []int) int {
	n := len(arr)
	coords := append([]int(nil), arr...)
	sort.Ints(coords)
	uniq := coords[:0]
	for i := 0; i < n; {
		v := coords[i]
		uniq = append(uniq, v)
		j := i + 1
		for j < n && coords[j] == v {
			j++
		}
		i = j
	}
	posMap := make(map[int]int, len(uniq))
	for idx, v := range uniq {
		posMap[v] = idx
	}
	dp := make([]int, n+1)
	last := make([]int, len(uniq))
	for i := 1; i <= n; i++ {
		dp[i] = dp[i-1]
		pos := posMap[arr[i-1]]
		if last[pos] != 0 {
			ptr := last[pos]
			if dp[ptr-1]+1 > dp[i] {
				dp[i] = dp[ptr-1] + 1
			}
		}
		last[pos] = i
	}
	return dp[n]
}

func isGoodSegment(arr []int, l, r int) bool {
	seen := make(map[int]bool)
	for i := l; i <= r; i++ {
		v := arr[i]
		if seen[v] {
			return true
		}
		seen[v] = true
	}
	return false
}

func verify(arr []int, out string) error {
	expected := maxSegments(arr)
	out = strings.TrimSpace(out)
	if expected == 0 {
		if out != "-1" {
			return fmt.Errorf("expected -1 got %q", out)
		}
		return nil
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	if k != expected {
		return fmt.Errorf("expected %d segments got %d", expected, k)
	}
	segs := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &segs[i][0], &segs[i][1]); err != nil {
			return fmt.Errorf("failed to read segment %d: %v", i+1, err)
		}
	}
	if _, err := fmt.Fscan(reader); err == nil {
		return fmt.Errorf("extra output")
	}
	n := len(arr)
	if segs[0][0] != 1 || segs[k-1][1] != n {
		return fmt.Errorf("segments must cover the array")
	}
	for i := 0; i < k; i++ {
		l, r := segs[i][0], segs[i][1]
		if l > r || l < 1 || r > n {
			return fmt.Errorf("segment %d out of range", i+1)
		}
		if i > 0 && l != segs[i-1][1]+1 {
			return fmt.Errorf("segments not contiguous")
		}
		if !isGoodSegment(arr, l-1, r-1) {
			return fmt.Errorf("segment %d is not good", i+1)
		}
	}
	return nil
}

func runCase(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := [][]int{
		{1, 2, 1, 2},
		{1, 1, 1, 1},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(5) + 1
		}
		tests = append(tests, arr)
	}

	for i, arr := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if err := verify(arr, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, sb.String(), out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
