package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func ceilDiv(a int64, b int64) int64 {
	if b == 0 {
		return 0
	}
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

func feasible(teachers []int, req []int64) bool {
	sort.Ints(teachers)
	reqInts := make([]int64, len(req))
	copy(reqInts, req)
	sort.Slice(reqInts, func(i, j int) bool { return reqInts[i] < reqInts[j] })
	off := len(teachers) - len(reqInts)
	if off < 0 {
		return false
	}
	for i := 0; i < len(reqInts); i++ {
		if int64(teachers[off+i]) < reqInts[i] {
			return false
		}
	}
	return true
}

func solveCaseE(n, m int, teachers []int, groups [][]int) string {
	req := make([]int64, m)
	sums := make([]int64, m)
	sizes := make([]int, m)
	for i := 0; i < m; i++ {
		sum := int64(0)
		for _, v := range groups[i] {
			sum += int64(v)
		}
		sums[i] = sum
		sizes[i] = len(groups[i])
		req[i] = ceilDiv(sum, int64(len(groups[i])))
	}
	res := ""
	for i := 0; i < m; i++ {
		for _, age := range groups[i] {
			newSum := sums[i] - int64(age)
			newSize := sizes[i] - 1
			newReq := ceilDiv(newSum, int64(newSize))
			old := req[i]
			req[i] = newReq
			if feasible(append([]int(nil), teachers...), req) {
				res += "1"
			} else {
				res += "0"
			}
			req[i] = old
		}
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	m := rng.Intn(n-1) + 1
	teachers := make([]int, n)
	for i := range teachers {
		teachers[i] = rng.Intn(10) + 1
	}
	groups := make([][]int, m)
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	for _, t := range teachers {
		input += fmt.Sprintf("%d ", t)
	}
	input = strings.TrimSpace(input) + "\n"
	for i := 0; i < m; i++ {
		k := rng.Intn(3) + 2
		groups[i] = make([]int, k)
		input += fmt.Sprintf("%d ", k)
		for j := 0; j < k; j++ {
			age := rng.Intn(10) + 1
			groups[i][j] = age
			input += fmt.Sprintf("%d ", age)
		}
		input = strings.TrimSpace(input) + "\n"
	}
	exp := solveCaseE(n, m, teachers, groups)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
