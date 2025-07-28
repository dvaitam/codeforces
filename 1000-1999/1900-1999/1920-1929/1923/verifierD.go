package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	n int
	a []int64
}

type state struct {
	arr []int64
	idx []int
	t   int
}

func encode(arr []int64, idx []int) string {
	var sb strings.Builder
	for i := range arr {
		sb.WriteString(strconv.FormatInt(arr[i], 10))
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(idx[i]))
		sb.WriteByte('|')
	}
	return sb.String()
}

func solve(a []int64) []int {
	n := len(a)
	best := make([]int, n)
	for i := range best {
		best[i] = -1
	}
	start := state{arr: append([]int64(nil), a...), idx: make([]int, n), t: 0}
	for i := range start.idx {
		start.idx[i] = i
	}
	visited := map[string]bool{encode(start.arr, start.idx): true}
	q := list.New()
	q.PushBack(start)
	for q.Len() > 0 {
		v := q.Remove(q.Front()).(state)
		arr, idxs, t := v.arr, v.idx, v.t
		for i := 0; i < len(arr)-1; i++ {
			if arr[i] > arr[i+1] {
				na := make([]int64, len(arr)-1)
				ni := make([]int, len(idxs)-1)
				copy(na[:i], arr[:i])
				na[i] = arr[i] + arr[i+1]
				copy(na[i+1:], arr[i+2:])
				copy(ni[:i], idxs[:i])
				ni[i] = idxs[i]
				copy(ni[i+1:], idxs[i+2:])
				if best[idxs[i+1]] == -1 || best[idxs[i+1]] > t+1 {
					best[idxs[i+1]] = t + 1
				}
				key := encode(na, ni)
				if !visited[key] {
					visited[key] = true
					q.PushBack(state{na, ni, t + 1})
				}
			}
			if arr[i+1] > arr[i] {
				na := make([]int64, len(arr)-1)
				ni := make([]int, len(idxs)-1)
				copy(na[:i], arr[:i])
				na[i] = arr[i] + arr[i+1]
				copy(na[i+1:], arr[i+2:])
				copy(ni[:i], idxs[:i])
				ni[i] = idxs[i+1]
				copy(ni[i+1:], idxs[i+2:])
				if best[idxs[i]] == -1 || best[idxs[i]] > t+1 {
					best[idxs[i]] = t + 1
				}
				key := encode(na, ni)
				if !visited[key] {
					visited[key] = true
					q.PushBack(state{na, ni, t + 1})
				}
			}
		}
	}
	return best
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(5) + 1)
		}
		tests = append(tests, test{n, arr})
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(t.n))
		sb.WriteString("\n")
		for j, v := range t.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteString("\n")
		expectedArr := solve(t.a)
		expected := make([]string, len(expectedArr))
		for j, v := range expectedArr {
			expected[j] = strconv.Itoa(v)
		}
		exp := strings.Join(expected, " ")
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
