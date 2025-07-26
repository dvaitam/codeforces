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

type emp struct{ l, r int }

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1251D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runProg(prog string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genCase() []byte {
	t := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for ; t > 0; t-- {
		n := rand.Intn(5) + 1
		s := rand.Int63n(int64(n*100)) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
		for i := 0; i < n; i++ {
			l := rand.Intn(100) + 1
			r := l + rand.Intn(100-l+1)
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genCase()
		want, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
