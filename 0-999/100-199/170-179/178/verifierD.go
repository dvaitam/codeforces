package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testD struct {
	n    int
	nums []int64
}

func genTests() []testD {
	rand.Seed(42)
	var tests []testD
	base := []int64{8, 1, 6, 3, 5, 7, 4, 9, 2}
	for i := 0; i < 100; i++ {
		switch i % 3 {
		case 0:
			n := 1
			val := int64(rand.Intn(10))
			tests = append(tests, testD{n, []int64{val}})
		case 1:
			n := 2
			val := int64(rand.Intn(5))
			tests = append(tests, testD{n, []int64{val, val, val, val}})
		default:
			n := 3
			nums := append([]int64(nil), base...)
			rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
			tests = append(tests, testD{n, nums})
		}
	}
	return tests
}

func inputString(t testD) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func check(t testD, out string) bool {
	reader := strings.NewReader(out)
	var s int64
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return false
	}
	grid := make([][]int64, t.n)
	counts := make(map[int64]int)
	for _, v := range t.nums {
		counts[v]++
	}
	for i := 0; i < t.n; i++ {
		row := make([]int64, t.n)
		for j := 0; j < t.n; j++ {
			if _, err := fmt.Fscan(reader, &row[j]); err != nil {
				return false
			}
			counts[row[j]]--
		}
		grid[i] = row
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	for i := 0; i < t.n; i++ {
		sum := int64(0)
		for j := 0; j < t.n; j++ {
			sum += grid[i][j]
		}
		if sum != s {
			return false
		}
	}
	for j := 0; j < t.n; j++ {
		sum := int64(0)
		for i := 0; i < t.n; i++ {
			sum += grid[i][j]
		}
		if sum != s {
			return false
		}
	}
	sum := int64(0)
	for i := 0; i < t.n; i++ {
		sum += grid[i][i]
	}
	if sum != s {
		return false
	}
	sum = 0
	for i := 0; i < t.n; i++ {
		sum += grid[i][t.n-1-i]
	}
	if sum != s {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		in := inputString(t)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if !check(t, out) {
			fmt.Printf("test %d failed\ninput:\n%soutput:\n%s\n", i+1, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
