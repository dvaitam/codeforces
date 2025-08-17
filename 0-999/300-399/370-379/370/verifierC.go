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

func buildRef() (string, error) {
	exe, err := os.CreateTemp("", "ref370C-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "370C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		x := rng.Intn(m) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)

		// Parse input to get counts of each color.
		in := strings.NewReader(input)
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
			os.Exit(1)
		}
		cnt := make([]int, m+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x]++
		}

		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expR := strings.NewReader(exp)
		var expGood int
		fmt.Fscan(expR, &expGood)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}

		gr := strings.NewReader(got)
		var gotGood int
		if _, err := fmt.Fscan(gr, &gotGood); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot read answer: %v\n", t+1, err)
			os.Exit(1)
		}
		leftCnt := make([]int, m+1)
		rightCnt := make([]int, m+1)
		good := 0
		for i := 0; i < n; i++ {
			var l, r int
			if _, err := fmt.Fscan(gr, &l, &r); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: not enough pairs\n", t+1)
				os.Exit(1)
			}
			if l < 1 || l > m || r < 1 || r > m {
				fmt.Fprintf(os.Stderr, "case %d: color out of range\n", t+1)
				os.Exit(1)
			}
			leftCnt[l]++
			rightCnt[r]++
			if l != r {
				good++
			}
		}
		var extra int
		if _, err := fmt.Fscan(gr, &extra); err == nil {
			fmt.Fprintf(os.Stderr, "case %d: trailing data in output\n", t+1)
			os.Exit(1)
		}

		if gotGood != expGood {
			fmt.Fprintf(os.Stderr, "case %d: reported %d good, expected %d\n", t+1, gotGood, expGood)
			os.Exit(1)
		}
		if good != gotGood {
			fmt.Fprintf(os.Stderr, "case %d: mismatch count is %d but reported %d\n", t+1, good, gotGood)
			os.Exit(1)
		}
		for c := 1; c <= m; c++ {
			if leftCnt[c] != cnt[c] || rightCnt[c] != cnt[c] {
				fmt.Fprintf(os.Stderr, "case %d: color %d count mismatch\n", t+1, c)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
