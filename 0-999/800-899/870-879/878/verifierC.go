package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n       int
	k       int
	players [][]int
	input   string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	k := rng.Intn(4) + 1
	players := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		players[i] = make([]int, k)
		for j := 0; j < k; j++ {
			players[i][j] = rng.Intn(10)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(players[i][j]))
		}
		sb.WriteByte('\n')
	}
	return Test{n: n, k: k, players: players, input: sb.String()}
}

func dominates(a, b []int) bool {
	for i := range a {
		if a[i] <= b[i] {
			return false
		}
	}
	return true
}

func solve(t Test) string {
	nondom := make([][]int, 0)
	var sb strings.Builder
	for i := 0; i < t.n; i++ {
		p := t.players[i]
		dominated := false
		for _, q := range nondom {
			if dominates(q, p) {
				dominated = true
				break
			}
		}
		if !dominated {
			newSet := nondom[:0]
			for _, q := range nondom {
				if !dominates(p, q) {
					newSet = append(newSet, q)
				}
			}
			nondom = append(newSet, p)
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(nondom)))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
