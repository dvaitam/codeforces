package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func randArr(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(100) + 1
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	refBin := "./_refD"
	if err := exec.Command("go", "build", "-o", refBin, "392D.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(refBin)

	rand.Seed(1)
	for tc := 1; tc <= 100; tc++ {
		n := rand.Intn(50) + 1
		a := randArr(n)
		b := randArr(n)
		c := randArr(n)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution error:", err)
			return
		}
		got, err := runCmd(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", tc, err)
			return
		}
		if got != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected:%s\ngot:%s\n", tc, input, exp, got)
			return
		}
	}
	fmt.Println("All tests passed!")
}
