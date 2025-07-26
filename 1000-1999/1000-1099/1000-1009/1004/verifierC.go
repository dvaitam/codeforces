package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(a []int) int64 {
	n := len(a)
	num := make([]int, n+1)
	seen := map[int]bool{}
	for i := n - 1; i >= 0; i-- {
		if !seen[a[i]] {
			num[i] = num[i+1] + 1
			seen[a[i]] = true
		} else {
			num[i] = num[i+1]
		}
	}
	seen = map[int]bool{}
	var ans int64
	for i := 0; i < n-1; i++ {
		if !seen[a[i]] {
			ans += int64(num[i+1])
			seen[a[i]] = true
		}
	}
	return ans
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases [][]int
	cases = append(cases, []int{1})
	cases = append(cases, []int{1, 2, 1})
	for i := 0; i < 98; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(10)
		}
		cases = append(cases, arr)
	}
	for idx, arr := range cases {
		input := fmt.Sprintf("%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := fmt.Sprintf("%d", expected(arr))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
