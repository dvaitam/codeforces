package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func buildSong(s0, t string, k int) string {
	s := s0
	for i := 0; i < k; i++ {
		s = s + string(t[i]) + s
		if len(s) > 2000 { // avoid blowup
			break
		}
	}
	return s
}

func countOccurrences(s, w string) int64 {
	count := 0
	for i := 0; i+len(w) <= len(s); i++ {
		if s[i:i+len(w)] == w {
			count++
		}
	}
	return int64(count)
}

func solveG(n int, s0, t string, queries [][2]string) string {
	res := make([]string, len(queries))
	for i, q := range queries {
		k := 0
		fmt.Sscan(q[0], &k)
		w := q[1]
		song := buildSong(s0, t, k)
		cnt := countOccurrences(song, w) % mod
		res[i] = fmt.Sprint(cnt)
	}
	return strings.Join(res, "\n")
}

func genCases() []string {
	rand.Seed(7)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 1
		s0Len := rand.Intn(3) + 1
		s0 := make([]byte, s0Len)
		for j := range s0 {
			s0[j] = byte('a' + rand.Intn(2))
		}
		t := make([]byte, n)
		for j := 0; j < n; j++ {
			t[j] = byte('a' + rand.Intn(2))
		}
		q := rand.Intn(3) + 1
		queries := make([][2]string, q)
		for j := 0; j < q; j++ {
			k := rand.Intn(n + 1)
			wLen := rand.Intn(3) + 1
			w := make([]byte, wLen)
			for p := range w {
				w[p] = byte('a' + rand.Intn(2))
			}
			queries[j] = [2]string{fmt.Sprint(k), string(w)}
		}
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		sb.WriteString(fmt.Sprintf("%s\n%s\n", string(s0), string(t)))
		for _, qu := range queries {
			sb.WriteString(qu[0])
			sb.WriteByte(' ')
			sb.WriteString(qu[1])
			sb.WriteByte('\n')
		}
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n, q int
		fmt.Sscan(lines[0], &n, &q)
		s0 := lines[1]
		t := lines[2]
		queries := make([][2]string, q)
		for j := 0; j < q; j++ {
			parts := strings.Fields(lines[3+j])
			queries[j] = [2]string{parts[0], parts[1]}
		}
		want := solveG(n, s0, t, queries)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
