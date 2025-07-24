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
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)
   pos := make([][]int, 8)
   for i, ch := range s {
       idx := int(ch - 'a')
       pos[idx] = append(pos[idx], i)
   }
   const INF = 1e9
   n8 := n
   // BFS helper
   dist := make([]int, n8)
   usedLetter := make([]bool, 8)
   var q []int
   bestD := 0
   var bestCnt int64
   for c := 0; c < 8; c++ {
       // init dist
       for i := 0; i < n8; i++ {
           dist[i] = int(INF)
       }
       for i := range usedLetter {
           usedLetter[i] = false
       }
       q = q[:0]
       // multi-source from pos[c]
       for _, i := range pos[c] {
           dist[i] = 0
           q = append(q, i)
       }
       usedLetter[c] = true
       // BFS
       for head := 0; head < len(q); head++ {
           u := q[head]
           d := dist[u]
           // neighbors
           if u > 0 && dist[u-1] > d+1 {
               dist[u-1] = d + 1
               q = append(q, u-1)
           }
           if u+1 < n8 && dist[u+1] > d+1 {
               dist[u+1] = d + 1
               q = append(q, u+1)
           }
           lc := int(s[u] - 'a')
           if !usedLetter[lc] {
               usedLetter[lc] = true
               for _, v := range pos[lc] {
                   if dist[v] > d+1 {
                       dist[v] = d + 1
                       q = append(q, v)
                   }
               }
           }
       }
       // compute max dist and count
       mx := 0
       var cnt int64
       for i := 0; i < n8; i++ {
           if dist[i] > mx {
               mx = dist[i]
               cnt = 1
           } else if dist[i] == mx {
               cnt++
           }
       }
       var D int
       var pairs int64
       sx := len(pos[c])
       if sx >= 2 {
           D = mx + 1
           pairs = cnt * int64(sx-1)
       } else {
           D = mx
           pairs = cnt
       }
       if D > bestD {
           bestD = D
           bestCnt = pairs
       } else if D == bestD {
           bestCnt += pairs
       }
   }
   fmt.Fprintln(writer, bestD, bestCnt)
}
