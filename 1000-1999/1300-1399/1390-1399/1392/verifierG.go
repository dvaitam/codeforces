package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func applySegment(k int, arr []byte, l, r int, a, b []int) []byte {
	res := make([]byte, k)
	copy(res, arr)
	for i := l; i <= r; i++ {
		x := a[i]
		y := b[i]
		res[x], res[y] = res[y], res[x]
	}
	return res
}

func matches(a, b []byte) int {
	cnt := 0
	for i := range a {
		if a[i] == b[i] {
			cnt++
		}
	}
	return cnt
}

func solveCase(n, m, k int, s0, s1 string, aList, bList []int) int {
	cur := []byte(s0)
	target := []byte(s1)
	best := -1
	for l := 0; l < n; l++ {
		tray := make([]byte, k)
		copy(tray, cur)
		for r := l; r < n; r++ {
			tray = applySegment(k, tray, r, r, aList, bList)
			if r-l+1 >= m {
				val := matches(tray, target)
				if val > best {
					best = val
				}
			}
		}
	}
	return best
}

func generateTest() (string, int, int, int, string, string, []int, []int) {
	n := rand.Intn(5) + 1
	m := rand.Intn(n) + 1
	k := rand.Intn(4) + 2
	s0 := make([]byte, k)
	s1 := make([]byte, k)
	for i := 0; i < k; i++ {
		if rand.Intn(2) == 0 {
			s0[i] = '0'
		} else {
			s0[i] = '1'
		}
		if rand.Intn(2) == 0 {
			s1[i] = '0'
		} else {
			s1[i] = '1'
		}
	}
	aList := make([]int, n)
	bList := make([]int, n)
	for i := 0; i < n; i++ {
		aList[i] = rand.Intn(k)
		bList[i] = rand.Intn(k)
		for aList[i] == bList[i] {
			bList[i] = rand.Intn(k)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n%s\n%s\n", n, m, k, string(s0), string(s1))
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", aList[i]+1, bList[i]+1)
	}
	return sb.String(), n, m, k, string(s0), string(s1), aList, bList
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	rand.Seed(7)
	for i := 0; i < 100; i++ {
		input, n, m, k, s0, s1, aList, bList := generateTest()
		bestScore := solveCase(n, m, k, s0, s1, aList, bList)

		out, err := runBinary(os.Args[1], input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		lines := strings.Split(out, "\n")
		if len(lines) < 2 {
			fmt.Printf("test %d: expected 2 lines of output, got %d\n", i+1, len(lines))
			os.Exit(1)
		}
		var gotScore, gotL, gotR int
		if _, err := fmt.Sscan(lines[0], &gotScore); err != nil {
			fmt.Printf("test %d: cannot parse score: %v\n", i+1, err)
			os.Exit(1)
		}
		if _, err := fmt.Sscan(lines[1], &gotL, &gotR); err != nil {
			fmt.Printf("test %d: cannot parse l r: %v\n", i+1, err)
			os.Exit(1)
		}
		if gotScore != bestScore {
			fmt.Printf("test %d: expected score %d, got %d\n", i+1, bestScore, gotScore)
			os.Exit(1)
		}
		// Validate l, r
		if gotL < 1 || gotR > n || gotL > gotR || gotR-gotL+1 < m {
			fmt.Printf("test %d: invalid l=%d r=%d (n=%d m=%d)\n", i+1, gotL, gotR, n, m)
			os.Exit(1)
		}
		// Verify the score by applying segment [gotL-1, gotR-1]
		cur := []byte(s0)
		result := applySegment(k, cur, gotL-1, gotR-1, aList, bList)
		actualScore := matches(result, []byte(s1))
		if actualScore != gotScore {
			fmt.Printf("test %d: claimed score %d but actual score is %d for l=%d r=%d\n", i+1, gotScore, actualScore, gotL, gotR)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
