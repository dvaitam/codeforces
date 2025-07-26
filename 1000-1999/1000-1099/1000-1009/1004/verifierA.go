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

func expected(n int, d int, a []int) int {
	count := 2
	for i := 0; i < n-1; i++ {
		gap := a[i+1] - a[i]
		if gap > 2*d {
			count += 2
		} else if gap == 2*d {
			count += 1
		}
	}
	return count
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	type test struct {
		n, d int
		a    []int
	}
	var cases []test
	// add deterministic edge cases
	cases = append(cases, test{1, 1, []int{0}})
	cases = append(cases, test{2, 1, []int{-2, 2}})
	cases = append(cases, test{3, 2, []int{-5, 0, 5}})
	// generate random cases
	for i := 0; i < 97; i++ {
		n := rng.Intn(10) + 1
		d := rng.Intn(10) + 1
		set := map[int]bool{}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			for {
				x := rng.Intn(200) - 100
				if !set[x] {
					set[x] = true
					arr[j] = x
					break
				}
			}
		}
		sort.Ints(arr)
		cases = append(cases, test{n, d, arr})
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.d)
		for i, v := range tc.a {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := fmt.Sprintf("%d", expected(tc.n, tc.d, tc.a))
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
