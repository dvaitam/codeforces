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

type TestCase struct {
	input string
}

func genTestCase() TestCase {
	n := rand.Intn(7) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Int63n(1e6)+1))
	}
	sb.WriteByte('\n')
	return TestCase{input: sb.String()}
}

func runBin(path string, in string) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	user := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "1667A.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	for i := 0; i < 100; i++ {
		tc := genTestCase()
		exp, err := runBin(ref, tc.input)
		if err != nil {
			fmt.Println("reference failed on test", i+1, err)
			os.Exit(1)
		}
		out, err := runBin(user, tc.input)
		if err != nil {
			fmt.Println("runtime error on test", i+1)
			os.Exit(1)
		}
		if strings.TrimSpace(string(out)) != strings.TrimSpace(string(exp)) {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected:%s\nfound:%s\n", i+1, tc.input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
