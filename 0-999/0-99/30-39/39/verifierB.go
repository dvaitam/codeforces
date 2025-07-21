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

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pos := make(map[int][]int)
	for i := 1; i <= n; i++ {
		v := a[i]
		pos[v] = append(pos[v], i)
	}
	var years []int
	last := 0
	for t := 1; ; t++ {
		found := -1
		for _, idx := range pos[t] {
			if idx > last {
				found = idx
				break
			}
		}
		if found == -1 {
			break
		}
		years = append(years, found)
		last = found
	}
	k := len(years)
	if k == 0 {
		return "0\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for i, idx := range years {
		year := 2000 + idx
		if i+1 < k {
			sb.WriteString(fmt.Sprintf("%d ", year))
		} else {
			sb.WriteString(fmt.Sprintf("%d\n", year))
		}
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseB(rng)
	}
	for i, tc := range cases {
		expect := solveB(tc)
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
