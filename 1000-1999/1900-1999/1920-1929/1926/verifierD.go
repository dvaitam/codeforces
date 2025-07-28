package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mask = (1 << 31) - 1

func expected(arr []int) string {
	freq := make(map[int]int)
	pairs := 0
	for _, x := range arr {
		c := mask ^ x
		if freq[c] > 0 {
			freq[c]--
			pairs++
		} else {
			freq[x]++
		}
	}
	return fmt.Sprint(len(arr) - pairs)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	var tests [][]int
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(mask)
		}
		tests = append(tests, arr)
	}
	for idx, arr := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(arr)
		if got != exp {
			fmt.Printf("test %d failed: expected=%s got=%s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
