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

func expected(s string) [][2]int {
	n := len(s)
	z := make([]int, n)
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			if r-i+1 < z[i-l] {
				z[i] = r - i + 1
			} else {
				z[i] = z[i-l]
			}
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	cnt := make([]int, n+2)
	for i := 1; i < n; i++ {
		cnt[1]++
		cnt[z[i]+1]--
	}
	for i := 1; i <= n; i++ {
		cnt[i] += cnt[i-1]
		cnt[i]++
	}
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	var borders []int
	k := n
	for k > 0 {
		borders = append(borders, k)
		if k == n {
			k = pi[n-1]
		} else {
			k = pi[k-1]
		}
	}
	for i, j := 0, len(borders)-1; i < j; i, j = i+1, j-1 {
		borders[i], borders[j] = borders[j], borders[i]
	}
	res := make([][2]int, len(borders))
	for i, b := range borders {
		res[i] = [2]int{b, cnt[b]}
	}
	return res
}

func runCase(exe string, s string) error {
	input := s + "\n"
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	pairs := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &pairs[i][0], &pairs[i][1]); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
	}
	expectedRes := expected(s)
	if len(expectedRes) != k {
		return fmt.Errorf("expected %d pairs got %d", len(expectedRes), k)
	}
	for i := 0; i < k; i++ {
		if pairs[i][0] != expectedRes[i][0] || pairs[i][1] != expectedRes[i][1] {
			return fmt.Errorf("pair %d expected %d %d got %d %d", i+1, expectedRes[i][0], expectedRes[i][1], pairs[i][0], pairs[i][1])
		}
	}
	return nil
}

func randomString(rng *rand.Rand, length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		length := rng.Intn(20) + 1
		s := randomString(rng, length)
		if err := runCase(exe, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
