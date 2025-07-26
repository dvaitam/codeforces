package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), "1097D_ref")
	cmd := exec.Command("go", "build", "-o", out, "1097D.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var primes = []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

func randomN() int64 {
	cnt := rand.Intn(4) + 1
	n := int64(1)
	for i := 0; i < cnt; i++ {
		n *= primes[rand.Intn(len(primes))]
	}
	return n
}

func genTest() string {
	n := randomN()
	k := rand.Intn(10) + 1
	return fmt.Sprintf("%d %d\n", n, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		test := genTest()
		exp, err := runBin(ref, test)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference run failed:", err)
			os.Exit(1)
		}
		got, err := runBin(cand, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate run failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:%s got:%s\n", t+1, test, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
