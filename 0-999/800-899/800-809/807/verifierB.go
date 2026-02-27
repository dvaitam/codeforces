package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest(rng *rand.Rand) []byte {
	p := rng.Intn(475) + 26  // 26..500
	x := rng.Intn(20000) + 1 // 1..20000
	y := rng.Intn(x) + 1     // 1..x
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", p, x, y))
	return []byte(sb.String())
}

func check(s, p int) bool {
	cur := (s / 50) % 475
	for i := 0; i < 25; i++ {
		cur = (cur*96 + 42) % 475
		if cur+26 == p {
			return true
		}
	}
	return false
}

func requiredSuccessHacks(x, s int) int {
	t := (s - x) / 50
	if t <= 0 {
		return 0
	}
	return (t + 1) / 2
}

func solve(p, x, y int) int {
	best := 1 << 30
	for s := y; s <= 100000; s++ {
		if (s-x)%50 != 0 {
			continue
		}
		if !check(s, p) {
			continue
		}
		need := requiredSuccessHacks(x, s)
		if need < best {
			best = need
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(1))
	tests := [][]byte{
		[]byte("120 530 82\n"),
		[]byte("470 1498 1309\n"),
	}
	for i := 0; i < 200; i++ {
		tests = append(tests, genTest(rng))
	}

	for i, input := range tests {
		var p, x, y int
		if _, err := fmt.Sscanf(string(input), "%d %d %d", &p, &x, &y); err != nil {
			fmt.Printf("internal verifier parse error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		want := fmt.Sprintf("%d\n", solve(p, x, y))
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
