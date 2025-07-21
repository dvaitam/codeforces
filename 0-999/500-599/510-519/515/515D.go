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
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   // degrees of empty cells
   deg := make([][]int, n)
   for i := 0; i < n; i++ {
       deg[i] = make([]int, m)
   }
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != '.' {
               continue
           }
           cnt := 0
           for _, d := range dirs {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
                   cnt++
               }
           }
           deg[i][j] = cnt
       }
   }
   // queue of forced cells
   type pair struct{ x, y int }
   queue := make([]pair, 0, n*m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' && deg[i][j] == 1 {
               queue = append(queue, pair{i, j})
           }
       }
   }
   head := 0
   for head < len(queue) {
       p := queue[head]; head++
       i, j := p.x, p.y
       if grid[i][j] != '.' {
           continue
       }
       // find the only neighbor
       placed := false
       for _, d := range dirs {
           ni, nj := i+d[0], j+d[1]
           if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
               // place domino between (i,j) and (ni,nj)
               if ni == i {
                   // horizontal
                   if nj == j+1 {
                       grid[i][j], grid[ni][nj] = '<', '>'
                   } else {
                       grid[ni][nj], grid[i][j] = '<', '>'
                   }
               } else {
                   // vertical
                   if ni == i+1 {
                       grid[i][j], grid[ni][nj] = '^', 'v'
                   } else {
                       grid[ni][nj], grid[i][j] = '^', 'v'
                   }
               }
               // update neighbors' degrees
               for _, d2 := range dirs {
                   xi, yj := i+d2[0], j+d2[1]
                   if xi >= 0 && xi < n && yj >= 0 && yj < m && grid[xi][yj] == '.' {
                       deg[xi][yj]--
                       if deg[xi][yj] == 1 {
                           queue = append(queue, pair{xi, yj})
                       }
                   }
                   xi, yj = ni+d2[0], nj+d2[1]
                   if xi >= 0 && xi < n && yj >= 0 && yj < m && grid[xi][yj] == '.' {
                       deg[xi][yj]--
                       if deg[xi][yj] == 1 {
                           queue = append(queue, pair{xi, yj})
                       }
                   }
               }
               placed = true
               break
           }
       }
       if !placed {
           fmt.Fprintln(writer, "Not unique")
           return
       }
   }
   // check if any empty remains
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' {
               fmt.Fprintln(writer, "Not unique")
               return
           }
       }
   }
   // output solution
   for i := 0; i < n; i++ {
       writer.Write(grid[i])
       writer.WriteByte('\n')
   }
}
