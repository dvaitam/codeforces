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

func solveC(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(r, &n, &k); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	totalStress := 0
	removed := make([]int, k+1)
	added := make([]int, k+1)
	for i := 1; i < n; i++ {
		if a[i] != a[i-1] {
			totalStress++
			removed[a[i]]++
			removed[a[i-1]]++
		}
	}
	for i := 0; i < n; {
		j := i
		for j+1 < n && a[j+1] == a[i] {
			j++
		}
		g := a[i]
		if i > 0 && j < n-1 {
			left := a[i-1]
			right := a[j+1]
			if left != right {
				added[g]++
			}
		}
		i = j + 1
	}
	best := 1
	minStress := totalStress - removed[1] + added[1]
	for x := 2; x <= k; x++ {
		stress := totalStress - removed[x] + added[x]
		if stress < minStress {
			minStress = stress
			best = x
		}
	}
	return fmt.Sprintf("%d\n", best)
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(20) + 2          // >=2
	k := rng.Intn(min(5, n-1)) + 2 // 2..min(5,n)
	if k > n {
		k = n
	}
	genres := make([]int, n)
	for i := 0; i < k; i++ {
		genres[i] = i + 1
	}
	for i := k; i < n; i++ {
		genres[i] = rng.Intn(k) + 1
	}
	rng.Shuffle(n, func(i, j int) { genres[i], genres[j] = genres[j], genres[i] })
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, g := range genres {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(g))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseC(rng)
	}
	for i, tc := range cases {
		expect := solveC(tc)
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
