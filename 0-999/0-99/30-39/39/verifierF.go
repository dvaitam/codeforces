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

func solveF(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int64
	var m, k int
	if _, err := fmt.Fscan(r, &n, &m, &k); err != nil {
		return ""
	}
	d := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &d[i])
	}
	p := make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(r, &p[i])
	}
	cnt := make([]int, m)
	minc := k + 1
	for i := 0; i < m; i++ {
		c := 0
		for j := 0; j < k; j++ {
			if p[j]%d[i] == 0 {
				c++
			}
		}
		cnt[i] = c
		if c < minc {
			minc = c
		}
	}
	var res []int
	for i := 0; i < m; i++ {
		if cnt[i] == minc {
			res = append(res, i+1)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCaseF(rng *rand.Rand) string {
	n := int64(rng.Intn(100) + 1)
	m := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(10) + 1))
	}
	sb.WriteByte('\n')
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseF(rng)
	}
	for i, tc := range cases {
		expect := solveF(tc)
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
