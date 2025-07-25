package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type TestA struct {
	n int
	k int
	a []int
}

func generateTests() []TestA {
	rand.Seed(1)
	tests := make([]TestA, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		k := rand.Intn(50)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(10) + 1
		}
		tests[i] = TestA{n, k, a}
	}
	return tests
}

func maxLearn(n, k int, a []int) int {
	type pair struct{ v, idx int }
	p := make([]pair, n)
	for i := 0; i < n; i++ {
		p[i] = pair{a[i], i + 1}
	}
	sort.Slice(p, func(i, j int) bool { return p[i].v < p[j].v })
	ans := 0
	for i := 0; i < n; i++ {
		if k >= p[i].v {
			k -= p[i].v
			ans++
		} else {
			break
		}
	}
	return ans
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.k)
		for j, v := range t.a {
			if j > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		// parse output
		scanner := bufio.NewScanner(strings.NewReader(got))
		if !scanner.Scan() {
			fmt.Printf("test %d: no output\n", i+1)
			os.Exit(1)
		}
		mStr := strings.TrimSpace(scanner.Text())
		m, err := strconv.Atoi(mStr)
		if err != nil {
			fmt.Printf("test %d: invalid integer output\n", i+1)
			os.Exit(1)
		}
		idxs := []int{}
		if scanner.Scan() {
			parts := strings.Fields(scanner.Text())
			for _, p := range parts {
				val, _ := strconv.Atoi(p)
				idxs = append(idxs, val)
			}
		}
		if len(idxs) != m {
			fmt.Printf("test %d: expected %d indices, got %d\n", i+1, m, len(idxs))
			os.Exit(1)
		}
		used := make(map[int]bool)
		sum := 0
		for _, idx := range idxs {
			if idx < 1 || idx > t.n {
				fmt.Printf("test %d: index %d out of range\n", i+1, idx)
				os.Exit(1)
			}
			if used[idx] {
				fmt.Printf("test %d: duplicate index %d\n", i+1, idx)
				os.Exit(1)
			}
			used[idx] = true
			sum += t.a[idx-1]
		}
		if sum > t.k {
			fmt.Printf("test %d: sum %d exceeds k %d\n", i+1, sum, t.k)
			os.Exit(1)
		}
		expectedM := maxLearn(t.n, t.k, t.a)
		if m != expectedM {
			fmt.Printf("test %d: expected m=%d got %d\n", i+1, expectedM, m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
