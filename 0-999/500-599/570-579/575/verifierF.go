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
	ref := "refF_bin"
	cmd := exec.Command("go", "build", "-o", ref, "575F.go")
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
	n := rand.Intn(4) + 1
	x := rand.Intn(20)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for i := 0; i < n; i++ {
		l := rand.Intn(20)
		r := l + rand.Intn(10)
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
