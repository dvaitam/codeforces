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

func solveCase(n, m, k int, s0, s1 string, aList, bList []int) (int, int, int) {
	cur := []byte(s0)
	target := []byte(s1)
	best := -1
	bestL, bestR := 1, m
	for l := 0; l < n; l++ {
		tray := make([]byte, k)
		copy(tray, cur)
		for r := l; r < n; r++ {
			tray = applySegment(k, tray, r, r, aList, bList)
			if r-l+1 >= m {
				val := matches(tray, target)
				if val > best {
					best = val
					bestL = l + 1
					bestR = r + 1
				}
			}
		}
	}
	return best, bestL, bestR
}

func generateTest() (string, string) {
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

	best, l, r := solveCase(n, m, k, string(s0), string(s1), aList, bList)
	exp := fmt.Sprintf("%d\n%d %d\n", best, l, r)
	return sb.String(), exp
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	for i := 0; i < t; i++ {
		ti, to := generateTest()
		in.WriteString(ti)
		out.WriteString(to)
	}
	return in.String(), out.String()
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
	in, exp := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
