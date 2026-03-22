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

// expected computes the number of valid (p, q) pairs for the robot problem.
// Robot 1 starts from the left, reads a[0], a[1], ... and stops at the first
// occurrence of p. Robot 2 starts from the right, reads a[n-1], a[n-2], ...
// and stops at the first occurrence of q. A pair (p, q) is valid if
// first_occurrence(p) < last_occurrence(q), i.e. robot 1 stops strictly to
// the left of robot 2.
//
// We count pairs of distinct values (p, q) where min_pos[p] < max_pos[q].
func expected(a []int) int64 {
	n := len(a)
	// num[i] = number of distinct values whose last occurrence is >= i
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
	cases = append(cases, []int{1, 5, 4, 1, 3})
	cases = append(cases, []int{1, 2, 3})
	for i := 0; i < 196; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(10) + 1
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

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput: %s", idx+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
