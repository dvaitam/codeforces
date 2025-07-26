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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func generateTest() []byte {
	t := rand.Intn(3) + 1 // up to 3 test cases
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rand.Intn(49) + 2 // 2..50
		arr := make([]int, n)
		arr[0] = rand.Intn(100) + 1
		for j := 1; j < n; j++ {
			arr[j] = arr[j-1] + rand.Intn(3)
			if arr[j] > 100 {
				arr[j] = 100
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	refBin := "./refA.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1783A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := generateTest()
		want, err := run(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
