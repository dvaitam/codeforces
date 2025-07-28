package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solve(profits []int64, ratings [][]int) int64 {
	n := len(profits)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		ra := ratings[0][order[a]]
		rb := ratings[0][order[b]]
		if ra == rb {
			return order[a] < order[b]
		}
		return ra < rb
	})

	dp := make([]int64, n)
	for _, idx := range order {
		dp[idx] = profits[idx]
	}

	m := len(ratings)
	for i := 0; i < n; i++ {
		a := order[i]
		for j := i + 1; j < n; j++ {
			b := order[j]
			ok := true
			for city := 1; city < m; city++ {
				if ratings[city][a] >= ratings[city][b] {
					ok = false
					break
				}
			}
			if ok {
				if dp[a]+profits[b] > dp[b] {
					dp[b] = dp[a] + profits[b]
				}
			}
		}
	}

	var ans int64
	for _, v := range dp {
		if v > ans {
			ans = v
		}
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(0)
	for t := 1; t <= 100; t++ {
		m := rand.Intn(3) + 1
		n := rand.Intn(6) + 1
		profits := make([]int64, n)
		for i := range profits {
			profits[i] = rand.Int63n(20) + 1
		}
		ratings := make([][]int, m)
		for i := 0; i < m; i++ {
			ratings[i] = make([]int, n)
			perm := rand.Perm(n)
			for j := 0; j < n; j++ {
				ratings[i][j] = perm[j] + 1
			}
		}
		input := fmt.Sprintf("%d %d\n", m, n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(profits[i])
		}
		input += "\n"
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					input += " "
				}
				input += fmt.Sprint(ratings[i][j])
			}
			input += "\n"
		}
		expected := fmt.Sprint(solve(profits, ratings))
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		if output != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
