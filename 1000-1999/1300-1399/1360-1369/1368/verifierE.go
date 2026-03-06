package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func genTest(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		if m > 8 {
			m = 8
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		used := make(map[[2]int]bool)
		for j := 0; j < m; j++ {
			var x, y int
			for {
				x = rng.Intn(n-1) + 1
				y = rng.Intn(n-x) + x + 1
				if !used[[2]int{x, y}] {
					used[[2]int{x, y}] = true
					break
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
	}
	return sb.String()
}

type testCase struct {
	n, m  int
	edges [][2]int
}

func parseInput(input string) []testCase {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	next := func() int {
		sc.Scan()
		v, _ := strconv.Atoi(sc.Text())
		return v
	}
	t := next()
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		n := next()
		m := next()
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			edges[j] = [2]int{next(), next()}
		}
		cases[i] = testCase{n, m, edges}
	}
	return cases
}

func verify(tc testCase, output string) error {
	sc := bufio.NewScanner(strings.NewReader(output))
	sc.Split(bufio.ScanWords)
	next := func() (int, bool) {
		if !sc.Scan() {
			return 0, false
		}
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, false
		}
		return v, true
	}

	k, ok := next()
	if !ok {
		return fmt.Errorf("missing k")
	}
	if 7*k > 4*tc.n {
		return fmt.Errorf("too many closed: %d (limit %d*4/7)", k, tc.n)
	}

	closed := make(map[int]bool)
	for i := 0; i < k; i++ {
		v, ok := next()
		if !ok {
			return fmt.Errorf("expected %d closed spots, got %d", k, i)
		}
		if v < 1 || v > tc.n {
			return fmt.Errorf("spot %d out of range [1,%d]", v, tc.n)
		}
		if closed[v] {
			return fmt.Errorf("duplicate spot %d", v)
		}
		closed[v] = true
	}

	hasIn := make([]bool, tc.n+1)
	hasOut := make([]bool, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		if closed[u] || closed[v] {
			continue
		}
		hasOut[u] = true
		hasIn[v] = true
	}
	for v := 1; v <= tc.n; v++ {
		if !closed[v] && hasIn[v] && hasOut[v] {
			return fmt.Errorf("remaining vertex %d has both incoming and outgoing edges (path of length >=2 exists)", v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input := genTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\ninput:\n%s", i+1, err, stderr.String(), input)
			os.Exit(1)
		}

		cases := parseInput(input)
		sc := bufio.NewScanner(strings.NewReader(stdout.String()))
		for ci, tc := range cases {
			var lines []string
			// Read k line
			if !sc.Scan() {
				fmt.Fprintf(os.Stderr, "case %d subcase %d: missing output\ninput:\n%s", i+1, ci+1, input)
				os.Exit(1)
			}
			lines = append(lines, sc.Text())
			// Read closed spots line
			if sc.Scan() {
				lines = append(lines, sc.Text())
			}
			caseOutput := strings.Join(lines, " ")
			if err := verify(tc, caseOutput); err != nil {
				fmt.Fprintf(os.Stderr, "case %d subcase %d: %v\ninput:\n%sgot:\n%s\n", i+1, ci+1, err, input, stdout.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
