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

func computeF(s string, c string) int64 {
	pos := make([][]int, 26)
	for i, ch := range s {
		pos[ch-'a'] = append(pos[ch-'a'], i+1)
	}
	var f int64
	for i, ch := range c {
		lst := pos[ch-'a']
		if len(lst) == 0 {
			f += int64(len(c))
		} else {
			target := i + 1
			idx := sort.Search(len(lst), func(j int) bool { return lst[j] >= target })
			d := 1<<31 - 1
			if idx < len(lst) {
				if diff := lst[idx] - target; diff < d {
					if diff < 0 {
						d = -diff
					} else {
						d = diff
					}
				}
			}
			if idx > 0 {
				diff := lst[idx-1] - target
				if diff < 0 {
					diff = -diff
				}
				if diff < d {
					d = diff
				}
			}
			f += int64(d)
		}
	}
	return f
}

func runCase(bin string, n, k int, s string, cands []string) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n%s\n", n, k, s))
	for _, c := range cands {
		sb.WriteString(c)
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcde")
	for t := 0; t < 100; t++ {
		k := rng.Intn(5) + 1
		n := rng.Intn(4) + 1
		var sb strings.Builder
		for i := 0; i < k; i++ {
			sb.WriteRune(letters[rng.Intn(len(letters))])
		}
		s := sb.String()
		cands := make([]string, n)
		for i := 0; i < n; i++ {
			clen := rng.Intn(5) + 1
			sb.Reset()
			for j := 0; j < clen; j++ {
				sb.WriteRune(letters[rng.Intn(len(letters))])
			}
			cands[i] = sb.String()
		}
		expectedVals := make([]string, n)
		for i, c := range cands {
			expectedVals[i] = fmt.Sprintf("%d", computeF(s, c))
		}
		expected := strings.Join(expectedVals, "\n")
		got, err := runCase(bin, n, k, s, cands)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\n", t+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
