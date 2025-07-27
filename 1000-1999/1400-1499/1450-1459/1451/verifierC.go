package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(n, k int, a, b string) bool {
	var ca, cb [26]int
	for i := 0; i < n; i++ {
		ca[a[i]-'a']++
		cb[b[i]-'a']++
	}
	carry := 0
	for i := 0; i < 26; i++ {
		available := ca[i] + carry
		if available < cb[i] {
			return false
		}
		diff := available - cb[i]
		if diff%k != 0 {
			return false
		}
		carry = diff
	}
	return carry == 0
}

func runCase(bin string, n, k int, a, b string) error {
	input := fmt.Sprintf("1\n%d %d\n%s\n%s\n", n, k, a, b)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := "No"
	if solve(n, k, a, b) {
		expect = "Yes"
	}
	if strings.ToLower(got) != strings.ToLower(expect) {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func randString(rng *rand.Rand, n int) string {
	bs := make([]byte, n)
	for i := range bs {
		bs[i] = byte('a' + rng.Intn(26))
	}
	return string(bs)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(5) + 1
		if k > n {
			k = n
		}
		a := randString(rng, n)
		b := randString(rng, n)
		if err := runCase(bin, n, k, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
