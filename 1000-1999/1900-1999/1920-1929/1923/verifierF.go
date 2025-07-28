package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	n int
	k int
	s string
}

func shrinkReverse(s string) string {
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	t := []byte(s[i:])
	for l, r := 0, len(t)-1; l < r; l, r = l+1, r-1 {
		t[l], t[r] = t[r], t[l]
	}
	return string(t)
}

func value(s string) int64 {
	var val int64
	for i := 0; i < len(s); i++ {
		val = val*2 + int64(s[i]-'0')
	}
	return val
}

func minValue(n int, k int, s string) int64 {
	type node struct {
		str  string
		step int
	}
	best := value(s)
	seen := map[string]int{s: 0}
	q := []node{{s, 0}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if val := value(cur.str); val < best {
			best = val
		}
		if cur.step == k {
			continue
		}
		// swap operations
		bytes := []byte(cur.str)
		for i := 0; i < len(bytes); i++ {
			for j := i + 1; j < len(bytes); j++ {
				b := append([]byte(nil), bytes...)
				b[i], b[j] = b[j], b[i]
				ns := string(b)
				if p, ok := seen[ns]; !ok || p > cur.step+1 {
					seen[ns] = cur.step + 1
					q = append(q, node{ns, cur.step + 1})
				}
			}
		}
		// shrink reverse
		ns := shrinkReverse(cur.str)
		if p, ok := seen[ns]; !ok || p > cur.step+1 {
			seen[ns] = cur.step + 1
			q = append(q, node{ns, cur.step + 1})
		}
	}
	return best
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(6) + 2
		k := rng.Intn(n) + 1
		b := make([]byte, n)
		hasOne := false
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 1 {
				b[i] = '1'
				hasOne = true
			} else {
				b[i] = '0'
			}
		}
		if !hasOne {
			b[rng.Intn(n)] = '1'
		}
		tests = append(tests, test{n, k, string(b)})
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", t.n, t.k, t.s)
		expected := strconv.FormatInt(minValue(t.n, t.k, t.s), 10)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
