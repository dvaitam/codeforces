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
   fmt.Fscan(reader, &n, &m)
   grid := make([]string, n)
   var si, sj int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
       if idx := indexRune(grid[i], 'S'); idx >= 0 {
           si, sj = i, idx
       }
   }

   // visited and offsets flattened
   NM := n * m
   visited := make([]bool, NM)
   offX := make([]int, NM)
   offY := make([]int, NM)
   // queue of positions
   queue := make([]int, 0, NM)
   start := si*m + sj
   visited[start] = true
   queue = append(queue, start)

   // directions
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for qi := 0; qi < len(queue); qi++ {
       v := queue[qi]
       i := v / m
       j := v % m
       bx := offX[v]
       by := offY[v]
       for _, d := range dirs {
           ni := i + d[0]
           nj := j + d[1]
           nx := bx
           ny := by
           if ni < 0 {
               ni += n
               nx--
           } else if ni >= n {
               ni -= n
               nx++
           }
           if nj < 0 {
               nj += m
               ny--
           } else if nj >= m {
               nj -= m
               ny++
           }
           if grid[ni][nj] == '#' {
               continue
           }
           u := ni*m + nj
           if !visited[u] {
               visited[u] = true
               offX[u] = nx
               offY[u] = ny
               queue = append(queue, u)
           } else if offX[u] != nx || offY[u] != ny {
               fmt.Fprint(writer, "Yes")
               return
           }
       }
   }
   fmt.Fprint(writer, "No")
}

// indexRune finds the first index of r in s, or -1 if not present
func indexRune(s string, r rune) int {
   for i, c := range s {
       if c == r {
           return i
       }
   }
   return -1
}
