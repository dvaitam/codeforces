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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type slot struct{ i, j int }

func expectedC(k int, s string) string {
	n := len(s)
	res := []byte(s)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		if res[i] != '?' && res[j] != '?' {
			if res[i] != res[j] {
				return "IMPOSSIBLE"
			}
		} else if res[i] == '?' && res[j] != '?' {
			res[i] = res[j]
		} else if res[i] != '?' && res[j] == '?' {
			res[j] = res[i]
		}
	}
	used := make([]bool, k)
	for _, c := range res {
		if c >= 'a' && int(c-'a') < k {
			used[c-'a'] = true
		}
	}
	var slots []slot
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		if res[i] == '?' && res[j] == '?' {
			slots = append(slots, slot{i, j})
		}
	}
	if n%2 == 1 && res[n/2] == '?' {
		slots = append(slots, slot{n / 2, n / 2})
	}
	var missing []byte
	for i := 0; i < k; i++ {
		if !used[i] {
			missing = append(missing, byte('a'+i))
		}
	}
	if len(missing) > len(slots) {
		return "IMPOSSIBLE"
	}
	for idx, c := range missing {
		sl := slots[idx]
		res[sl.i] = c
		res[sl.j] = c
	}
	for idx := len(missing); idx < len(slots); idx++ {
		sl := slots[idx]
		res[sl.i] = 'a'
		res[sl.j] = 'a'
	}
	for _, c := range res {
		if c == '?' {
			return "IMPOSSIBLE"
		}
	}
	finalUsed := make([]bool, k)
	for _, c := range res {
		if int(c-'a') < k {
			finalUsed[c-'a'] = true
		}
	}
	for i := 0; i < k; i++ {
		if !finalUsed[i] {
			return "IMPOSSIBLE"
		}
	}
	return string(res)
}

func generateCaseC(rng *rand.Rand) (string, string) {
	k := rng.Intn(5) + 1
	n := rng.Intn(20) + k
	letters := make([]byte, n)
	for i := 0; i < n; i++ {
		r := rng.Intn(k + 1) // include ?
		if r == k {
			letters[i] = '?'
		} else {
			letters[i] = byte('a' + r)
		}
	}
	s := string(letters)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", k))
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String(), expectedC(k, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
