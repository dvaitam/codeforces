package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type card struct {
	a int
	b int
}

type testCase struct {
	cards []card
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveB(cards []card) string {
	n := len(cards)
	fronts := make(map[int]int)
	backs := make(map[int]int)
	for _, c := range cards {
		fronts[c.a]++
		if c.b != c.a {
			backs[c.b]++
		}
	}
	needed := (n + 1) / 2
	const INF = 1 << 60
	ans := INF
	for color, fcnt := range fronts {
		if fcnt >= needed {
			ans = 0
			break
		}
		need := needed - fcnt
		if bcnt, ok := backs[color]; ok && bcnt >= need {
			if need < ans {
				ans = need
			}
		}
	}
	for color, bcnt := range backs {
		if _, seen := fronts[color]; seen {
			continue
		}
		if bcnt >= needed && needed < ans {
			ans = needed
		}
	}
	if ans == INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{
		{cards: []card{{4, 4}, {4, 4}, {7, 7}}},
		{cards: []card{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}},
	}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		cards := make([]card, n)
		for i := 0; i < n; i++ {
			a := rng.Intn(10) + 1
			b := rng.Intn(10) + 1
			cards[i] = card{a: a, b: b}
		}
		tests = append(tests, testCase{cards: cards})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", len(t.cards)))
		for _, c := range t.cards {
			input.WriteString(fmt.Sprintf("%d %d\n", c.a, c.b))
		}
		expect := strings.TrimSpace(solveB(t.cards))
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
