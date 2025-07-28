package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func solve(arr []int64) int64 {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	ans := arr[0]
	for i := 1; i < len(arr); i++ {
		diff := arr[i] - arr[i-1]
		if diff > ans {
			ans = diff
		}
	}
	return ans
}

func randCase() []int64 {
	n := rand.Intn(50) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rand.Int63n(2_000_000_001) - 1_000_000_000
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	rand.Seed(3)
	cases := [][]int64{
		{10},
		{-5, -5},
		{-1, 0, 1},
	}
	for len(cases) < 100 {
		cases = append(cases, randCase())
	}

	for i, arr := range cases {
		input := fmt.Sprintf("1\n%d\n", len(arr))
		for j, v := range arr {
			if j+1 == len(arr) {
				input += fmt.Sprintf("%d\n", v)
			} else {
				input += fmt.Sprintf("%d ", v)
			}
		}
		expected := solve(append([]int64(nil), arr...))
		gotStr, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr = strings.TrimSpace(gotStr)
		if gotStr != fmt.Sprint(expected) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %d\ngot: %s\n", i+1, input, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
