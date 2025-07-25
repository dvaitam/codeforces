package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(cities, towers []int) int {
	sort.Ints(towers)
	m := len(towers)
	result := 0
	for _, x := range cities {
		idx := sort.SearchInts(towers, x)
		dist := int(1<<31 - 1)
		if idx < m {
			d := towers[idx] - x
			if d < dist {
				dist = d
			}
		}
		if idx > 0 {
			d := x - towers[idx-1]
			if d < dist {
				dist = d
			}
		}
		if dist > result {
			result = dist
		}
	}
	return result
}

func buildCase(cities, towers []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(cities), len(towers)))
	for i, v := range cities {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range towers {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	expect := strconv.Itoa(solve(cities, towers))
	return testCase{input: sb.String(), expected: expect}
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	cities := make([]int, n)
	towers := make([]int, m)
	for i := range cities {
		cities[i] = rng.Intn(100)
	}
	for i := range towers {
		towers[i] = rng.Intn(100)
	}
	return buildCase(cities, towers)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int{1}, []int{1}))
	cases = append(cases, buildCase([]int{0, 10}, []int{5}))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
