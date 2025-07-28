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
	n := rand.Intn(3) + 2 // 2..4
	m := n
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		used := make(map[int]bool)
		b := rand.Intn(n)
		used[b] = true
		c := rand.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", i, b, c)
	}
	return sb.String()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref := "refD.bin"
	if out, err := exec.Command("go", "build", "-o", ref, "1528D.go").CombinedOutput(); err != nil {
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
