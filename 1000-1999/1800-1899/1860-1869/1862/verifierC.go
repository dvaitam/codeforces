package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func isSymmetric(a []int) bool {
	n := len(a)
	if a[0] != n {
		return false
	}
	j := n - 1
	for i := 1; i <= n; i++ {
		for j >= 0 && a[j] < i {
			j--
		}
		if j+1 != a[i-1] {
			return false
		}
	}
	return true
}

func genCases() []string {
	rand.Seed(3)
	cases := make([]string, 100)
	for idx := 0; idx < 100; idx++ {
		n := rand.Intn(7) + 1
		a := make([]int, n)
		a[0] = rand.Intn(10) + n
		for i := 1; i < n; i++ {
			dec := rand.Intn(3)
			val := a[i-1] - dec
			if val < 1 {
				val = 1
			}
			a[i] = val
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		cases[idx] = sb.String()
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		parts := strings.Fields(lines[2])
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &a[j])
		}
		want := "NO"
		if isSymmetric(a) {
			want = "YES"
		}
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
