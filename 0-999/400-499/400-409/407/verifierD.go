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

// Correct solver for 407D embedded directly.
func solve407D(input string) string {
	data := []byte(input)
	pos := 0
	readByte := func() byte {
		if pos >= len(data) {
			return 0
		}
		b := data[pos]
		pos++
		return b
	}
	readInt := func() int {
		b := readByte()
		for b > 0 && b <= ' ' {
			b = readByte()
		}
		if b == 0 {
			return 0
		}
		res := 0
		for b > ' ' {
			res = res*10 + int(b-'0')
			b = readByte()
		}
		return res
	}

	n := readInt()
	if n == 0 {
		return "0"
	}
	m := readInt()

	a_t := make([][]int, m)
	for c := 0; c < m; c++ {
		a_t[c] = make([]int, n)
	}
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			a_t[c][r] = readInt()
		}
	}

	first_dup := make([][]int, n)
	seen_r := make([]int, 160005)
	cv_fd := 0
	for i1 := 0; i1 < n; i1++ {
		first_dup[i1] = make([]int, m)
		for c := 0; c < m; c++ {
			cv_fd++
			fd := n
			col := a_t[c]
			for r := i1; r < n; r++ {
				v := col[r]
				if seen_r[v] == cv_fd {
					fd = r
					break
				}
				seen_r[v] = cv_fd
			}
			first_dup[i1][c] = fd
		}
	}

	max_area := 0
	seen := make([]uint32, 160005)
	current_version := 0

	for i1 := 0; i1 < n; i1++ {
		if (n-i1)*m <= max_area {
			break
		}
		for i2 := i1; i2 < n; i2++ {
			h := i2 - i1 + 1
			if h*m <= max_area {
				continue
			}

			max_consec := 0
			consec := 0
			for c := 0; c < m; c++ {
				if first_dup[i1][c] > i2 {
					consec++
					if consec > max_consec {
						max_consec = consec
					}
				} else {
					consec = 0
				}
			}
			if h*max_consec <= max_area {
				continue
			}

			current_version++
			cv_shifted := uint32(current_version) << 10
			limit := 0

			for c := 0; c < m; c++ {
				if h*(m-limit) <= max_area {
					break
				}

				val_to_store := cv_shifted | uint32(c+1)
				col := a_t[c]

				for r := i1; r <= i2; r++ {
					v := col[r]
					s := seen[v]
					seen[v] = val_to_store
					if s > cv_shifted {
						ls := int(s & 1023)
						if ls > limit {
							limit = ls
						}
					}
				}
				cand := h * (c - limit + 1)
				if cand > max_area {
					max_area = cand
				}
			}
		}
	}

	return fmt.Sprintf("%d", max_area)
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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := solve407D(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
