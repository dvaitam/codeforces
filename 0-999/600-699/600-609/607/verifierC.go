package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `7 NENNWW ENNESS
10 SSSSWWNWS NEEESSSSW
6 ENNEN SESEE
4 WWN SWN
8 SSSSEEN NNNWWNN
8 NESWWSW SSEENWS
3 SE EN
8 WNESSEN WSWNESW
7 NWWSSE NWSEEN
10 NNWSSEESE NWSWSSSES
8 SWSWNES SWSEEEN
5 ENWW ESEE
2 E W
2 W E
8 SENNNWS SSEESEN
3 WW SS
4 SSS NWN
9 EEESWSWN WSSEEESW
4 WWW ENE
7 WSWSES NWNNNW
10 NWWWNWWWW NEEEENENN
2 N E
9 ENWSEEES WNNWWNWN
2 S W
5 WWWS ENWN
8 SEENENW WSSWSSE
3 EE WW
8 WSSSWSS EENENNE
8 WWNESSS ENENNWN
4 SSS WWW
9 NEESENWN SSESENNE
7 ENNEES SWSENN
6 SWSWS ESWNN
8 EESSWWS ESESEES
3 WS NE
4 ENW SES
5 NWSE WWSW
5 WWSE SEEE
2 N W
7 WNENWN WNENWS
4 NNW ENN
2 N W
10 NNNWWNNWS SSWWSSWNN
5 WSEN WNEN
6 NWSWS EESWN
10 NENWNENWN SENWWSSEE
10 NWNWSSWSE SEESSESEN
10 ENENESENE NWSENENES
9 SSESEESW SWWNNEEE
3 NN EN
7 SESWNE WNNENE
10 SWSSWWWWW WNWWWSWWN
9 NNEESWWN NNNNNEEE
5 NWSW EESE
5 SEES WSEE
7 WNNNNN WSWWNE
2 N W
5 NEEE WWNW
2 S W
6 WWWNN WNNEE
7 ENWNNW NWWSSE
6 NWSEN WSWWS
5 ENNN EEEE
9 NWSESEEN NEEENNNN
8 SSWSESS ESEEENW
10 NEEENNENW ESSENWWSS
6 NWNWN WNWNE
3 SS SE
5 WSSE EENN
9 ENEEENNW ESWSSSWN
6 NNWSS ESSEE
10 WNNWSENEN WWWSSSSWS
3 ES EN
5 SEEN WSEE
5 ESWS WWNE
7 SWNNWS WNNNEN
4 SEN EES
3 SE ES
8 WSWNNWW ENNNNEN
4 EEN SSW
8 SSWNWWN SWSWWNW
9 SSENEESS WWWSSSSW
2 N N
5 WSWS WSWS
6 NNWNW SSSWN
9 SEEENWWW ENWNEESE
6 ENNEE SWSES
6 WNESW ESWWS
8 WWWSWWN SENEENW
4 WSE WSE
10 NWNESWNNN NENNNWNWS
6 NWWNN SENES
7 ENWSSW NNNWWS
7 SENENN WWNWWW
8 SSSEENW SWNESEE
5 NENE ESEE
7 EENENN EEESWS
8 SWSSWWN NNWSEEE
3 WS EN
7 NNNWNN NESENE`

// solve reproduces the logic from 607C.go to compute the expected answer.
func solve(n int, s1, s2 string) string {
	opp := map[byte]byte{'N': 'S', 'S': 'N', 'E': 'W', 'W': 'E'}
	type pair struct{ i, j int }
	vis := make(map[pair]struct{})
	q := make([]pair, 0)
	start := pair{0, 0}
	q = append(q, start)
	vis[start] = struct{}{}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.i == n-1 && cur.j == n-1 {
			return "YES"
		}
		dirs := make(map[byte]struct{})
		if cur.i < n-1 {
			dirs[s1[cur.i]] = struct{}{}
		}
		if cur.i > 0 {
			dirs[opp[s1[cur.i-1]]] = struct{}{}
		}
		if cur.j < n-1 {
			dirs[s2[cur.j]] = struct{}{}
		}
		if cur.j > 0 {
			dirs[opp[s2[cur.j-1]]] = struct{}{}
		}
		for d := range dirs {
			ni, nj := cur.i, cur.j
			if cur.i < n-1 && s1[cur.i] == d {
				ni = cur.i + 1
			} else if cur.i > 0 && opp[s1[cur.i-1]] == d {
				ni = cur.i - 1
			}
			if cur.j < n-1 && s2[cur.j] == d {
				nj = cur.j + 1
			} else if cur.j > 0 && opp[s2[cur.j-1]] == d {
				nj = cur.j - 1
			}
			np := pair{ni, nj}
			if _, ok := vis[np]; !ok {
				vis[np] = struct{}{}
				q = append(q, np)
			}
		}
	}
	return "NO"
}

func parseLine(line string) (int, string, string, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return 0, "", "", fmt.Errorf("expected 3 fields got %d", len(parts))
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", "", err
	}
	if len(parts[1]) != n-1 || len(parts[2]) != n-1 {
		return 0, "", "", fmt.Errorf("invalid string lengths for n=%d", n)
	}
	return n, parts[1], parts[2], nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesC))
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, s1, s2, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		expected := solve(n, s1, s2)

		var input strings.Builder
		input.WriteString(strconv.Itoa(n))
		input.WriteByte('\n')
		input.WriteString(s1)
		input.WriteByte('\n')
		input.WriteString(s2)
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
