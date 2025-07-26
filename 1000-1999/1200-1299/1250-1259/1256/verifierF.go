package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(s, t string) string {
	n := len(s)
	cntS := [26]int{}
	cntT := [26]int{}
	for i := 0; i < n; i++ {
		cntS[s[i]-'a']++
		cntT[t[i]-'a']++
	}
	equal := true
	dup := false
	for i := 0; i < 26; i++ {
		if cntS[i] != cntT[i] {
			equal = false
			break
		}
		if cntS[i] >= 2 {
			dup = true
		}
	}
	if !equal {
		return "NO"
	}
	if dup {
		return "YES"
	}
	pos := make([]int, 26)
	for i := 0; i < n; i++ {
		pos[t[i]-'a'] = i
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = pos[s[i]-'a']
	}
	invParity := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				invParity ^= 1
			}
		}
	}
	if invParity == 0 {
		return "YES"
	}
	return "NO"
}

func generate() (string, string) {
	const T = 100
	rand.Seed(6)
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	letters := []rune("abcdef")
	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 1
		sb1 := make([]rune, n)
		sb2 := make([]rune, n)
		for j := 0; j < n; j++ {
			sb1[j] = letters[rand.Intn(len(letters))]
			sb2[j] = letters[rand.Intn(len(letters))]
		}
		s := string(sb1)
		t := string(sb2)
		fmt.Fprintf(&in, "%d\n%s\n%s\n", n, s, t)
		fmt.Fprintln(&out, solveCase(s, t))
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Fprintln(os.Stderr, "wrong answer")
		fmt.Fprintln(os.Stderr, "expected:\n"+exp)
		fmt.Fprintln(os.Stderr, "got:\n"+out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
