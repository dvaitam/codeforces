package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runExe(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func encode(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func bfs(start []int) map[string]int {
	key := encode(start)
	dist := map[string]int{key: 0}
	q := [][]int{append([]int(nil), start...)}
	for front := 0; front < len(q); front++ {
		cur := q[front]
		d := dist[encode(cur)]
		n := len(cur)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				cur[i], cur[j] = cur[j], cur[i]
				k := encode(cur)
				if _, ok := dist[k]; !ok {
					cp := append([]int(nil), cur...)
					dist[k] = d + 1
					q = append(q, cp)
				}
				cur[i], cur[j] = cur[j], cur[i]
			}
		}
	}
	return dist
}

func genCase(rng *rand.Rand) (string, []int, []int) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n) + 1
	}
	b := make([]int, n)
	copy(b, a)
	rng.Shuffle(n, func(i, j int) { b[i], b[j] = b[j], b[i] })

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, a, b := genCase(rng)
		dist := bfs(a)
		maxSad := 0
		for _, d := range dist {
			if d > maxSad {
				maxSad = d
			}
		}
		sadB := dist[encode(b)]
		want := "WA"
		if sadB == maxSad {
			want = "AC"
		}
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out = strings.Fields(out)[0]
		if out != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
