package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(n int, b []int) string {
	res := make([]int, 0, 2*n)
	res = append(res, b[0])
	for i := 1; i < n; i++ {
		if b[i] < b[i-1] {
			res = append(res, b[i])
		}
		res = append(res, b[i])
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprint(len(res)))
	sb.WriteByte('\n')
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genCases() []string {
	rand.Seed(2)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(50) + 1
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
		parts := strings.Fields(lines[2])
		b := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &b[j])
		}
		want := solveB(n, b)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
