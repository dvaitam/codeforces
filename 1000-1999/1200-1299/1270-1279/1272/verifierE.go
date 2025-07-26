package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testE struct {
	a []int
}

func genTestsE() []testE {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testE, 100)
	for i := range tests {
		n := r.Intn(100) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = r.Intn(n) + 1
		}
		tests[i] = testE{a: arr}
	}
	return tests
}

func solveE(tc testE) []int {
	n := len(tc.a)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		if i-tc.a[i] >= 0 {
			adj[i-tc.a[i]] = append(adj[i-tc.a[i]], i)
		}
		if i+tc.a[i] < n {
			adj[i+tc.a[i]] = append(adj[i+tc.a[i]], i)
		}
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		found := false
		if i-tc.a[i] >= 0 && (tc.a[i-tc.a[i]]%2 != tc.a[i]%2) {
			found = true
		}
		if !found && i+tc.a[i] < n && (tc.a[i+tc.a[i]]%2 != tc.a[i]%2) {
			found = true
		}
		if found {
			dist[i] = 1
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue = append(queue, to)
			}
		}
	}
	return dist
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", len(tc.a))
		for j, v := range tc.a {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expectedSlice := solveE(tc)
		var expBuilder strings.Builder
		for j, v := range expectedSlice {
			if j > 0 {
				expBuilder.WriteByte(' ')
			}
			expBuilder.WriteString(strconv.Itoa(v))
		}
		expected := expBuilder.String()
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input.String(), expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
