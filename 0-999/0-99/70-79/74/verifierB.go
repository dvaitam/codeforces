package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expectedOutcome(n, m, k int, dirStr, st string) string {
	dir := 0
	if strings.Contains(dirStr, "tail") {
		dir = 1
	} else {
		dir = -1
	}
	S := make([]bool, n+1)
	S[m] = true
	cpos := k
	for i := 1; i <= len(st); i++ {
		ch := st[i-1]
		if ch == '0' {
			U := make([]bool, n+1)
			for p := 1; p <= n; p++ {
				if !S[p] {
					continue
				}
				for _, dp := range []int{-1, 0, 1} {
					q := p + dp
					if q < 1 || q > n {
						continue
					}
					if q == cpos {
						continue
					}
					U[q] = true
				}
			}
			newpos := cpos + dir
			newdir := dir
			if newpos == 1 || newpos == n {
				newdir = -dir
			}
			if newpos >= 1 && newpos <= n {
				U[newpos] = false
			}
			ok := false
			for p := 1; p <= n; p++ {
				if U[p] {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Sprintf("Controller %d", i)
			}
			S = U
			cpos = newpos
			dir = newdir
		} else {
			newpos := cpos + dir
			newdir := dir
			if newpos == 1 || newpos == n {
				newdir = -dir
			}
			if i == len(st) {
				return "Stowaway"
			}
			U := make([]bool, n+1)
			for p := 1; p <= n; p++ {
				if p != newpos {
					U[p] = true
				}
			}
			S = U
			cpos = newpos
			dir = newdir
		}
	}
	return "Stowaway"
}

func generateTest() (string, string) {
	n := rand.Intn(8) + 2 // 2..9
	m := rand.Intn(n) + 1
	k := rand.Intn(n) + 1
	for k == m {
		k = rand.Intn(n) + 1
	}
	dirStr := "to tail"
	if k > 1 && k < n {
		if rand.Intn(2) == 0 {
			dirStr = "to head"
		}
	} else if k == n {
		dirStr = "to head"
	}
	if k == 1 {
		dirStr = "to tail"
	}
	tlen := rand.Intn(9) + 1
	stBytes := make([]byte, tlen)
	for i := 0; i < tlen-1; i++ {
		if rand.Intn(2) == 0 {
			stBytes[i] = '0'
		} else {
			stBytes[i] = '1'
		}
	}
	stBytes[tlen-1] = '1'
	st := string(stBytes)
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d %d\n", n, m, k)
	fmt.Fprintln(&buf, dirStr)
	fmt.Fprintln(&buf, st)
	want := expectedOutcome(n, m, k, dirStr, st)
	return buf.String(), want
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest()
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(inp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
