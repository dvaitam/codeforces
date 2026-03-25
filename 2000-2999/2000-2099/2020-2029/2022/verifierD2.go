package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n     int
	roles []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	tests := generateTests()

	// Compute expected answers directly
	expect := make([]int, len(tests))
	for i, tc := range tests {
		expect[i] = impostorIndex(tc.roles)
	}

	// Run candidate with an interaction simulator
	candAnswers, err := runInteractive(candidate, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != candAnswers[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, expect[i], candAnswers[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// runInteractive simulates the interactive judge for the candidate binary.
// The candidate reads t, then for each test reads n, then issues "? u v" queries
// and finally "! ans".
func runInteractive(binPath string, tests []testCase) ([]int, error) {
	cmd := exec.Command(binPath)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %v", err)
	}

	answers := make([]int, len(tests))
	scanner := newLineScanner(stdoutPipe)

	// Write number of test cases
	fmt.Fprintf(stdinPipe, "%d\n", len(tests))

	for ti, tc := range tests {
		// Write n
		fmt.Fprintf(stdinPipe, "%d\n", tc.n)

		// Now interact: read lines from candidate
		for {
			line, err := scanner.ReadLine()
			if err != nil {
				stdinPipe.Close()
				cmd.Wait()
				return nil, fmt.Errorf("test %d: read error: %v", ti+1, err)
			}
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if line[0] == '!' {
				// Final answer
				parts := strings.Fields(line)
				if len(parts) < 2 {
					stdinPipe.Close()
					cmd.Wait()
					return nil, fmt.Errorf("test %d: bad answer line: %s", ti+1, line)
				}
				val, err := strconv.Atoi(parts[1])
				if err != nil {
					stdinPipe.Close()
					cmd.Wait()
					return nil, fmt.Errorf("test %d: bad answer value: %s", ti+1, parts[1])
				}
				answers[ti] = val
				break
			} else if line[0] == '?' {
				// Query: ? u v
				// u asks about v: "are you and v the same type?"
				// Crewmate answers truthfully, impostor lies
				parts := strings.Fields(line)
				if len(parts) < 3 {
					stdinPipe.Close()
					cmd.Wait()
					return nil, fmt.Errorf("test %d: bad query: %s", ti+1, line)
				}
				u, _ := strconv.Atoi(parts[1])
				v, _ := strconv.Atoi(parts[2])

				uRole := tc.roles[u-1] // 0 or 1 for crewmate, -1 for impostor
				vRole := tc.roles[v-1]

				// Compute truth: are u and v the same type?
				// For the impostor, we assign them an effective color of 0
				uColor := uRole
				if uColor == -1 {
					uColor = 0
				}
				vColor := vRole
				if vColor == -1 {
					vColor = 0
				}

				// Truth: 0 if same color, 1 if different
				truth := 0
				if uColor != vColor {
					truth = 1
				}

				response := truth
				if uRole == -1 {
					// Impostor lies
					response = 1 - truth
				}

				fmt.Fprintf(stdinPipe, "%d\n", response)
			} else {
				stdinPipe.Close()
				cmd.Wait()
				return nil, fmt.Errorf("test %d: unexpected output: %s", ti+1, line)
			}
		}
	}

	stdinPipe.Close()
	cmd.Wait()
	return answers, nil
}

type lineScanner struct {
	buf  []byte
	pos  int
	r    interface{ Read([]byte) (int, error) }
	data []byte
}

func newLineScanner(r interface{ Read([]byte) (int, error) }) *lineScanner {
	return &lineScanner{r: r, buf: make([]byte, 4096)}
}

func (s *lineScanner) ReadLine() (string, error) {
	for {
		for i := s.pos; i < len(s.data); i++ {
			if s.data[i] == '\n' {
				line := string(s.data[s.pos:i])
				s.pos = i + 1
				return line, nil
			}
		}
		if s.pos > 0 {
			s.data = append(s.data[:0], s.data[s.pos:]...)
			s.pos = 0
		}
		n, err := s.r.Read(s.buf)
		if n > 0 {
			s.data = append(s.data, s.buf[:n]...)
		}
		if err != nil {
			if len(s.data) > s.pos {
				line := string(s.data[s.pos:])
				s.pos = len(s.data)
				return line, nil
			}
			return "", err
		}
	}
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2022))
	var tests []testCase
	total := 0
	add := func(tc testCase) {
		if total+tc.n > 100000 {
			return
		}
		tests = append(tests, tc)
		total += tc.n
	}

	add(makeCase([]int{0, 1, 0, -1, 0, 1, 0}))
	add(makeCase([]int{0, 1, -1, 0}))
	add(makeCase([]int{1, 1, 1, -1, 0, 0}))
	add(makeCase([]int{-1, 1, 1}))
	add(makeCase([]int{0, 0, 0, -1, 0}))

	for total < 100000 {
		n := rng.Intn(200) + 3
		roles := make([]int, n)
		pos := rng.Intn(n)
		for i := range roles {
			if i == pos {
				roles[i] = -1
			} else {
				if rng.Intn(2) == 0 {
					roles[i] = 0
				} else {
					roles[i] = 1
				}
			}
		}
		add(makeCase(roles))
		if len(tests) > 500 {
			break
		}
	}

	if total < 100000 {
		roles := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			roles[i] = 1
		}
		roles[0] = -1
		add(makeCase(roles))
	}

	return tests
}

func makeCase(roles []int) testCase {
	cp := make([]int, len(roles))
	copy(cp, roles)
	return testCase{n: len(cp), roles: cp}
}

func impostorIndex(arr []int) int {
	for i, v := range arr {
		if v == -1 {
			return i + 1
		}
	}
	return -1
}
