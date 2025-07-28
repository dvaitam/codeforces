package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func processArray(arr []int) int {
	sort.Ints(arr)
	arr = unique(arr)
	for len(arr) > 1 {
		n := len(arr)
		for i := 0; i < n; i++ {
			arr[i] += n - i
		}
		sort.Ints(arr)
		arr = unique(arr)
	}
	return arr[0]
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func solveG(n int, a []int, ops [][2]int) string {
	res := make([]string, len(ops))
	for i, op := range ops {
		a[op[0]] = op[1]
		arrCopy := make([]int, len(a))
		copy(arrCopy, a)
		val := processArray(arrCopy)
		res[i] = fmt.Sprint(val)
	}
	return strings.Join(res, "\n")
}

func genCases() []string {
	rand.Seed(7)
	cases := make([]string, 100)
	for idx := 0; idx < 100; idx++ {
		n := rand.Intn(4) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(10) + 1
		}
		q := rand.Intn(4) + 1
		ops := make([][2]int, q)
		for i := 0; i < q; i++ {
			ops[i][0] = rand.Intn(n)
			ops[i][1] = rand.Intn(10) + 1
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(a[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprint(q))
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", ops[i][0]+1, ops[i][1]))
		}
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
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		idx := 0
		var n int
		fmt.Sscan(lines[1], &n)
		idx = 2
		parts := strings.Fields(lines[idx])
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &a[j])
		}
		idx++
		var q int
		fmt.Sscan(lines[idx], &q)
		idx++
		ops := make([][2]int, q)
		for j := 0; j < q; j++ {
			fmt.Sscan(lines[idx+j], &ops[j][0], &ops[j][1])
			ops[j][0]--
		}
		want := solveG(n, append([]int(nil), a...), ops)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
