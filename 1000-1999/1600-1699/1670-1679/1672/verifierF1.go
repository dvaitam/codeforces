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

func buildRef() (string, error) {
	ref := "./refF1.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1672F1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

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

func parseInts(s string, n int) ([]int, error) {
	fields := strings.Fields(s)
	if len(fields) < n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(7) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n) + 1
	}
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
	return sb.String(), a
}

func isPerm(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	cnt := make(map[int]int)
	for _, v := range a {
		cnt[v]++
	}
	for _, v := range b {
		cnt[v]--
	}
	for _, c := range cnt {
		if c != 0 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, arr := genCase(rng)
		dist := bfs(arr)
		refOut, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		refPerm, err := parseInts(refOut, len(arr))
		if err != nil || !isPerm(arr, refPerm) {
			fmt.Fprintf(os.Stderr, "bad reference output on case %d\n", i+1)
			os.Exit(1)
		}
		maxSad := dist[encode(refPerm)]
		for _, d := range dist {
			if d > maxSad {
				maxSad = d
			}
		}

		candOut, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candPerm, err := parseInts(candOut, len(arr))
		if err != nil {
			fmt.Printf("case %d failed: cannot parse output\ninput:\n%soutput:\n%s\n", i+1, input, candOut)
			os.Exit(1)
		}
		if !isPerm(arr, candPerm) {
			fmt.Printf("case %d failed: output is not a permutation\n", i+1)
			os.Exit(1)
		}
		sad := dist[encode(candPerm)]
		if sad != maxSad {
			fmt.Printf("case %d failed: sadness %d expected %d\n", i+1, sad, maxSad)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
