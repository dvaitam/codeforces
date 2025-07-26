package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solve(n int, a []int) int {
	prefixOdd := make([]int, n+1)
	prefixEven := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefixOdd[i] = prefixOdd[i-1]
		prefixEven[i] = prefixEven[i-1]
		if i%2 == 0 {
			prefixEven[i] += a[i-1]
		} else {
			prefixOdd[i] += a[i-1]
		}
	}
	ans := 0
	for i := 1; i <= n; i++ {
		odd := prefixOdd[i-1] + (prefixEven[n] - prefixEven[i])
		even := prefixEven[i-1] + (prefixOdd[n] - prefixOdd[i])
		if odd == even {
			ans++
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(2)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(10) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(50)
			sb.WriteString(fmt.Sprintf("%d ", a[i]))
		}
		sb.WriteString("\n")
		expected := solve(n, a)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(string(out))
			return
		}
		var got int
		fmt.Fscan(strings.NewReader(out), &got)
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected %d got %d\n", t, sb.String(), expected, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
