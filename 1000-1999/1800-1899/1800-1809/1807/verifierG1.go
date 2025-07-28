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

func runProg(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	arr []int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1807))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(8) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(10) + 1
		}
		cases[i] = Case{arr: a}
	}
	return cases
}

func expected(a []int) string {
	b := make([]int, len(a))
	copy(b, a)
	sort.Ints(b)
	if b[0] != 1 {
		return "NO"
	}
	sum := 1
	for i := 1; i < len(b); i++ {
		if b[i] > sum {
			return "NO"
		}
		sum += b[i]
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d\n", len(c.arr)))
		for j, v := range c.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		exp := expected(c.arr)
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != exp {
			fmt.Printf("case %d failed: expected %s got %s (arr %v)\n", i+1, exp, got, c.arr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
