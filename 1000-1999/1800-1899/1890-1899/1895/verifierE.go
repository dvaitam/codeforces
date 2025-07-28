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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n int, ax, ay []int, bx, by []int) string {
	total := n + len(bx)
	edges := make([][]int, total)
	rev := make([][]int, total)
	outdeg := make([]int, total)

	m := len(bx)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if bx[j] > ay[i] {
				edges[i] = append(edges[i], n+j)
				rev[n+j] = append(rev[n+j], i)
			}
		}
		outdeg[i] = len(edges[i])
	}
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if ax[i] > by[j] {
				edges[n+j] = append(edges[n+j], i)
				rev[i] = append(rev[i], n+j)
			}
		}
		outdeg[n+j] = len(edges[n+j])
	}

	state := make([]int, total)
	queue := make([]int, 0)
	for i := 0; i < total; i++ {
		if outdeg[i] == 0 {
			state[i] = 2
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, u := range rev[v] {
			if state[u] != 0 {
				continue
			}
			if state[v] == 2 {
				state[u] = 1
				queue = append(queue, u)
			} else {
				outdeg[u]--
				if outdeg[u] == 0 {
					state[u] = 2
					queue = append(queue, u)
				}
			}
		}
	}
	win, draw, lose := 0, 0, 0
	for i := 0; i < n; i++ {
		if state[i] == 1 {
			win++
		} else if state[i] == 2 {
			lose++
		} else {
			draw++
		}
	}
	return fmt.Sprintf("%d %d %d", win, draw, lose)
}

func genCase(rng *rand.Rand) (int, []int, []int, []int, []int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	ax := make([]int, n)
	ay := make([]int, n)
	for i := 0; i < n; i++ {
		ax[i] = rng.Intn(50)
		ay[i] = rng.Intn(50)
	}
	bx := make([]int, m)
	by := make([]int, m)
	for i := 0; i < m; i++ {
		bx[i] = rng.Intn(50)
		by[i] = rng.Intn(50)
	}
	return n, ax, ay, bx, by
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, ax, ay, bx, by := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(ax[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(ay[j]))
		}
		sb.WriteByte('\n')
		m := len(bx)
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(bx[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(by[j]))
		}
		sb.WriteByte('\n')
		expect := solveE(n, ax, ay, bx, by)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
