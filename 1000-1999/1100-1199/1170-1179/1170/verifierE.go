package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), "1170E_ref")
	cmd := exec.Command("go", "build", "-o", out, "1170E.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func randSubset(n, k int) []int {
	perm := rand.Perm(n)
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[i] = perm[i] + 1
	}
	sort.Ints(res)
	return res
}

func genTest() string {
	n := rand.Intn(3) + 1
	m := rand.Intn(10) + n
	widths := make([]int, n)
	rem := m
	for i := 0; i < n; i++ {
		maxW := rem - (n - i - 1)
		w := rand.Intn(maxW) + 1
		widths[i] = w
		rem -= w
	}
	q := rand.Intn(5) + 1
	lines := []string{fmt.Sprintf("%d %d", n, m)}
	wstr := make([]string, n)
	for i, v := range widths {
		wstr[i] = fmt.Sprintf("%d", v)
	}
	lines = append(lines, strings.Join(wstr, " "))
	lines = append(lines, fmt.Sprintf("%d", q))
	for i := 0; i < q; i++ {
		c := rand.Intn(m) + 1
		subset := randSubset(m, c)
		s := make([]string, c)
		for j, v := range subset {
			s[j] = fmt.Sprintf("%d", v)
		}
		lines = append(lines, fmt.Sprintf("%d %s", c, strings.Join(s, " ")))
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
