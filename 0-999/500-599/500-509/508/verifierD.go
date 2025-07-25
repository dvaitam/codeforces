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

func solveD(n int, subs []string) (bool, string) {
	const maxv = 200 * 200
	graph := make([][]int, maxv)
	indeg := make([]int, maxv)
	outdeg := make([]int, maxv)
	for _, s := range subs {
		u1 := int(s[0])*200 + int(s[1])
		u2 := int(s[1])*200 + int(s[2])
		graph[u1] = append(graph[u1], u2)
		outdeg[u1]++
		indeg[u2]++
	}
	start := -1
	cntPos := 0
	cntNeg := 0
	for i := 0; i < maxv; i++ {
		diff := outdeg[i] - indeg[i]
		if diff == 1 {
			start = i
			cntPos++
		} else if diff == -1 {
			cntNeg++
		} else if diff != 0 {
			if diff > 1 || diff < -1 {
				return false, ""
			}
		}
	}
	if !(cntPos == cntNeg && (cntPos == 0 || cntPos == 1)) {
		return false, ""
	}
	if start == -1 {
		for i := 0; i < maxv; i++ {
			if outdeg[i] > 0 {
				start = i
				break
			}
		}
	}
	if start == -1 {
		return false, ""
	}
	stack := []int{start}
	var path []int
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		if len(graph[v]) > 0 {
			nxt := graph[v][len(graph[v])-1]
			graph[v] = graph[v][:len(graph[v])-1]
			stack = append(stack, nxt)
		} else {
			stack = stack[:len(stack)-1]
			path = append(path, v)
		}
	}
	if len(path) != n+1 {
		return false, ""
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	ans := make([]byte, 0, n+2)
	ans = append(ans, byte(path[0]/200))
	for _, v := range path {
		ans = append(ans, byte(v%200))
	}
	return true, string(ans)
}

func validD(ans string, subs []string) bool {
	if len(ans) != len(subs)+2 {
		return false
	}
	cnt := make(map[string]int)
	for _, s := range subs {
		cnt[s]++
	}
	for i := 0; i+3 <= len(ans); i++ {
		sub := ans[i : i+3]
		if cnt[sub] == 0 {
			return false
		}
		cnt[sub]--
	}
	for _, c := range cnt {
		if c != 0 {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, []string, bool) {
	n := rng.Intn(18) + 2
	subs := make([]string, n)
	chars := []byte("abcXYZ012")
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		b := []byte{chars[rng.Intn(len(chars))], chars[rng.Intn(len(chars))], chars[rng.Intn(len(chars))]}
		subs[i] = string(b)
		sb.WriteString(subs[i])
		sb.WriteByte('\n')
	}
	ok, _ := solveD(n, subs)
	return sb.String(), subs, ok
}

func runCase(bin, input string, subs []string, expected bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := scanner.Text()
	if !expected {
		if strings.ToUpper(first) != "NO" {
			return fmt.Errorf("expected NO got %s", first)
		}
		return nil
	}
	if strings.ToUpper(first) != "YES" {
		return fmt.Errorf("expected YES got %s", first)
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing sequence")
	}
	ans := scanner.Text()
	if !validD(ans, subs) {
		return fmt.Errorf("invalid sequence %s", ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, subs, ok := generateCase(rng)
		if err := runCase(bin, input, subs, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
