package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, m, k1 int
   if _, err := fmt.Fscan(reader, &n, &m, &k1); err != nil {
       return
   }
   total := n * m / 2
   k2 := total - k1
   origN, origM := n, m
   // allocate grid
   grid := make([][]rune, origN)
   for i := range grid {
       grid[i] = make([]rune, origM)
   }
   // handle odd row
   if n%2 == 1 {
       now := m / 2
       if k1 < now {
           fmt.Fprintln(writer, "NO")
           return
       }
       // fill last row with horizontal dominos
       row := origN - 1
       for j := 0; j < origM; j += 4 {
           if j < origM {
               grid[row][j] = 'a'
           }
           if j+1 < origM {
               grid[row][j+1] = 'a'
           }
           if j+2 < origM {
               grid[row][j+2] = 'b'
           }
           if j+3 < origM {
               grid[row][j+3] = 'b'
           }
       }
       k1 -= now
       n--
   } else if m%2 == 1 {
       // handle odd column
       now := n / 2
       if k2 < now {
           fmt.Fprintln(writer, "NO")
           return
       }
       // fill last column with vertical dominos
       col := origM - 1
       for i := 0; i < origN; i += 4 {
           if i < origN {
               grid[i][col] = 'a'
           }
           if i+1 < origN {
               grid[i+1][col] = 'a'
           }
           if i+2 < origN {
               grid[i+2][col] = 'b'
           }
           if i+3 < origN {
               grid[i+3][col] = 'b'
           }
       }
       k2 -= now
       m--
   }
   // parity check
   if k1%2 != 0 || k2%2 != 0 {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   // fill 2x2 blocks
   c := 'c'
   rows := n
   cols := m
   for i := 0; i < rows; i += 2 {
       for j := 0; j < cols; j += 2 {
           // choose orientation
           if k1 > 0 {
               k1 -= 2
               // horizontal dominos
               grid[i][j] = c
               grid[i][j+1] = c
               grid[i+1][j] = c + 1
               grid[i+1][j+1] = c + 1
           } else {
               // vertical dominos
               grid[i][j] = c
               grid[i+1][j] = c
               grid[i][j+1] = c + 1
               grid[i+1][j+1] = c + 1
           }
           // advance and fix color if adjacent match
           // wrap c
           wrap := func() {
               if int(c-'a') == 26 {
                   c = 'c'
               }
           }
           c += 2; wrap()
           for {
               conflict := false
               if i > 0 {
                   if grid[i][j] == grid[i-1][j] || grid[i][j+1] == grid[i-1][j+1] {
                       conflict = true
                   }
               }
               if j > 0 {
                   if grid[i][j] == grid[i][j-1] || grid[i+1][j] == grid[i+1][j-1] {
                       conflict = true
                   }
               }
               if !conflict {
                   break
               }
               // reassign block
               grid[i][j] = c
               grid[i][j+1] = c
               grid[i+1][j] = c + 1
               grid[i+1][j+1] = c + 1
               c += 2; wrap()
           }
       }
   }
   // output grid
   for i := 0; i < origN; i++ {
       for j := 0; j < origM; j++ {
           writer.WriteRune(grid[i][j])
       }
       writer.WriteByte('\n')
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for tc := 0; tc < T; tc++ {
       solve(reader, writer)
   }
}
