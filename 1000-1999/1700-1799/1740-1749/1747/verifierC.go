package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type testCase struct {
	n      int
	arr    []int
	winner string
}

func generateTests() (string, []testCase) {
	rand.Seed(1)
	t := 100
	var buf bytes.Buffer
	var cases []testCase
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 2 // n in [2,6]
		fmt.Fprintln(&buf, n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(9) + 1
			fmt.Fprint(&buf, arr[j])
			if j+1 < n {
				fmt.Fprint(&buf, " ")
			}
		}
		fmt.Fprintln(&buf)
		winner := "Bob"
		mn := arr[0]
		for _, v := range arr[1:] {
			if v < mn {
				winner = "Alice"
				break
			}
		}
		cases = append(cases, testCase{n: n, arr: arr, winner: winner})
	}
	return buf.String(), cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	input, cases := generateTests()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for idx, tc := range cases {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		got := scanner.Text()
		if got != tc.winner {
			fmt.Printf("case %d: expected %s, got %s\n", idx+1, tc.winner, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
