package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const (
	seed uint64 = 357
	base uint64 = 10317
)

func hsv(s string) uint64 {
	var h uint64
	for i := 0; i < len(s); i++ {
		h = h*seed + uint64(s[i])
	}
	return h
}

func solveC(urls []string) string {
	hostPaths := make(map[string]map[string]struct{}, len(urls))
	for _, url := range urls {
		s := strings.TrimPrefix(url, "http://")
		j := 0
		for j < len(s) {
			c := s[j]
			if (c >= 'a' && c <= 'z') || c == '.' {
				j++
			} else {
				break
			}
		}
		host := s[:j]
		k := j
		for k < len(s) {
			c := s[k]
			if (c >= 'a' && c <= 'z') || c == '.' || c == '/' {
				k++
			} else {
				break
			}
		}
		path := s[j:k]
		if len(path) == 0 {
			path = "$"
		}
		m, ok := hostPaths[host]
		if !ok {
			m = make(map[string]struct{})
			hostPaths[host] = m
		}
		m[path] = struct{}{}
	}
	groups := make(map[uint64][]string)
	for host, pathsMap := range hostPaths {
		paths := make([]string, 0, len(pathsMap))
		for p := range pathsMap {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		var hv uint64
		for _, p := range paths {
			hv = hv*base + hsv(p)
		}
		groups[hv] = append(groups[hv], host)
	}
	keys := make([]uint64, 0)
	for hv, hs := range groups {
		if len(hs) > 1 {
			keys = append(keys, hv)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(len(keys)))
	for _, hv := range keys {
		hs := groups[hv]
		sort.Strings(hs)
		for i, host := range hs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString("http://")
			sb.WriteString(host)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randHost(rng *rand.Rand) string {
	letters := []byte{'a', 'b', 'c', 'd', 'e'}
	l := rng.Intn(3) + 1
	var sb strings.Builder
	for i := 0; i < l; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	if rng.Intn(2) == 0 {
		sb.WriteByte('.')
		l2 := rng.Intn(3) + 1
		for i := 0; i < l2; i++ {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
	}
	return sb.String()
}

func randPath(rng *rand.Rand) string {
	letters := []byte{'a', 'b', 'c'}
	if rng.Intn(2) == 0 {
		return ""
	}
	segs := rng.Intn(2) + 1
	var sb strings.Builder
	for s := 0; s < segs; s++ {
		sb.WriteByte('/')
		l := rng.Intn(3) + 1
		for i := 0; i < l; i++ {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	urls := make([]string, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		h := randHost(rng)
		p := randPath(rng)
		url := "http://" + h + p
		urls[i] = url
		sb.WriteString(url)
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solveC(urls)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
