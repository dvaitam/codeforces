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

func solveCase(n, m int, s string, inds []int, c string) string {
	letters := []byte(c)
	sort.Ints(inds)
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	bs := []byte(s)
	uniq := make([]int, 0, len(inds))
	for _, v := range inds {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	for i := 0; i < len(uniq); i++ {
		idx := uniq[i] - 1
		if i < len(letters) {
			bs[idx] = letters[i]
		}
	}
	return string(bs)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	letters := make([]byte, n)
	for i := 0; i < n; i++ {
		letters[i] = byte('a' + rng.Intn(5))
	}
	inds := make([]int, m)
	for i := 0; i < m; i++ {
		inds[i] = rng.Intn(n) + 1
	}
	cs := make([]byte, m)
	for i := 0; i < m; i++ {
		cs[i] = byte('a' + rng.Intn(5))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n%s\n", n, m, string(letters))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", inds[i])
	}
	sb.WriteByte('\n')
	sb.WriteString(string(cs))
	sb.WriteByte('\n')
	expected := solveCase(n, m, string(letters), inds, string(cs)) + "\n"
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
