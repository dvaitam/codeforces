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

func toStr(val int64, k int) string {
	b := make([]byte, k)
	for i := k - 1; i >= 0; i-- {
		b[i] = byte(val%26) + 'a'
		val /= 26
	}
	return string(b)
}

func genTest() []byte {
	k := rand.Intn(10) + 1
	max := int64(1)
	for i := 0; i < k; i++ {
		max *= 26
	}
	if max < 3 {
		max = 3
	}
	a := rand.Int63n(max - 2)
	diff := rand.Int63n(max - a - 1)
	if diff%2 == 1 {
		diff++
	}
	if a+diff >= max {
		a = max - diff - 1
	}
	b := a + diff
	s := toStr(a, k)
	t := toStr(b, k)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%s\n%s\n", k, s, t))
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "1144E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
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
