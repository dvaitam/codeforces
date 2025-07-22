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

type testCaseB struct {
	n       int
	friends [][]int
	alex    []int
}

func solveB(tc testCaseB) string {
	n := tc.n
	posF := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		posF[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			posF[i][tc.friends[i-1][j-1]] = j
		}
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = tc.alex[i-1]
	}
	b1 := make([]int, n+1)
	b2 := make([]int, n+1)
	for k := 1; k <= n; k++ {
		cnt := 0
		for t := 1; t <= n; t++ {
			card := a[t]
			if card <= k {
				cnt++
				if cnt == 1 {
					b1[k] = card
				} else if cnt == 2 {
					b2[k] = card
					break
				}
			}
		}
	}
	res := make([]int, n+1)
	for j := 1; j <= n; j++ {
		bestRank := n + 1
		bestK := 1
		for k := 1; k <= n; k++ {
			cj := b1[k]
			if cj == j {
				cj = b2[k]
			}
			if cj == 0 {
				continue
			}
			if posF[j][cj] < bestRank {
				bestRank = posF[j][cj]
				bestK = k
			}
		}
		res[j] = bestK
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(res[i]))
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	friends := make([][]int, n)
	for i := 0; i < n; i++ {
		perm := rng.Perm(n)
		for j := range perm {
			perm[j]++
		}
		friends[i] = perm
	}
	alexPerm := rng.Perm(n)
	for j := range alexPerm {
		alexPerm[j]++
	}
	tc := testCaseB{n: n, friends: friends, alex: alexPerm}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(friends[i][j]))
		}
		sb.WriteByte('\n')
	}
	for j := 0; j < n; j++ {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(alexPerm[j]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveB(tc)
	return input, expect
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
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseB(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
