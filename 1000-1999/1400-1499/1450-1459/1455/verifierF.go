package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type item struct {
	idx int
	ch  byte
}

func brute(n, k int, s string) string {
	arr := make([]item, n)
	for i := 0; i < n; i++ {
		arr[i] = item{i, s[i]}
	}
	best := ""
	var rec func(int, []item)
	rec = func(pos int, cur []item) {
		if pos == n {
			var sb strings.Builder
			for i := 0; i < n; i++ {
				sb.WriteByte(cur[i].ch)
			}
			cand := sb.String()
			if best == "" || cand < best {
				best = cand
			}
			return
		}
		j := 0
		for cur[j].idx != pos {
			j++
		}
		tmp := make([]item, n)
		copy(tmp, cur)
		rec(pos+1, tmp)
		if j > 0 {
			copy(tmp, cur)
			tmp[j], tmp[j-1] = tmp[j-1], tmp[j]
			rec(pos+1, tmp)
		}
		if j < n-1 {
			copy(tmp, cur)
			tmp[j], tmp[j+1] = tmp[j+1], tmp[j]
			rec(pos+1, tmp)
		}
		copy(tmp, cur)
		tmp[j].ch = byte(int(tmp[j].ch-'a'-1+k)%k + 'a')
		rec(pos+1, tmp)
		copy(tmp, cur)
		tmp[j].ch = byte(int(tmp[j].ch-'a'+1)%k + 'a')
		rec(pos+1, tmp)
	}
	rec(0, arr)
	return best
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	k := rng.Intn(4) + 2
	bytesArr := make([]byte, n)
	for i := 0; i < n; i++ {
		bytesArr[i] = byte('a' + rng.Intn(k))
	}
	s := string(bytesArr)
	return fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		var n, k int
		var s string
		fmt.Sscanf(strings.Split(in, "\n")[1], "%d %d", &n, &k)
		s = strings.TrimSpace(strings.Split(in, "\n")[2])
		expected := brute(n, k, s)
		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
