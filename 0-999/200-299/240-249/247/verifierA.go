package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveA implements the logic of 247A using the code from 247A.go
func solveA(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	var sizes []int
	currLen := 0
	negCnt := 0
	for _, v := range a {
		if v < 0 {
			if negCnt == 2 {
				sizes = append(sizes, currLen)
				currLen = 1
				negCnt = 1
			} else {
				currLen++
				negCnt++
			}
		} else {
			currLen++
		}
	}
	if currLen > 0 {
		sizes = append(sizes, currLen)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(sizes)))
	for i, sz := range sizes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(sz))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(20) + 1 // 1..20
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(41) - 20 // -20..20
		sb.WriteString(fmt.Sprint(val))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseA(rng)
	}
	for i, tc := range cases {
		expect := solveA(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
