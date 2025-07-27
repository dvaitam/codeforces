package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(s string) string {
	adj := make([][]bool, 26)
	for i := range adj {
		adj[i] = make([]bool, 26)
	}
	n := len(s)
	for i := 0; i < n; i++ {
		cur := s[i] - 'a'
		if i > 0 {
			prev := s[i-1] - 'a'
			adj[cur][prev] = true
			adj[prev][cur] = true
		}
		if i < n-1 {
			next := s[i+1] - 'a'
			adj[cur][next] = true
			adj[next][cur] = true
		}
	}
	visited := make([]bool, 26)
	deg := make([]int, 26)
	poss := true
	for i := 0; i < 26; i++ {
		cnt := 0
		for j := 0; j < 26; j++ {
			if adj[i][j] {
				cnt++
			}
		}
		if cnt > 2 {
			poss = false
		}
		deg[i] = cnt
	}
	var ans []byte
	for i := 0; i < 26; i++ {
		if !visited[i] && deg[i] <= 1 {
			stack := [][2]int{{i, -1}}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				cur, parent := top[0], top[1]
				if visited[cur] {
					continue
				}
				visited[cur] = true
				ans = append(ans, byte(cur)+'a')
				for c := 0; c < 26; c++ {
					if adj[cur][c] {
						if c == parent {
							continue
						}
						if visited[c] {
							poss = false
						} else {
							stack = append(stack, [2]int{c, cur})
						}
					}
				}
			}
		}
	}
	for i := 0; i < 26; i++ {
		if !visited[i] {
			ans = append(ans, byte(i)+'a')
		}
	}
	if !poss {
		return "NO\n"
	}
	return fmt.Sprintf("YES\n%s\n", string(ans))
}

func genCase(rng *rand.Rand) (string, string) {
	// generate random string without adjacent equal letters
	length := rng.Intn(20) + 1
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		for {
			ch := byte(rng.Intn(26)) + 'a'
			if i == 0 || ch != b[i-1] {
				b[i] = ch
				break
			}
		}
	}
	s := string(b)
	in := fmt.Sprintf("1\n%s\n", s)
	out := solveCase(s)
	return in, out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
