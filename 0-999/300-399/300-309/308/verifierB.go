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

func runCandidate(bin, input string) (string, error) {
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

func expectedB(n, r, c int, text string) string {
	words := strings.Fields(text)
	if len(words) > n {
		words = words[:n]
	}
	lens := make([]int, len(words))
	for i, w := range words {
		lens[i] = len(w)
	}
	N := len(words)
	nxt := make([]uint32, N+1)
	sum := 0
	j := 0
	for i := 0; i < N; i++ {
		if j < i {
			j = i
			sum = 0
		}
		for j < N && sum+lens[j]+(j-i) <= c {
			sum += lens[j]
			j++
		}
		nxt[i] = uint32(j)
		sum -= lens[i]
	}
	nxt[N] = uint32(N)
	maxB := 0
	tmp := r
	for tmp > 0 {
		maxB++
		tmp >>= 1
	}
	dp := make([][]uint32, maxB)
	dp[0] = nxt
	for b := 1; b < maxB; b++ {
		dp[b] = make([]uint32, N+1)
		prev := dp[b-1]
		cur := dp[b]
		for i := 0; i <= N; i++ {
			cur[i] = prev[prev[i]]
		}
	}
	bestLen := 0
	bestStart := 0
	var bestEnd uint32
	for i := 0; i < N; i++ {
		pos := uint32(i)
		rem := r
		b := 0
		for rem > 0 {
			if rem&1 != 0 {
				pos = dp[b][pos]
			}
			rem >>= 1
			b++
		}
		length := int(pos) - i
		if length > bestLen {
			bestLen = length
			bestStart = i
			bestEnd = pos
		}
		if bestLen == N {
			break
		}
	}
	var outBuf bytes.Buffer
	cur := uint32(bestStart)
	lines := 0
	for lines < r && cur < bestEnd {
		nxtPos := nxt[cur]
		if nxtPos > bestEnd {
			nxtPos = bestEnd
		}
		for k := cur; k < nxtPos; k++ {
			if k > cur {
				outBuf.WriteByte(' ')
			}
			outBuf.WriteString(words[k])
		}
		if lines+1 < r && nxtPos < bestEnd {
			outBuf.WriteByte('\n')
		}
		cur = nxtPos
		lines++
	}
	return strings.TrimSpace(outBuf.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	r := rng.Intn(5) + 1
	c := rng.Intn(15) + 5
	wordCount := n
	words := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		words[i] = string(b)
	}
	text := strings.Join(words, " ")
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, r, c)
	sb.WriteString(text)
	if !strings.HasSuffix(text, "\n") {
		sb.WriteByte('\n')
	}
	expect := expectedB(n, r, c, text)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
