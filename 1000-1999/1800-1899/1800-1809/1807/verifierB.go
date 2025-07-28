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
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(20) + 1
		}
		cases[i] = Case{arr: a}
	}
	return cases
}

func expected(a []int) string {
	even, odd := 0, 0
	for _, x := range a {
		if x%2 == 0 {
			even += x
		} else {
			odd += x
		}
	}
	if even > odd {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
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
			fmt.Printf("case %d failed: expected %s got %s (array %v)\n", i+1, exp, got, c.arr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
