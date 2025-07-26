package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveE(D []int) string {
	n := len(D) - 1
	sort.Slice(D[1:], func(i, j int) bool { return D[1+i] > D[1+j] })
	S1 := make([]int64, n+2)
	S2 := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		S1[i] = S1[i-1] + int64(D[i])
	}
	for i := n; i >= 1; i-- {
		S2[i] = S2[i+1] + int64(D[i])
	}
	T := make([]int, n+2)
	j := n + 1
	for i := 0; i <= n; i++ {
		for j > 1 && D[j-1] <= i {
			j--
		}
		T[i] = j
	}
	get := func(i, k int) int64 {
		t := T[k]
		if i > t {
			t = i
		}
		return int64(t-i)*int64(k) + S2[t]
	}
	P1 := make([]int64, n+2)
	P2 := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		pi1 := S1[i] - int64(i)*(int64(i)-1) - get(i+1, i)
		if pi1 > int64(i) {
			P1[i] = int64(n) + 1
		} else {
			P1[i] = pi1
		}
		P2[i] = int64(i+1)*int64(i) + get(i+1, i+1) - S1[i]
	}
	P1[0] = 0
	P2[n+1] = int64(n) + 1
	for i := 1; i <= n; i++ {
		if P1[i-1] > P1[i] {
			P1[i] = P1[i-1]
		}
	}
	for i := n; i >= 1; i-- {
		if P2[i+1] < P2[i] {
			P2[i] = P2[i+1]
		}
	}
	total := S1[n]
	var res []int
	j = n + 1
	for i := int(total & 1); i <= n; i += 2 {
		for j > 1 && i >= D[j-1] {
			j--
		}
		cond1 := P1[j-1] <= int64(i)
		cond2 := P2[j] >= int64(i)
		cond3 := S1[j-1]+int64(i) <= int64(j)*(int64(j)-1)+get(j, j)
		if cond1 && cond2 && cond3 {
			res = append(res, i)
		}
	}
	if len(res) == 0 {
		return "-1"
	}
	var sb strings.Builder
	for idx, v := range res {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(50) + 1
		D := make([]int, n+1)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			D[i] = rand.Intn(n + 1)
			input.WriteString(fmt.Sprintf("%d ", D[i]))
		}
		input.WriteByte('\n')
		expect := solveE(D)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Println("input:\n", input.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("wrong answer on test %d\n", t+1)
			fmt.Println("input:\n", input.String())
			fmt.Printf("expected: %s\n got: %s\n", expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
