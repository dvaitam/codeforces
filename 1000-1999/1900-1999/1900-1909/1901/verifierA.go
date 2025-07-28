package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type result struct {
	out string
	err error
}

func buildRef() (string, error) {
	refSrc := filepath.Join(filepath.Dir(os.Args[0]), "1901A.go")
	bin := filepath.Join(os.TempDir(), fmt.Sprintf("refA_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", bin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return bin, nil
}

func runBinary(path, input string) result {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, stderr.String())
	}
	return result{strings.TrimSpace(out.String()), err}
}

func genTest() string {
	n := rand.Intn(50) + 1       // 1..50
	x := rand.Intn(99) + 2       // 2..100
	vals := rand.Perm(x - 1)[:n] // choose n distinct from 0..x-2
	for i := range vals {
		vals[i]++ // shift to 1..x-1
	}
	sort.Ints(vals)
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "1\n%d %d\n", n, x)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(1)

	refBin, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		tc := genTest()
		exp := runBinary(refBin, tc)
		if exp.err != nil {
			fmt.Fprintf(os.Stderr, "reference run failed on test %d: %v\n", i+1, exp.err)
			os.Exit(1)
		}
		got := runBinary(binary, tc)
		if got.err != nil {
			fmt.Fprintf(os.Stderr, "binary failed on test %d: %v\n", i+1, got.err)
			os.Exit(1)
		}
		if exp.out != got.out {
			fmt.Printf("mismatch on test %d\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc, exp.out, got.out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
