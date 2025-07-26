package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	reader := bufio.NewReader(stdout)
	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("time limit")
			}
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		guess := line
		if len(guess) != 4 {
			return fmt.Errorf("invalid guess %q", guess)
		}
		for i := 0; i < 4; i++ {
			if guess[i] < '0' || guess[i] > '9' {
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
			return fmt.Errorf("too many queries: %d", queries)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	secrets := generateSecrets()
	for i, s := range secrets {
		if err := runCase(bin, s, 50); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (secret %s)\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
