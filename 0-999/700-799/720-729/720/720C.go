package main

import (
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	T := readInt(reader)
	for tc := 0; tc < T; tc++ {
		n := readInt(reader)
		m := readInt(reader)
		y := readInt(reader)
		ok, mp := build(n, m, y)
		if ok {
			for i := 1; i <= n; i++ {
				writer.Write(mp[i][1 : m+1])
				writer.WriteByte('\n')
			}
		} else {
			ok2, mp2 := build(m, n, y)
			if !ok2 {
				writer.WriteString("-1\n")
			} else {
				for i := 1; i <= n; i++ {
					// transpose
					for j := 1; j <= m; j++ {
						writer.WriteByte(mp2[j][i])
					}
					writer.WriteByte('\n')
				}
			}
		}
		if tc < T-1 {
			writer.WriteByte('\n')
		}
	}
}

// readInt reads next integer from bufio.Reader
func readInt(r *bufio.Reader) int {
	var x int
	b, _ := r.ReadByte()
	for b < '0' || b > '9' {
		b, _ = r.ReadByte()
	}
	for b >= '0' && b <= '9' {
		x = x*10 + int(b-'0')
		b, _ = r.ReadByte()
	}
	return x
}

// build attempts to place t stars in an n x m grid; returns grid or false
func build(n, m, t0 int) (bool, [][]byte) {
	t := t0
	if t > (n-1)*(m-1)*4 {
		return false, nil
	}
	// use 1-based indexing with padding
	mp := make([][]byte, n+2)
	for i := range mp {
		mp[i] = make([]byte, m+2)
		for j := range mp[i] {
			mp[i][j] = '.'
		}
	}
	mp[1][1], mp[2][1] = '*', '*'
	var i, j int
	for i = 2; i <= n; i++ {
		if i > 2 {
			t--
			mp[i][1] = '*'
		}
		j = 2
		if t < 4 {
			break
		}
		for j = 2; j <= m && t >= 4; j++ {
			mp[i-1][j-1], mp[i-1][j], mp[i][j-1], mp[i][j] = '*', '*', '*', '*'
			if i > 2 && j == m {
				t -= 3
			} else {
				t -= 4
			}
		}
		if t < 4 {
			break
		}
	}
	if i == 2 && j <= m {
		switch t {
		case 1:
			mp[1][j] = '*'
		case 2:
			if j == 2 {
				mp[2][2], mp[3][1] = '*', '*'
			} else {
				mp[1][j], mp[3][1] = '*', '*'
			}
		case 3:
			if j == 2 {
				mp[3][3], mp[2][3] = '*', '*'
			}
			mp[1][j], mp[3][j-1], mp[3][j] = '*', '*', '*'
		}
	} else {
		switch t {
		case 1:
			if j < m {
				mp[i][m] = '*'
			} else if i < n {
				mp[i+1][1] = '*'
			} else {
				return false, nil
			}
		case 2:
			if j < m-1 {
				mp[i][m-1] = '*'
			} else if i < n && j > 3 {
				mp[i+1][2] = '*'
			} else if i < n && j <= m {
				mp[i+1][j], mp[i+1][j-1] = '*', '*'
			} else {
				return false, nil
			}
		case 3:
			if j < m-2 {
				mp[i][m], mp[i][m-2] = '*', '*'
			} else if j == m {
				mp[i][j] = '*'
			} else if i < n {
				if j < m {
					mp[i][m] = '*'
					mp[i+1][m-1], mp[i+1][j], mp[i+1][j-1] = '*', '*', '*'
					if j == 2 {
						mp[i+1][m] = '*'
					}
				} else if m >= 4 {
					mp[i+1][1], mp[i+1][m-1] = '*', '*'
					if j == m {
						mp[i+1][m] = '*'
					}
				} else {
					return false, nil
				}
			} else {
				return false, nil
			}
		}
	}
	return true, mp
}
