package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(8) + 1
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	var sb strings.Builder
	sRunes := make([]rune, n)
	for i := 0; i < n; i++ {
		sRunes[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sRunes)
	sb.WriteString(fmt.Sprintf("%s\n", s))
	q := rng.Intn(3) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	answers := make([]int, q)
	for i := 0; i < q; i++ {
		alen := rng.Intn(4) + 1
		blen := rng.Intn(4) + 1
		a := make([]rune, alen)
		b := make([]rune, blen)
		for j := 0; j < alen; j++ {
			a[j] = letters[rng.Intn(len(letters))]
		}
		for j := 0; j < blen; j++ {
			b[j] = letters[rng.Intn(len(letters))]
		}
		as := string(a)
		bs := string(b)
		sb.WriteString(fmt.Sprintf("%s %s\n", as, bs))
		ans := -1
		for l := 0; l < len(s); l++ {
			for r := l; r < len(s); r++ {
				sub := s[l : r+1]
				if strings.Contains(sub, as) && strings.Contains(sub, bs) {
					length := r - l + 1
					if ans == -1 || length < ans {
						ans = length
					}
				}
			}
		}
		answers[i] = ans
	}
	return sb.String(), answers
}

func runCase(bin string, input string, expected []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	got := make([]int, 0, len(expected))
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		got = append(got, v)
	}
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(got))
	}
	for i, v := range expected {
		if got[i] != v {
			return fmt.Errorf("query %d expected %d got %d", i, v, got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
