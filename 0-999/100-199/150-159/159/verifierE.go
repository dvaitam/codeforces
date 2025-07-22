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

type Cube struct {
	id    int
	color int
	size  int
}

func solveE(cubes []Cube) string {
	colors := make(map[int]struct{})
	for _, c := range cubes {
		colors[c.color] = struct{}{}
	}
	bestHeight := 0
	var bestSeq []int
	colorList := make([]int, 0, len(colors))
	for c := range colors {
		colorList = append(colorList, c)
	}
	for i := 0; i < len(colorList); i++ {
		for j := i + 1; j < len(colorList); j++ {
			c1 := colorList[i]
			c2 := colorList[j]
			var list1, list2 []Cube
			for _, c := range cubes {
				if c.color == c1 {
					list1 = append(list1, c)
				} else if c.color == c2 {
					list2 = append(list2, c)
				}
			}
			sort.Slice(list1, func(i, j int) bool { return list1[i].size > list1[j].size })
			sort.Slice(list2, func(i, j int) bool { return list2[i].size > list2[j].size })
			for start := 0; start < 2; start++ {
				h := 0
				seq := []int{}
				i1, i2 := 0, 0
				cur := start
				for {
					if cur == 0 {
						if i1 >= len(list1) {
							break
						}
						seq = append(seq, list1[i1].id)
						h += list1[i1].size
						i1++
						cur = 1
					} else {
						if i2 >= len(list2) {
							break
						}
						seq = append(seq, list2[i2].id)
						h += list2[i2].size
						i2++
						cur = 0
					}
				}
				if h > bestHeight && len(seq) >= 2 {
					bestHeight = h
					bestSeq = append([]int(nil), seq...)
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n", bestHeight, len(bestSeq)))
	for i, id := range bestSeq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", id))
	}
	sb.WriteString("\n")
	return strings.TrimSpace(sb.String())
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		cubes := make([]Cube, n)
		for j := 0; j < n; j++ {
			cubes[j] = Cube{id: j + 1, color: rng.Intn(3) + 1, size: rng.Intn(9) + 1}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, c := range cubes {
			sb.WriteString(fmt.Sprintf("%d %d\n", c.color, c.size))
		}
		input := sb.String()
		expected := solveE(cubes)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
