package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var good string
	if _, err := fmt.Fscan(reader, &good); err != nil {
		return
	}
	goodSet := make([]bool, 26)
	for _, ch := range good {
		goodSet[ch-'a'] = true
	}

	var pattern string
	fmt.Fscan(reader, &pattern)
	var n int
	fmt.Fscan(reader, &n)

	star := -1
	for i, ch := range pattern {
		if ch == '*' {
			star = i
			break
		}
	}

	for ; n > 0; n-- {
		var q string
		fmt.Fscan(reader, &q)
		if star == -1 {
			if len(q) != len(pattern) {
				fmt.Fprintln(writer, "NO")
				continue
			}
			ok := true
			for i := 0; i < len(pattern); i++ {
				pc := pattern[i]
				qc := q[i]
				switch pc {
				case '?':
					if !goodSet[qc-'a'] {
						ok = false
					}
				default:
					if pc != qc {
						ok = false
					}
				}
				if !ok {
					break
				}
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		} else {
			prefixLen := star
			suffixLen := len(pattern) - star - 1
			if len(q) < prefixLen+suffixLen {
				fmt.Fprintln(writer, "NO")
				continue
			}
			ok := true
			// prefix
			for i := 0; i < prefixLen; i++ {
				pc := pattern[i]
				qc := q[i]
				switch pc {
				case '?':
					if !goodSet[qc-'a'] {
						ok = false
					}
				default:
					if pc != qc {
						ok = false
					}
				}
				if !ok {
					break
				}
			}
			// suffix
			if ok {
				for i := 0; i < suffixLen; i++ {
					pc := pattern[len(pattern)-1-i]
					qc := q[len(q)-1-i]
					switch pc {
					case '?':
						if !goodSet[qc-'a'] {
							ok = false
						}
					default:
						if pc != qc {
							ok = false
						}
					}
					if !ok {
						break
					}
				}
			}
			// middle for star
			if ok {
				start := prefixLen
				end := len(q) - suffixLen
				for i := start; i < end; i++ {
					qc := q[i]
					if goodSet[qc-'a'] {
						ok = false
						break
					}
				}
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
