package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	nQuestions  = 5000
	kWrong      = 2000
	maxAttempts = 100
)

type testCase struct {
	name    string
	answers string
}

func makeAnswers(pattern func(i int) byte) string {
	var b strings.Builder
	b.Grow(nQuestions)
	for i := 0; i < nQuestions; i++ {
		b.WriteByte(pattern(i))
	}
	return b.String()
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(684))
	var tests []testCase
	tests = append(tests,
		testCase{name: "all_zero", answers: makeAnswers(func(i int) byte { return '0' })},
		testCase{name: "all_one", answers: makeAnswers(func(i int) byte { return '1' })},
		testCase{name: "alternating", answers: makeAnswers(func(i int) byte {
			if i%2 == 0 {
				return '0'
			}
			return '1'
		})},
	)
	for i := 0; i < 3; i++ {
		answers := makeAnswers(func(int) byte {
			if rng.Intn(2) == 0 {
				return '0'
			}
			return '1'
		})
		tests = append(tests, testCase{
			name:    fmt.Sprintf("random_%d", i+1),
			answers: answers,
		})
	}
	return tests
}

func calcResult(ans []byte, attempt string) int {
	wrong := 0
	data := []byte(attempt)
	for i := 0; i < nQuestions; i++ {
		if data[i] != ans[i] {
			wrong++
			if wrong == kWrong {
				return i + 1
			}
		}
	}
	return nQuestions + 1
}

func runTest(bin string, tc testCase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture stdout: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to capture stdin: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to capture stderr: %w", err)
	}
	var stderrBuf bytes.Buffer
	go io.Copy(&stderrBuf, stderrPipe)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start binary: %w", err)
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)
	defer writer.Flush()

	attempts := 0
	madeAttempt := false
	answers := []byte(tc.answers)
	best := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed reading attempt output: %w", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if len(line) == 0 {
			return fmt.Errorf("empty attempt detected")
		}
		attempts++
		madeAttempt = true
		if attempts > maxAttempts {
			return fmt.Errorf("exceeded %d attempts", maxAttempts)
		}
		if len(line) != nQuestions {
			return fmt.Errorf("attempt length %d != %d", len(line), nQuestions)
		}
		for i := 0; i < len(line); i++ {
			if line[i] != '0' && line[i] != '1' {
				return fmt.Errorf("invalid character %q in attempt", line[i])
			}
		}
		result := calcResult(answers, line)
		if result > best {
			best = result
		}
		if _, err := fmt.Fprintf(writer, "%d\n", result); err != nil {
			return fmt.Errorf("failed to send result: %w", err)
		}
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush result: %w", err)
		}
		if result == nQuestions+1 {
			// optional: candidate may choose to stop early; continue loop to handle graceful exit.
		}
	}

	writer.Close()
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("process timed out on %s", tc.name)
		}
		return fmt.Errorf("process exited with error: %v\nstderr: %s", err, stderrBuf.String())
	}

	if !madeAttempt {
		return fmt.Errorf("no attempts made for test %s", tc.name)
	}

	if best < kWrong {
		return fmt.Errorf("suspiciously low best result (%d) for test %s", best, tc.name)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := buildTests()
	for idx, tc := range tests {
		if err := runTest(bin, tc); err != nil {
			fmt.Printf("test %d (%s) failed: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
