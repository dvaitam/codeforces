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
	"strconv"
	"strings"
	"time"
)

type phrases struct {
	found    string
	closer   string
	further  string
	same     string
	notFound string
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func randomWord(rng *rand.Rand, minLen, maxLen int) string {
	length := rng.Intn(maxLen-minLen+1) + minLen
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func newLanguage(rng *rand.Rand) phrases {
	used := make(map[string]bool)
	next := func(minLen, maxLen int) string {
		for {
			w := randomWord(rng, minLen, maxLen)
			if !used[w] {
				used[w] = true
				return w
			}
		}
	}
	notFound := next(3, 10)
	closer := next(3, 10)
	further := next(3, 10)
	same := next(3, 10)
	found := next(3, 9) + "!"
	return phrases{
		found:    found,
		closer:   closer,
		further:  further,
		same:     same,
		notFound: notFound,
	}
}

func nextInt(scanner *bufio.Scanner) (int, error) {
	for scanner.Scan() {
		text := scanner.Text()
		if v, err := strconv.Atoi(text); err == nil {
			return v, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return 0, io.EOF
}

func runCase(bin string, targetX, targetY int, ph phrases) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() {
		if cmd.ProcessState == nil {
			cmd.Wait()
		}
	}()

	writer := bufio.NewWriter(stdin)
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)

	queries := 0
	var prevDist int64 = -1

	for {
		x, err := nextInt(scanner)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("timeout waiting for query %d", queries+1)
			}
			return fmt.Errorf("failed to read x for query %d: %v stderr:%s", queries+1, err, stderr.String())
		}
		y, err := nextInt(scanner)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("timeout waiting for query %d", queries+1)
			}
			return fmt.Errorf("failed to read y for query %d: %v stderr:%s", queries+1, err, stderr.String())
		}
		queries++
		if queries > 64 {
			return fmt.Errorf("too many queries (>64)")
		}
		if x < 0 || x > 1_000_000 || y < 0 || y > 1_000_000 {
			return fmt.Errorf("query %d out of bounds: %d %d", queries, x, y)
		}

		dx := int64(x - targetX)
		dy := int64(y - targetY)
		dist := dx*dx + dy*dy

		var resp string
		if dist == 0 {
			resp = ph.found
			fmt.Fprintln(writer, resp)
			writer.Flush()
			stdin.Close()
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("runtime error after found: %v stderr:%s", err, stderr.String())
			}
			return nil
		}

		if queries == 1 {
			resp = ph.notFound
		} else if dist < prevDist {
			resp = ph.closer
		} else if dist > prevDist {
			resp = ph.further
		} else {
			resp = ph.same
		}

		if _, err := fmt.Fprintln(writer, resp); err != nil {
			return fmt.Errorf("failed to write response: %v", err)
		}
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush response: %v", err)
		}
		prevDist = dist
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(1773))
	tests := []struct {
		x int
		y int
	}{
		{0, 0},
		{1_000_000, 1_000_000},
		{1_000_000, 0},
		{0, 1_000_000},
		{500_000, 500_000},
		{123_456, 654_321},
		{765_432, 234_567},
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, struct {
			x int
			y int
		}{
			rng.Intn(1_000_001),
			rng.Intn(1_000_001),
		})
	}

	for i, tc := range tests {
		lang := newLanguage(rng)
		if err := runCase(bin, tc.x, tc.y, lang); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (treasure at %d %d): %v\n", i+1, tc.x, tc.y, err)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
