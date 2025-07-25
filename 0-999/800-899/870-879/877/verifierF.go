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
	n := rand.Intn(40) + 1
	k := rand.Intn(21) - 10
	tvals := make([]int, n)
	for i := range tvals {
		if rand.Intn(2) == 0 {
			tvals[i] = 1
		} else {
			tvals[i] = 2
		}
	}
	avals := make([]int, n)
	for i := range avals {
		avals[i] = rand.Intn(21)
	}
	q := rand.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range tvals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range avals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n-l+1) + l
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "877F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := generateTest()
		want, err := run(ref, input)
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
