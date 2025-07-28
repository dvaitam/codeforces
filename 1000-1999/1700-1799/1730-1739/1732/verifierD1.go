package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solve(reader *bufio.Reader) string {
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return ""
	}
	set := make(map[int64]struct{})
	set[0] = struct{}{}
	next := make(map[int64]int64)
	var sb strings.Builder
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var x int64
			fmt.Fscan(reader, &x)
			set[x] = struct{}{}
		} else if op == "?" {
			var k int64
			fmt.Fscan(reader, &k)
			v := next[k]
			for {
				if _, found := set[v]; !found {
					fmt.Fprintf(&sb, "%d\n", v)
					next[k] = v
					break
				}
				v += k
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	q := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	set := map[int64]bool{0: true}
	ask := false
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) // 0 add, 1 query
		if i == q-1 && !ask {
			typ = 1
		}
		if typ == 0 {
			var x int64
			for {
				x = int64(rng.Intn(1000000) + 1)
				if !set[x] {
					break
				}
			}
			set[x] = true
			fmt.Fprintf(&sb, "+ %d\n", x)
		} else {
			k := int64(rng.Intn(10) + 1)
			fmt.Fprintf(&sb, "? %d\n", k)
			ask = true
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
