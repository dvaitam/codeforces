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

type Test struct {
	n int
	a []int
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc Test) string {
	n := tc.n
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = tc.a[i-1]
	}
	leftLess := make([]int, n+2)
	rightLess := make([]int, n+2)
	leftGreater := make([]int, n+2)
	rightGreater := make([]int, n+2)
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			leftLess[i] = stack[len(stack)-1]
		} else {
			leftLess[i] = 0
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			rightLess[i] = stack[len(stack)-1]
		} else {
			rightLess[i] = n + 1
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			leftGreater[i] = stack[len(stack)-1]
		} else {
			leftGreater[i] = 0
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			rightGreater[i] = stack[len(stack)-1]
		} else {
			rightGreater[i] = n + 1
		}
		stack = append(stack, i)
	}
	res := make([]int, n+1)
	for i := 1; i <= n; i++ {
		lenMin := rightLess[i] - leftLess[i] - 1
		distMin := lenMin / 2
		lenMax := rightGreater[i] - leftGreater[i] - 1
		distMax := (lenMax+1)/2 - 1
		if distMin > distMax {
			res[i] = distMin
		} else {
			res[i] = distMax
		}
	}
	var out strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", res[i]))
	}
	return out.String()
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(50)
	}
	return Test{n: n, a: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
