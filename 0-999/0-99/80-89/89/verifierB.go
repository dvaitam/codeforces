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

type Widget struct {
	kind     int //0 simple,1 HBox,2 VBox
	w, h     int
	border   int
	spacing  int
	children []string
}

func compute(name string, widgets map[string]*Widget, memoW, memoH map[string]int) (int, int) {
	if w, ok := memoW[name]; ok {
		return w, memoH[name]
	}
	wgt := widgets[name]
	var w, h int
	switch wgt.kind {
	case 0:
		w, h = wgt.w, wgt.h
	case 1, 2:
		n := len(wgt.children)
		if n == 0 {
			w, h = 0, 0
		} else {
			sumW, sumH := 0, 0
			maxW, maxH := 0, 0
			for _, ch := range wgt.children {
				cw, chh := compute(ch, widgets, memoW, memoH)
				sumW += cw
				sumH += chh
				if cw > maxW {
					maxW = cw
				}
				if chh > maxH {
					maxH = chh
				}
			}
			if wgt.kind == 1 {
				w = sumW + wgt.spacing*(n-1) + 2*wgt.border
				h = maxH + 2*wgt.border
			} else {
				w = maxW + 2*wgt.border
				h = sumH + wgt.spacing*(n-1) + 2*wgt.border
			}
		}
	}
	memoW[name] = w
	memoH[name] = h
	return w, h
}

type scriptCase struct {
	n       int
	lines   []string
	widgets map[string]*Widget
	names   []string
}

func genCaseB(rng *rand.Rand) scriptCase {
	wc := rng.Intn(5) + 3
	lines := make([]string, 0)
	widgets := make(map[string]*Widget)
	names := make([]string, 0)
	for i := 0; i < wc; i++ {
		name := fmt.Sprintf("w%d", i)
		names = append(names, name)
		kind := rng.Intn(3)
		switch kind {
		case 0:
			w := rng.Intn(10)
			h := rng.Intn(10)
			lines = append(lines, fmt.Sprintf("Widget %s(%d,%d)", name, w, h))
			widgets[name] = &Widget{kind: 0, w: w, h: h}
		case 1:
			lines = append(lines, fmt.Sprintf("HBox %s", name))
			widgets[name] = &Widget{kind: 1, border: rng.Intn(4), spacing: rng.Intn(3)}
			lines = append(lines, fmt.Sprintf("%s.set_border(%d)", name, widgets[name].border))
			lines = append(lines, fmt.Sprintf("%s.set_spacing(%d)", name, widgets[name].spacing))
		case 2:
			lines = append(lines, fmt.Sprintf("VBox %s", name))
			widgets[name] = &Widget{kind: 2, border: rng.Intn(4), spacing: rng.Intn(3)}
			lines = append(lines, fmt.Sprintf("%s.set_border(%d)", name, widgets[name].border))
			lines = append(lines, fmt.Sprintf("%s.set_spacing(%d)", name, widgets[name].spacing))
		}
	}
	// pack operations
	for i, name := range names {
		wgt := widgets[name]
		if wgt.kind == 0 {
			continue
		}
		childCount := rng.Intn(len(names))
		for j := 0; j < childCount; j++ {
			idx := rng.Intn(i + 1) // ensure child created before container
			if names[idx] == name {
				continue
			}
			wgt.children = append(wgt.children, names[idx])
			lines = append(lines, fmt.Sprintf("%s.pack(%s)", name, names[idx]))
		}
	}
	return scriptCase{n: len(lines), lines: lines, widgets: widgets, names: names}
}

func solveB(sc scriptCase) string {
	memoW := make(map[string]int)
	memoH := make(map[string]int)
	names := append([]string(nil), sc.names...)
	sort.Strings(names)
	var sb strings.Builder
	for _, name := range names {
		w, h := compute(name, sc.widgets, memoW, memoH)
		sb.WriteString(fmt.Sprintf("%s %d %d\n", name, w, h))
	}
	return sb.String()
}

func runCaseB(bin string, sc scriptCase) error {
	input := fmt.Sprintf("%d\n%s\n", sc.n, strings.Join(sc.lines, "\n"))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	expected := solveB(sc)
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%sbut got:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		sc := genCaseB(rng)
		if err := runCaseB(bin, sc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
