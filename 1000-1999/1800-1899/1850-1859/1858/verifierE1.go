package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded solver for 1858E1
func solve1858E1(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	if !scanner.Scan() {
		return ""
	}
	q, _ := strconv.Atoi(scanner.Text())

	up := make([][20]int, q+1)
	val := make([]int, q+1)

	stack := make([]int, 0, q+1)
	stack = append(stack, 0)

	queries := make([]int, 0, q)

	id := 0

	for i := 0; i < q; i++ {
		scanner.Scan()
		op := scanner.Text()
		if op == "+" {
			scanner.Scan()
			x, _ := strconv.Atoi(scanner.Text())
			id++
			val[id] = x
			p := stack[len(stack)-1]
			up[id][0] = p
			for j := 1; j < 20; j++ {
				up[id][j] = up[up[id][j-1]][j-1]
			}
			stack = append(stack, id)
		} else if op == "-" {
			scanner.Scan()
			k, _ := strconv.Atoi(scanner.Text())
			curr := stack[len(stack)-1]
			for j := 19; j >= 0; j-- {
				if (k & (1 << j)) != 0 {
					curr = up[curr][j]
				}
			}
			stack = append(stack, curr)
		} else if op == "!" {
			stack = stack[:len(stack)-1]
		} else if op == "?" {
			queries = append(queries, stack[len(stack)-1])
		}
	}

	head := make([]int, id+1)
	for i := range head {
		head[i] = -1
	}
	next := make([]int, id+1)

	for i := 1; i <= id; i++ {
		p := up[i][0]
		next[i] = head[p]
		head[p] = i
	}

	ansAt := make([]int, id+1)
	freq := make([]int, 1000005)
	distinct := 0

	var dfs func(int)
	dfs = func(u int) {
		ansAt[u] = distinct
		for v := head[u]; v != -1; v = next[v] {
			x := val[v]
			if freq[x] == 0 {
				distinct++
			}
			freq[x]++

			dfs(v)

			freq[x]--
			if freq[x] == 0 {
				distinct--
			}
		}
	}

	dfs(0)

	var sb strings.Builder
	for i, qNode := range queries {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(ansAt[qNode]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func deterministicCases() []string {
	return []string{
		"5\n+ 1\n+ 2\n?\n- 1\n?\n",
		"6\n+ 1\n+ 1\n+ 2\n?\n!\n?\n",
	}
}

func randomCase(rng *rand.Rand) string {
	q := rng.Intn(30) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	arrLen := 0
	history := 0
	for i := 0; i < q; i++ {
		opType := rng.Intn(4)
		if opType == 0 || arrLen == 0 { // push
			x := rng.Intn(1000000) + 1
			fmt.Fprintf(&sb, "+ %d\n", x)
			arrLen++
			history++
			continue
		}
		if opType == 1 && arrLen > 0 { // pop
			k := rng.Intn(arrLen) + 1
			fmt.Fprintf(&sb, "- %d\n", k)
			arrLen -= k
			history++
			continue
		}
		if opType == 2 && history > 0 { // rollback
			sb.WriteString("!\n")
			history--
			// arrLen cannot exceed q but we ignore exact arr state for generator
			continue
		}
		sb.WriteString("?\n")
	}
	return sb.String()
}

func verify(userBin string, cases []string) {
	for i, in := range cases {
		want := solve1858E1(in)
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	verify(userBin, cases)
}
