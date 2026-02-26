package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runTest(bin string, n, t int, arr []int, queries []int) error {
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

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	inScan := bufio.NewScanner(stdout)
	
	fmt.Fprintf(stdin, "%d %d\n", n, t)

	queryCount := 0

	for _, k := range queries {
		fmt.Fprintf(stdin, "%d\n", k)

		for {
			if !inScan.Scan() {
				return fmt.Errorf("unexpected EOF from candidate")
			}
			line := strings.TrimSpace(inScan.Text())
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}

			if parts[0] == "?" {
				queryCount++
				if queryCount > 60000 {
					return fmt.Errorf("too many queries: %d", queryCount)
				}
				if len(parts) != 3 {
					return fmt.Errorf("invalid query format: %s", line)
				}
				l, err1 := strconv.Atoi(parts[1])
				r, err2 := strconv.Atoi(parts[2])
				if err1 != nil || err2 != nil || l < 1 || r > n || l > r {
					return fmt.Errorf("invalid query range: %d %d", l, r)
				}

				sum := 0
				for i := l; i <= r; i++ {
					sum += arr[i-1]
				}
				fmt.Fprintf(stdin, "%d\n", sum)
			} else if parts[0] == "!" {
				if len(parts) != 2 {
					return fmt.Errorf("invalid answer format: %s", line)
				}
				x, err := strconv.Atoi(parts[1])
				if err != nil || x < 1 || x > n {
					return fmt.Errorf("invalid answer index: %d", x)
				}

				// Find k-th zero
				zeros := 0
				expectedX := -1
				for i := 0; i < n; i++ {
					if arr[i] == 0 {
						zeros++
						if zeros == k {
							expectedX = i + 1
							break
						}
					}
				}

				if x != expectedX {
					return fmt.Errorf("wrong answer: expected %d, got %d", expectedX, x)
				}
				arr[x-1] = 1
				break
			} else {
				return fmt.Errorf("unknown command: %s", line)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 20; i++ {
		n := rng.Intn(200000) + 1
		t := rng.Intn(500) + 1
		if t > n {
			t = n
		}

		arr := make([]int, n)
		zerosCount := n
		for j := 0; j < n; j++ {
			if rng.Intn(4) == 0 && zerosCount > t {
				arr[j] = 1
				zerosCount--
			}
		}

		queries := make([]int, t)
		for j := 0; j < t; j++ {
			queries[j] = rng.Intn(zerosCount) + 1
			zerosCount--
		}

		if err := runTest(bin, n, t, arr, queries); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
