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
	tmp := "refA_bin"
	cmd := exec.Command("go", "build", "-o", tmp, "1765A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp, nil
}

func runProg(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func genTest(seed int64) string {
	r := rand.New(rand.NewSource(seed))
	n := r.Intn(3) + 1
	m := r.Intn(3) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if r.Intn(2) == 1 {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest(int64(i))
		expectOut, err := runProg("./"+ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		gotOut, err := runProg(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate crashed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if normalize(expectOut) != normalize(gotOut) {
			fmt.Printf("mismatch on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i, input, expectOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
