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

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest() string {
	n := rand.Intn(5) + 2 // 2..6
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 1; i <= n; i++ {
		l := rand.Intn(50) + 1
		r := l + rand.Intn(50)
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		fmt.Fprintf(&sb, "%d %d\n", i, p)
	}
	return sb.String()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	// build reference solution
	ref := "refA.bin"
	if out, err := exec.Command("go", "build", "-o", ref, "1528A.go").CombinedOutput(); err != nil {
		fmt.Println("failed to build reference:", string(out))
		return
	}
	defer os.Remove(ref)

	for i := 1; i <= 100; i++ {
		input := genTest()
		expect, err := runBinary("./"+ref, input)
		if err != nil {
			fmt.Println("reference run error on test", i, err)
			return
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Println("candidate run error on test", i, err)
			return
		}
		if expect != got {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, input, expect, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
