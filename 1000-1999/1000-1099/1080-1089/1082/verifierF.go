package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(name string, args []string, input string) (string, error) {
	cmd := exec.Command(name, args...)
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

func randStr() string {
	l := rand.Intn(3) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('0' + rand.Intn(10))
	}
	return string(b)
}

func genTests() string {
	rand.Seed(1)
	const t = 100
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rand.Intn(6) + 1
		k := rand.Intn(3)
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%s %d\n", randStr(), rand.Intn(10))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	inp := genTests()
	candOut, err := runCmd(os.Args[1], nil, inp)
	if err != nil {
		fmt.Println("candidate error:", err)
		fmt.Print(candOut)
		os.Exit(1)
	}
	refOut, err := runCmd("go", []string{"run", "1082F.go"}, inp)
	if err != nil {
		fmt.Println("reference error:", err)
		fmt.Print(refOut)
		os.Exit(1)
	}
	if normalize(candOut) != normalize(refOut) {
		fmt.Println("outputs differ")
		fmt.Println("candidate:\n", candOut)
		fmt.Println("expected:\n", refOut)
		os.Exit(1)
	}
	fmt.Println("OK")
}
