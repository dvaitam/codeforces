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
	n := rand.Intn(20) + 1
	k := rand.Intn(10) + 1
	return fmt.Sprintf("%d %d\n", n, k)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref := "refF.bin"
	if out, err := exec.Command("go", "build", "-o", ref, "1528F.go").CombinedOutput(); err != nil {
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
			fmt.Printf("mismatch on test %d\ninput:%sexpected:%s\ngot:%s\n", i, input, expect, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
