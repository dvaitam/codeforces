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

func buildRef() string {
	ref := "refA_bin"
	cmd := exec.Command("go", "build", "-o", ref, "575A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("failed to build reference: %v\n%s", err, string(out)))
	}
	return ref
}

func run(bin, input string) (string, error) {
	c := exec.Command(bin)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	return out.String(), err
}

func genCase() string {
	K := rand.Int63n(20)
	P := rand.Int63n(1000) + 1
	N := rand.Intn(3) + 1
	base := make([]int64, N)
	for i := 0; i < N; i++ {
		base[i] = rand.Int63n(P) + 1
	}
	M := rand.Intn(3)
	lines := fmt.Sprintf("%d %d\n%d\n", K, P, N)
	for i := 0; i < N; i++ {
		lines += fmt.Sprintf("%d ", base[i])
	}
	lines = strings.TrimSpace(lines) + "\n"
	lines += fmt.Sprintf("%d\n", M)
	for i := 0; i < M; i++ {
		j := rand.Int63n(K+int64(N)) + int64(N)
		v := rand.Int63n(P) + 1
		lines += fmt.Sprintf("%d %d\n", j, v)
	}
	return lines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	ref := buildRef()
	defer os.Remove(ref)
	for i := 0; i < 100; i++ {
		input := genCase()
		exp, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("binary failed on case %d: %v\n", i, err)
			return
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("mismatch on case %d\ninput:\n%s\nexpected:%s\nactual:%s\n", i, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
