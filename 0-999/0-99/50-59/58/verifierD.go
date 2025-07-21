package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(names []string, sep string) []string {
	groups := make(map[int][]string)
	lengths := []int{}
	for _, s := range names {
		l := len(s)
		if _, ok := groups[l]; !ok {
			lengths = append(lengths, l)
		}
		groups[l] = append(groups[l], s)
	}
	sort.Ints(lengths)
	minL := lengths[0]
	maxL := lengths[len(lengths)-1]
	K := minL + maxL
	for _, l := range lengths {
		sort.Strings(groups[l])
	}
	res := []string{}
	i, j := 0, len(lengths)-1
	for i <= j {
		li := lengths[i]
		lj := lengths[j]
		sum := li + lj
		if sum < K {
			i++
		} else if sum > K {
			j--
		} else {
			if i < j {
				A := groups[li]
				B := groups[lj]
				for idx := 0; idx < len(A) && idx < len(B); idx++ {
					a := A[idx]
					b := B[idx]
					s1 := a + sep + b
					s2 := b + sep + a
					if s1 < s2 {
						res = append(res, s1)
					} else {
						res = append(res, s2)
					}
				}
			} else {
				A := groups[li]
				for idx := 0; idx+1 < len(A); idx += 2 {
					a := A[idx]
					b := A[idx+1]
					s1 := a + sep + b
					s2 := b + sep + a
					if s1 < s2 {
						res = append(res, s1)
					} else {
						res = append(res, s2)
					}
				}
			}
			i++
			j--
		}
	}
	sort.Strings(res)
	return res
}

func runCase(bin string, names []string, sep string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(names)))
	for _, n := range names {
		sb.WriteString(n)
		sb.WriteByte('\n')
	}
	sb.WriteString(sep)
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := []string{}
	sc := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}
	exp := expected(names, sep)
	if len(gotLines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(gotLines))
	}
	for i := range exp {
		if gotLines[i] != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], gotLines[i])
		}
	}
	return nil
}

func randName(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		names := make([]string, n)
		for j := range names {
			names[j] = randName(rng)
		}
		sep := string('a' + rng.Intn(26))
		if err := runCase(bin, names, sep); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v sep:%s\n", i+1, err, names, sep)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
