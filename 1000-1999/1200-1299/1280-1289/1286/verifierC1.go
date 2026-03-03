package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n int
	s string
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Process.Kill()

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	fmt.Fprintf(writer, "%d\n", tc.n)
	writer.Flush()

	totalSubstrings := 0
	queriesCount := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("unexpected EOF")
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "?") {
			queriesCount++
			if queriesCount > 3 {
				return fmt.Errorf("too many queries")
			}
			var l, r int
			_, err := fmt.Sscanf(line, "? %d %d", &l, &r)
			if err != nil {
				return fmt.Errorf("invalid query format: %s", line)
			}
			if l < 1 || r > tc.n || l > r {
				return fmt.Errorf("invalid query range: %d %d", l, r)
			}

			num := (r - l + 1) * (r - l + 2) / 2
			totalSubstrings += num
			if totalSubstrings > (tc.n+1)*(tc.n+1) {
				return fmt.Errorf("too many substrings: %d", totalSubstrings)
			}

			substrings := []string{}
			for i := l - 1; i < r; i++ {
				for j := i + 1; j <= r; j++ {
					sub := tc.s[i:j]
					chars := strings.Split(sub, "")
					sort.Strings(chars)
					substrings = append(substrings, strings.Join(chars, ""))
				}
			}
			rand.Shuffle(len(substrings), func(i, j int) {
				substrings[i], substrings[j] = substrings[j], substrings[i]
			})
			for _, sub := range substrings {
				fmt.Fprintln(writer, sub)
			}
			writer.Flush()
		} else if strings.HasPrefix(line, "!") {
			got := strings.TrimSpace(line[1:])
			if got != tc.s {
				return fmt.Errorf("wrong answer: expected %s, got %s", tc.s, got)
			}
			break
		} else {
			return fmt.Errorf("unexpected output: %s", line)
		}
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timeout")
	}
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{n: 1, s: "a"},
		{n: 2, s: "ab"},
		{n: 3, s: "abc"},
	}
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 20; i++ {
		n := rng.Intn(10) + 1
		if i > 15 {
			n = rng.Intn(90) + 10
		}
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
		cases = append(cases, testCase{n: n, s: sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (n=%d, s=%s): %v\n", i+1, tc.n, tc.s, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
