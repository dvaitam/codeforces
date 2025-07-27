package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
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
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	input    string
	expected []int
}

func solveCase(a, b []int) int {
	n := len(a)
	pairs := make([]struct{ a, b int }, n)
	for i := 0; i < n; i++ {
		pairs[i].a = a[i]
		pairs[i].b = b[i]
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].a < pairs[j].a })
	sumB := 0
	for i := 0; i < n; i++ {
		sumB += pairs[i].b
	}
	ans := sumB
	for i := 0; i < n; i++ {
		sumB -= pairs[i].b
		time := sumB
		if pairs[i].a > time {
			time = pairs[i].a
		}
		if time < ans {
			ans = time
		}
	}
	return ans
}

func genTest(rng *rand.Rand) Test {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]int, 0, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		a := make([]int, n)
		b := make([]int, n)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(100)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			b[j] = rng.Intn(100)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(b[j]))
		}
		sb.WriteByte('\n')
		expected = append(expected, solveCase(a, b))
	}
	return Test{input: sb.String(), expected: expected}
}

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		nums, err := parseInts(out)
		if err != nil {
			fmt.Printf("test %d bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if len(nums) != len(tc.expected) {
			fmt.Printf("test %d expected %d numbers got %d\ninput:\n%s", i+1, len(tc.expected), len(nums), tc.input)
			os.Exit(1)
		}
		for j, v := range tc.expected {
			if nums[j] != v {
				fmt.Printf("test %d failed at case %d\ninput:\n%s\nexpected %d got %d\n", i+1, j+1, tc.input, v, nums[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
