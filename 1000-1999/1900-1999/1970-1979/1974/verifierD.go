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

func solve(n int, s string) string {
	mp := map[byte][]int{'N': {}, 'S': {}, 'E': {}, 'W': {}}
	for i := 0; i < n; i++ {
		mp[s[i]] = append(mp[s[i]], i)
	}
	x := len(mp['E']) - len(mp['W'])
	y := len(mp['N']) - len(mp['S'])
	if x%2 != 0 || y%2 != 0 {
		return "NO"
	}
	ans := make([]byte, n)
	for i := range ans {
		ans[i] = 'R'
	}
	if x != 0 || y != 0 {
		if x > 0 {
			for i := 0; i < x/2; i++ {
				ans[mp['E'][i]] = 'H'
			}
		} else {
			for i := 0; i < (-x)/2; i++ {
				ans[mp['W'][i]] = 'H'
			}
		}
		if y > 0 {
			for i := 0; i < y/2; i++ {
				ans[mp['N'][i]] = 'H'
			}
		} else {
			for i := 0; i < (-y)/2; i++ {
				ans[mp['S'][i]] = 'H'
			}
		}
	} else {
		if n == 2 {
			return "NO"
		}
		if len(mp['E']) > 0 {
			ans[mp['E'][0]] = 'H'
			ans[mp['W'][0]] = 'H'
		} else {
			ans[mp['N'][0]] = 'H'
			ans[mp['S'][0]] = 'H'
		}
	}
	return string(ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dirs := []byte{'N', 'S', 'E', 'W'}
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = dirs[rng.Intn(len(dirs))]
		}
		s := string(b)
		input := fmt.Sprintf("1\n%d\n%s\n", n, s)
		expected := solve(n, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
