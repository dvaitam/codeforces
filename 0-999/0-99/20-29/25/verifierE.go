package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func prefixFunc(p []byte) []int {
	n := len(p)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && p[i] != p[j] {
			j = pi[j-1]
		}
		if p[i] == p[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func contains(a, b string) bool {
	la, lb := len(a), len(b)
	if lb > la {
		return false
	}
	pat := []byte(b)
	pi := prefixFunc(pat)
	j := 0
	for i := 0; i < la; i++ {
		for j > 0 && a[i] != b[j] {
			j = pi[j-1]
		}
		if a[i] == b[j] {
			j++
			if j == lb {
				return true
			}
		}
	}
	return false
}

func overlap(a, b string) int {
	la, lb := len(a), len(b)
	start := 0
	if la > lb {
		start = la - lb
	}
	pat := []byte(b)
	pi := prefixFunc(pat)
	j := 0
	for i := start; i < la; i++ {
		for j > 0 && a[i] != b[j] {
			j = pi[j-1]
		}
		if a[i] == b[j] {
			j++
			if j == lb {
				break
			}
		}
	}
	return j
}

func merge(a, b string) string {
	k := overlap(a, b)
	return a + b[k:]
}

func expected(strs []string) int {
	keep := make([]string, 0, 3)
	for i, s := range strs {
		skip := false
		for j, t := range strs {
			if i == j {
				continue
			}
			if contains(t, s) {
				if len(t) > len(s) || (len(t) == len(s) && j < i) {
					skip = true
					break
				}
			}
		}
		if !skip {
			keep = append(keep, s)
		}
	}
	
	var result int
	n := len(keep)
	if n == 0 {
		result = len(strs[0])
	} else if n == 1 {
		result = len(keep[0])
	} else if n == 2 {
		a, b := keep[0], keep[1]
		l1 := len(a) + len(b) - overlap(a, b)
		l2 := len(a) + len(b) - overlap(b, a)
		if l1 < l2 {
			result = l1
		} else {
			result = l2
		}
	} else {
		perm := []int{0, 1, 2}
		result = -1
		var gen func(int)
		gen = func(idx int) {
			if idx == 3 {
				a, b, c := keep[perm[0]], keep[perm[1]], keep[perm[2]]
				t := merge(a, b)
				t = merge(t, c)
				l := len(t)
				if result < 0 || l < result {
					result = l
				}
				return
			}
			for i := idx; i < 3; i++ {
				perm[idx], perm[i] = perm[i], perm[idx]
				gen(idx + 1)
				perm[idx], perm[i] = perm[i], perm[idx]
			}
		}
		gen(0)
	}
	
	return result
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateRandomString(minLen, maxLen int, rng *rand.Rand) string {
	length := rng.Intn(maxLen-minLen+1) + minLen
	chars := make([]byte, length)
	for i := 0; i < length; i++ {
		chars[i] = byte('a' + rng.Intn(26))
	}
	return string(chars)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		strs := make([]string, 3)
		for j := 0; j < 3; j++ {
			strs[j] = generateRandomString(1, 20, rng)
		}
		
		var input strings.Builder
		for j := 0; j < 3; j++ {
			input.WriteString(strs[j])
			input.WriteString("\n")
		}
		
		expectedOut := expected(strs)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		gotInt, parseErr := strconv.Atoi(got)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input.String())
			os.Exit(1)
		}
		
		if gotInt != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expectedOut, gotInt, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}