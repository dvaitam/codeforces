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

func genInput() string {
	t := 100
	opts := []int{8, 12}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := opts[rand.Intn(len(opts))]
		m := opts[rand.Intn(len(opts))]
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		line := strings.Repeat(".", m)
		for r := 0; r < n; r++ {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runBin(path, in string) ([]byte, error) {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	user := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "1667F.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	input := genInput()
	exp, err := runBin(ref, input)
	if err != nil {
		fmt.Println("reference failed:", err)
		os.Exit(1)
	}
	out, err := runBin(user, input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(string(out)) != strings.TrimSpace(string(exp)) {
		fmt.Printf("wrong answer\ninput:\n%s\nexpected:\n%s\nfound:\n%s\n", input, exp, out)
		os.Exit(1)
	}
	fmt.Println("ok")
}
