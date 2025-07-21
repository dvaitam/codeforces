package main

import (
   "bufio"
   "fmt"
   "os"
)

type state struct {
   i, j, mask int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &grid[i])
   }
   // find start, treasures count, bombs
   var si, sj int
   maxDigit := '0'
   bombCount := 0
   for i := 0; i < n; i++ {
       for j, c := range grid[i] {
           if c == 'S' {
               si, sj = i, j
           } else if c >= '1' && c <= '8' {
               if c > maxDigit {
                   maxDigit = c
               }
           } else if c == 'B' {
               bombCount++
           }
       }
   }
   t := int(maxDigit - '0')
   // read treasure values
   treasureVals := make([]int, t)
   for i := 0; i < t; i++ {
       fmt.Fscan(in, &treasureVals[i])
   }
   // build objects: treasures then bombs
   o := t + bombCount
   xs := make([]int, o)
   ys := make([]int, o)
   vals := make([]int, o)
   bombMask := 0
   // treasures
   for i := 0; i < n; i++ {
       for j, c := range grid[i] {
           if c >= '1' && c <= '8' {
               tid := int(c - '1')
               xs[tid] = j
               ys[tid] = i
               vals[tid] = treasureVals[tid]
           }
       }
   }
   // bombs
   bi := t
   for i := 0; i < n; i++ {
       for j, c := range grid[i] {
           if c == 'B' {
               xs[bi] = j
               ys[bi] = i
               vals[bi] = 0
               bombMask |= 1 << bi
               bi++
           }
       }
   }
   maxMask := 1 << o
   // distances
   const INF = -1
   dist := make([][][]int, n)
   for i := 0; i < n; i++ {
       dist[i] = make([][]int, m)
       for j := 0; j < m; j++ {
           dist[i][j] = make([]int, maxMask)
           for k := range dist[i][j] {
               dist[i][j][k] = INF
           }
       }
   }
   // BFS
   q := make([]state, 0, n*m*4)
   dist[si][sj][0] = 0
   q = append(q, state{si, sj, 0})
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for head := 0; head < len(q); head++ {
       cur := q[head]
       d := dist[cur.i][cur.j][cur.mask]
       for _, dir := range dirs {
           ni := cur.i + dir[0]
           nj := cur.j + dir[1]
           if ni < 0 || ni >= n || nj < 0 || nj >= m {
               continue
           }
           c := grid[ni][nj]
           if c == '#' || c == 'B' || (c >= '1' && c <= '8') {
               continue
           }
           newMask := cur.mask
           // check crossings
           if cur.j == nj {
               // vertical move
               x := cur.j
               y1, y2 := cur.i, ni
               if y1 > y2 {
                   y1, y2 = y2, y1
               }
               for k := 0; k < o; k++ {
                   if xs[k] < x && ys[k] >= y1 && ys[k] < y2 {
                       newMask ^= 1 << k
                   }
               }
           }
           if dist[ni][nj][newMask] == INF {
               dist[ni][nj][newMask] = d + 1
               q = append(q, state{ni, nj, newMask})
           }
       }
   }
   // compute answer
   ans := 0
   for mask := 0; mask < maxMask; mask++ {
       k := dist[si][sj][mask]
       if k <= 0 {
           continue
       }
       if mask&bombMask != 0 {
           continue
       }
       sum := 0
       for i := 0; i < o; i++ {
           if mask&(1<<i) != 0 {
               sum += vals[i]
           }
       }
       profit := sum - k
       if profit > ans {
           ans = profit
       }
   }
   fmt.Fprintln(out, ans)
}
