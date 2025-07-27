package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveI(arr []int) string {
	mx := 0
	for _, v := range arr {
		if v > mx {
			mx = v
		}
	}
	return fmt.Sprint(mx)
}

func genCases() []string {
	rand.Seed(9)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		b := 20
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(1 << uint(b))
		}
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, b))
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
		fmt.Println("Usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n, b int
		fmt.Sscan(lines[0], &n, &b)
		arr := make([]int, n)
		parts := strings.Fields(lines[1])
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &arr[j])
		}
		want := solveI(arr)
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
