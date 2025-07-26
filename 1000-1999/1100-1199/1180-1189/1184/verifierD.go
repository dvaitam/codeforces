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

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1184D1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(9) + 2   // 2..10
	m := n + rng.Intn(10)  // at least n
	k := rng.Intn(n-1) + 2 // 2..n
	tOps := rng.Intn(20) + 1
	if m > 250 {
		m = 250
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, m, tOps))
	l := n
	for i := 0; i < tOps; i++ {
		typ := 1
		if l > 1 && (l == m || rng.Intn(2) == 0) {
			typ = 0
		} else if l < m && rng.Intn(2) == 0 {
			typ = 1
		}
		if typ == 1 { // insert
			pos := rng.Intn(l+1) + 1
			sb.WriteString(fmt.Sprintf("1 %d\n", pos))
			if pos <= k {
				k++
			}
			l++
		} else { // cut
			pos := rng.Intn(l-1) + 1
			sb.WriteString(fmt.Sprintf("0 %d\n", pos))
			if k <= pos {
				l = pos
			} else {
				l = l - pos
				k = k - pos
			}
		}
		if l > m {
			l = m
		}
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", t+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", t+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
