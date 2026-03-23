package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const oracleSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Pair struct {
	val, id int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		if !scanner.Scan() {
			return 0
		}
		res, _ := strconv.Atoi(scanner.Text())
		return res
	}

	t := readInt()
	var out strings.Builder
	for tc := 0; tc < t; tc++ {
		n := readInt()
		a := make([]int, n)
		pairs := make([]Pair, n)
		for i := 0; i < n; i++ {
			a[i] = readInt()
			pairs[i] = Pair{val: a[i], id: i}
		}

		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].val < pairs[j].val
		})

		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[pairs[i].id] = i + 1
		}

		inv := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if b[i] > b[j] {
					inv++
				}
			}
		}

		if inv%2 != 0 {
			dupFound := false
			for i := 0; i < n-1; i++ {
				if pairs[i].val == pairs[i+1].val {
					id1, id2 := pairs[i].id, pairs[i+1].id
					b[id1], b[id2] = b[id2], b[id1]
					dupFound = true
					break
				}
			}
			if !dupFound {
				out.WriteString("-1\n")
				continue
			}
		}

		var ops []int
		for v := 1; v <= n-2; v++ {
			p := -1
			for i := 0; i < n; i++ {
				if b[i] == v {
					p = i
					break
				}
			}

			for p > v {
				ops = append(ops, p-1)
				tmp := b[p]
				b[p] = b[p-1]
				b[p-1] = b[p-2]
				b[p-2] = tmp
				p -= 2
			}
			if p == v {
				ops = append(ops, v, v)
				for k := 0; k < 2; k++ {
					tmp := b[v+1]
					b[v+1] = b[v]
					b[v] = b[v-1]
					b[v-1] = tmp
				}
			}
		}

		out.WriteString(strconv.Itoa(len(ops)) + "\n")
		for i, op := range ops {
			if i > 0 {
				out.WriteString(" ")
			}
			out.WriteString(strconv.Itoa(op))
		}
		out.WriteString("\n")
	}
	fmt.Print(out.String())
}
`

func buildOracle() (string, error) {
	dir := os.TempDir()
	src := filepath.Join(dir, fmt.Sprintf("oracle1374F_%d.go", time.Now().UnixNano()))
	if err := os.WriteFile(src, []byte(oracleSource), 0644); err != nil {
		return "", fmt.Errorf("write oracle source: %v", err)
	}
	defer os.Remove(src)
	oracle := src[:len(src)-3]
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 3
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
