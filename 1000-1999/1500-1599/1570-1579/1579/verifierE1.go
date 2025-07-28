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

func solveE1(p []int) string {
	n := len(p)
	dq := make([]int, 2*n+5)
	l := n
	r := n
	dq[l] = p[0]
	for i := 1; i < n; i++ {
		if p[i] < dq[l] {
			l--
			dq[l] = p[i]
		} else {
			r++
			dq[r] = p[i]
		}
	}
	res := dq[l : r+1]
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCaseE1(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	arr := rand.Perm(n)
	for i := range arr {
		arr[i]++
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveE1(arr)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseE1(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
