package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isVowel(c byte) bool {
	switch c {
	case 'I', 'E', 'A', 'O', 'U', 'Y':
		return true
	}
	return false
}

func solvePrettiness(s string) string {
	n := len(s)
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1]
		if isVowel(s[i-1]) {
			pref[i]++
		}
	}
	prev := float64(pref[n])
	ret := prev
	for l := 2; l <= n; l++ {
		delta := float64(pref[n-l+1] - pref[l-1])
		ans := prev + delta
		ret += ans / float64(l)
		prev = ans
	}
	return fmt.Sprintf("%.8f", ret)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solvePrettiness(input)
	gotVal, err1 := strconv.ParseFloat(got, 64)
	expVal, err2 := strconv.ParseFloat(expected, 64)
	if err1 != nil {
		return fmt.Errorf("cannot parse candidate output %q: %v", got, err1)
	}
	if err2 != nil {
		return fmt.Errorf("cannot parse expected output %q: %v", expected, err2)
	}
	if math.Abs(gotVal-expVal) > 1e-6*math.Max(1.0, math.Abs(expVal)) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
