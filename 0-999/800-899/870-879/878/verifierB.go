package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	k     int
	m     int
	arr   []int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	k := rng.Intn(4) + 2
	m := rng.Intn(4) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, m))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(3) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteByte('\n')
	return Test{n: n, k: k, m: m, arr: arr, input: sb.String()}
}

func solve(t Test) string {
	type pair struct{ val, cnt int }
	stack := []pair{}
	for i := 0; i < t.m; i++ {
		for _, v := range t.arr {
			if len(stack) > 0 && stack[len(stack)-1].val == v {
				stack[len(stack)-1].cnt++
				if stack[len(stack)-1].cnt == t.k {
					stack = stack[:len(stack)-1]
				}
			} else {
				stack = append(stack, pair{v, 1})
			}
		}
	}
	res := 0
	for _, p := range stack {
		res += p.cnt
	}
	return strconv.Itoa(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
