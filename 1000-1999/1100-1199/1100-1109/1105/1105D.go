package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, p int
   fmt.Fscan(reader, &n, &m, &p)
   steps := make([]int, p+1)
   for i := 1; i <= p; i++ {
       fmt.Fscan(reader, &steps[i])
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }
   // initialize frontiers
   frontier := make([][]int, p+1)
   locked := make([]bool, p+1)
   lockedCount := 0
   for r := 0; r < n; r++ {
       for c := 0; c < m; c++ {
           ch := grid[r][c]
           if ch >= '1' && ch <= byte('0'+p) {
               pid := int(ch - '0')
               frontier[pid] = append(frontier[pid], r*m+c)
           }
       }
   }
   // simulate turns
   for lockedCount < p {
       progress := false
       for pid := 1; pid <= p; pid++ {
           if locked[pid] {
               continue
           }
           // perform up to steps[pid] BFS layers
           for s := 0; s < steps[pid]; s++ {
               cur := frontier[pid]
               if len(cur) == 0 {
                   locked[pid] = true
                   lockedCount++
                   break
               }
               progress = true
               next := make([]int, 0)
               for _, code := range cur {
                   r := code / m
                   c := code % m
                   // expand neighbors
                   if r > 0 && grid[r-1][c] == '.' {
                       grid[r-1][c] = byte('0' + pid)
                       next = append(next, (r-1)*m+c)
                   }
                   if r+1 < n && grid[r+1][c] == '.' {
                       grid[r+1][c] = byte('0' + pid)
                       next = append(next, (r+1)*m+c)
                   }
                   if c > 0 && grid[r][c-1] == '.' {
                       grid[r][c-1] = byte('0' + pid)
                       next = append(next, r*m+(c-1))
                   }
                   if c+1 < m && grid[r][c+1] == '.' {
                       grid[r][c+1] = byte('0' + pid)
                       next = append(next, r*m+(c+1))
                   }
               }
               frontier[pid] = next
           }
       }
       if !progress {
           break
       }
   }
   // count territories
   ans := make([]int, p+1)
   for r := 0; r < n; r++ {
       for c := 0; c < m; c++ {
           ch := grid[r][c]
           if ch >= '1' && ch <= byte('0'+p) {
               ans[ch-'0']++
           }
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 1; i <= p; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", ans[i])
   }
   writer.WriteByte('\n')
}
