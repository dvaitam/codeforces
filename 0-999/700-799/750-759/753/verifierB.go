package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func bullsAndCows(secret, guess string) (int, int) {
	bulls := 0
	for i := 0; i < 4; i++ {
		if secret[i] == guess[i] {
			bulls++
		}
	}
	countS := make(map[byte]int)
	countG := make(map[byte]int)
	for i := 0; i < 4; i++ {
		countS[secret[i]]++
		countG[guess[i]]++
	}
	cows := 0
	for d, cs := range countS {
		if cg, ok := countG[d]; ok {
			if cs < cg {
				cows += cs
			} else {
				cows += cg
			}
		}
	}
	cows -= bulls
	return bulls, cows
}

func generateSecrets() []string {
	secrets := make([]string, 0, 5040)
	digits := []byte("0123456789")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 10; k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < 10; l++ {
					if l == i || l == j || l == k {
						continue
					}
					secrets = append(secrets, string([]byte{digits[i], digits[j], digits[k], digits[l]}))
				}
			}
		}
	}
	if len(secrets) > 100 {
		secrets = secrets[:100]
	}
	return secrets
}

func runCase(bin, secret string, limit int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
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
	cleanup := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	}
	reader := bufio.NewReader(stdout)
	queries := 0
	for {
		var guess string
		_, err = fmt.Fscan(reader, &guess)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("time limit")
			}
			if errors.Is(err, io.EOF) {
				cleanup()
				return fmt.Errorf("program terminated before guessing secret stderr:%s", stderr.String())
			}
			cleanup()
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		if len(guess) != 4 {
			cleanup()
			return fmt.Errorf("invalid guess %q", guess)
		}
		for i := 0; i < 4; i++ {
			if guess[i] < '0' || guess[i] > '9' {
				cleanup()
				return fmt.Errorf("invalid guess %q", guess)
			}
		}
		bulls, cows := bullsAndCows(secret, guess)
		fmt.Fprintf(stdin, "%d %d\n", bulls, cows)
		queries++
		if bulls == 4 {
			stdin.Close()
			err := cmd.Wait()
			if err != nil {
				return fmt.Errorf("program error: %v stderr:%s", err, stderr.String())
			}
			if queries > limit {
				return fmt.Errorf("too many queries: %d", queries)
			}
			return nil
		}
		if queries >= limit {
			cleanup()
			return fmt.Errorf("too many queries: %d", queries)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	secrets := generateSecrets()

	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8
	}
	if workers < 1 {
		workers = 1
	}
	sem := make(chan struct{}, workers)
	var mu sync.Mutex
	var firstErr string
	var wg sync.WaitGroup

	for i, s := range secrets {
		if mu.Lock(); firstErr != "" {
			mu.Unlock()
			break
		}
		mu.Unlock()

		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, secret string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := runCase(bin, secret, 50); err != nil {
				mu.Lock()
				if firstErr == "" {
					firstErr = fmt.Sprintf("case %d failed: %v (secret %s)", idx+1, err, secret)
				}
				mu.Unlock()
			}
		}(i, s)
	}
	wg.Wait()

	if firstErr != "" {
		fmt.Fprintln(os.Stderr, firstErr)
		fmt.Println(firstErr)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
