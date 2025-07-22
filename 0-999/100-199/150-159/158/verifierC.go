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

func solveC(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	n := 0
	fmt.Sscan(lines[0], &n)
	path := make([]string, 0)
	var out []string
	for i := 1; i <= n; i++ {
		fields := strings.Fields(lines[i])
		if fields[0] == "pwd" {
			if len(path) == 0 {
				out = append(out, "/")
			} else {
				out = append(out, "/"+strings.Join(path, "/")+"/")
			}
		} else {
			p := fields[1]
			if strings.HasPrefix(p, "/") {
				path = path[:0]
			}
			parts := strings.Split(p, "/")
			for _, part := range parts {
				if part == "" {
					continue
				}
				if part == ".." {
					if len(path) > 0 {
						path = path[:len(path)-1]
					}
				} else {
					path = append(path, part)
				}
			}
		}
	}
	return strings.Join(out, "\n") + "\n"
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
	err := cmd.Run()
	return out.String(), err
}

func randName(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	path := make([]string, 0)
	var lines []string
	hasPwd := false
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 { // pwd
			lines = append(lines, "pwd")
			hasPwd = true
		} else { // cd
			absolute := rng.Intn(2) == 0
			var segs []string
			temp := make([]string, len(path))
			copy(temp, path)
			count := rng.Intn(5)
			if absolute {
				temp = temp[:0]
			}
			for j := 0; j < count; j++ {
				if rng.Intn(3) == 0 && len(temp) > 0 {
					segs = append(segs, "..")
					temp = temp[:len(temp)-1]
				} else {
					name := randName(rng)
					segs = append(segs, name)
					temp = append(temp, name)
				}
			}
			if absolute {
				path = temp
				if len(segs) == 0 {
					lines = append(lines, "cd /")
				} else {
					lines = append(lines, "cd /"+strings.Join(segs, "/"))
				}
			} else {
				path = temp
				if len(segs) == 0 {
					lines = append(lines, "cd .")
				} else {
					lines = append(lines, "cd "+strings.Join(segs, "/"))
				}
			}
		}
	}
	if !hasPwd {
		lines = append(lines, "pwd")
	}
	return fmt.Sprintf("%d\n%s\n", len(lines), strings.Join(lines, "\n"))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		expected := solveC(input)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%sinput:\n%s", t+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
