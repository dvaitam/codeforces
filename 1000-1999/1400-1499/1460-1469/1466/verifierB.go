package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(n int, arr []int) string {
	used := make(map[int]bool)
	for _, x := range arr {
		if !used[x] {
			used[x] = true
		} else if !used[x+1] {
			used[x+1] = true
		}
	}
	return fmt.Sprint(len(used))
}

func genCases() []string {
	rand.Seed(2)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(2*n) + 1
		}
		// sort arr
		for j := 0; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if arr[k] < arr[j] {
					arr[k], arr[j] = arr[j], arr[k]
				}
			}
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(arr[j]))
		}
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		arr := make([]int, n)
		parts := strings.Fields(lines[2])
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &arr[j])
		}
		want := solveB(n, arr)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
