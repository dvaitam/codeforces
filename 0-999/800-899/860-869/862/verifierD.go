package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	s string
}

func generateTests() []Test {
	rand.Seed(45)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(999) + 2 // n up to 1000
		for {
			b := make([]byte, n)
			has0, has1 := false, false
			for j := 0; j < n; j++ {
				if rand.Intn(2) == 0 {
					b[j] = '0'
					has0 = true
				} else {
					b[j] = '1'
					has1 = true
				}
			}
			if has0 && has1 {
				tests = append(tests, Test{s: string(b)})
				break
			}
		}
	}
	tests = append(tests, Test{s: "01"})
	tests = append(tests, Test{s: "10"})
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierD <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := generateTests()

	passed := 0
	for i, t := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		cmd := exec.CommandContext(ctx, bin)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Printf("Test %d start err %v\n", i+1, err)
			cancel()
			continue
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Test %d start err %v\n", i+1, err)
			cancel()
			continue
		}
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			fmt.Printf("Test %d start err %v\n", i+1, err)
			cancel()
			continue
		}

		n := len(t.s)
		fmt.Fprintf(stdin, "%d\n", n)

		scanner := bufio.NewScanner(stdout)
		queries := 0
		success := false
		failReason := ""

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}

			if parts[0] == "?" {
				queries++
				if queries > 15 {
					failReason = "Query limit exceeded"
					break
				}
				if len(parts) != 2 || len(parts[1]) != n {
					failReason = "Invalid query format"
					break
				}
				q := parts[1]
				dist := 0
				valid := true
				for j := 0; j < n; j++ {
					if q[j] != '0' && q[j] != '1' {
						valid = false
						break
					}
					if q[j] != t.s[j] {
						dist++
					}
				}
				if !valid {
					failReason = "Invalid query string"
					break
				}
				fmt.Fprintf(stdin, "%d\n", dist)
			} else if parts[0] == "!" {
				if len(parts) != 3 {
					failReason = "Invalid answer format"
					break
				}
				var p0, p1 int
				fmt.Sscanf(parts[1], "%d", &p0)
				fmt.Sscanf(parts[2], "%d", &p1)

				if p0 >= 1 && p0 <= n && p1 >= 1 && p1 <= n {
					if t.s[p0-1] == '0' && t.s[p1-1] == '1' {
						success = true
					} else {
						failReason = fmt.Sprintf("Wrong answer: expected 0 at %d, 1 at %d. string was %s", p0, p1, t.s)
					}
				} else {
					failReason = "Answer indices out of bounds"
				}
				break
			} else {
				failReason = "Invalid command: " + parts[0]
				break
			}
		}

		stdin.Close()
		cmd.Wait()
		
		if !success && failReason == "" {
			if ctx.Err() == context.DeadlineExceeded {
				failReason = "Time limit exceeded"
			} else {
				failReason = "Program exited unexpectedly"
			}
		}

		cancel()

		if success {
			passed++
		} else {
			fmt.Printf("Test %d failed: %s\n", i+1, failReason)
			continue
		}
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
