package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded reference solver for 1819B

type Pair struct {
	h int64
	w int64
}

func refNextInt(data []byte, idx *int) int64 {
	n := len(data)
	for *idx < n {
		c := data[*idx]
		if c >= '0' && c <= '9' {
			break
		}
		*idx++
	}
	var v int64
	for *idx < n {
		c := data[*idx]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		*idx++
	}
	return v
}

func refVerify(a, b []int64, H, W int64) bool {
	n := len(a)
	byH := make(map[int64][]int, n)
	byW := make(map[int64][]int, n)
	for i := 0; i < n; i++ {
		byH[a[i]] = append(byH[a[i]], i)
		byW[b[i]] = append(byW[b[i]], i)
	}
	active := make([]bool, n)
	for i := 0; i < n; i++ {
		active[i] = true
	}
	remain := n

	popH := func(key int64) int {
		s := byH[key]
		for len(s) > 0 && !active[s[len(s)-1]] {
			s = s[:len(s)-1]
		}
		byH[key] = s
		if len(s) == 0 {
			return -1
		}
		return s[len(s)-1]
	}

	popW := func(key int64) int {
		s := byW[key]
		for len(s) > 0 && !active[s[len(s)-1]] {
			s = s[:len(s)-1]
		}
		byW[key] = s
		if len(s) == 0 {
			return -1
		}
		return s[len(s)-1]
	}

	curH, curW := H, W
	for remain > 0 {
		id := popH(curH)
		if id != -1 {
			active[id] = false
			curW -= b[id]
			if curW < 0 {
				return false
			}
			remain--
			continue
		}
		id = popW(curW)
		if id == -1 {
			return false
		}
		active[id] = false
		curH -= a[id]
		if curH < 0 {
			return false
		}
		remain--
	}
	return true
}

func solveReference(input string) string {
	data := []byte(input)
	idx := 0
	t := int(refNextInt(data, &idx))
	var out bytes.Buffer

	for ; t > 0; t-- {
		n := int(refNextInt(data, &idx))
		a := make([]int64, n)
		b := make([]int64, n)
		var sum int64
		var maxA, maxB int64
		for i := 0; i < n; i++ {
			a[i] = refNextInt(data, &idx)
			b[i] = refNextInt(data, &idx)
			sum += a[i] * b[i]
			if a[i] > maxA {
				maxA = a[i]
			}
			if b[i] > maxB {
				maxB = b[i]
			}
		}

		ans := make([]Pair, 0, 2)

		if sum%maxA == 0 {
			H := maxA
			W := sum / maxA
			if refVerify(a, b, H, W) {
				ans = append(ans, Pair{H, W})
			}
		}
		if sum%maxB == 0 {
			H := sum / maxB
			W := maxB
			dup := false
			for _, p := range ans {
				if p.h == H && p.w == W {
					dup = true
					break
				}
			}
			if !dup && refVerify(a, b, H, W) {
				ans = append(ans, Pair{H, W})
			}
		}

		out.WriteString(strconv.Itoa(len(ans)))
		out.WriteByte('\n')
		for _, p := range ans {
			out.WriteString(strconv.FormatInt(p.h, 10))
			out.WriteByte(' ')
			out.WriteString(strconv.FormatInt(p.w, 10))
			out.WriteByte('\n')
		}
	}

	return strings.TrimSpace(out.String())
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	H := rng.Intn(5) + 1
	W := rng.Intn(5) + 1
	rects := [][2]int{{H, W}}
	for len(rects) < n {
		idx := rng.Intn(len(rects))
		h, w := rects[idx][0], rects[idx][1]
		if h == 1 && w == 1 {
			continue
		}
		if (rng.Intn(2) == 0 && h > 1) || w == 1 {
			x := rng.Intn(h-1) + 1
			rects[idx][0] = h - x
			rects = append(rects, [2]int{x, w})
		} else {
			y := rng.Intn(w-1) + 1
			rects[idx][1] = w - y
			rects = append(rects, [2]int{h, y})
		}
	}
	rng.Shuffle(len(rects), func(i, j int) { rects[i], rects[j] = rects[j], rects[i] })
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(rects)))
	for _, r := range rects {
		sb.WriteString(fmt.Sprintf("%d %d\n", r[0], r[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	candPath := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		input := genCase(rng)
		exp := solveReference(input)
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		// Compare as sets of pairs (order may differ)
		expLines := strings.Split(strings.TrimSpace(exp), "\n")
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(expLines) == 0 || len(gotLines) == 0 || strings.TrimSpace(expLines[0]) != strings.TrimSpace(gotLines[0]) {
			fmt.Printf("case %d failed (count mismatch)\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
		if len(expLines) != len(gotLines) {
			fmt.Printf("case %d failed (line count mismatch)\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
		expSet := make(map[string]bool)
		for _, l := range expLines[1:] {
			expSet[strings.TrimSpace(l)] = true
		}
		match := true
		for _, l := range gotLines[1:] {
			if !expSet[strings.TrimSpace(l)] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
