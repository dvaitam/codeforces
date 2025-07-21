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
   parent := make([]int, n+1)
   addTime := make([]int, n+1)
   // packets
   s := make([]int, 1, m+1)
   packetTime := make([]int, 1, m+1)
   type Query struct { x, pid, idx int }
   var queries []Query
   ansCount := 0
   // read events
   for t := 1; t <= m; t++ {
       var tp int
       fmt.Fscan(reader, &tp)
       if tp == 1 {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           parent[x] = y
           addTime[x] = t
       } else if tp == 2 {
           var x int
           fmt.Fscan(reader, &x)
           s = append(s, x)
           packetTime = append(packetTime, t)
       } else if tp == 3 {
           var x, pid int
           fmt.Fscan(reader, &x, &pid)
           queries = append(queries, Query{x: x, pid: pid, idx: ansCount})
           ansCount++
       }
   }
   // build children
   children := make([][]int, n+1)
   for v := 1; v <= n; v++ {
       p := parent[v]
       if p != 0 {
           children[p] = append(children[p], v)
       }
   }
   // DFS for tin/tout and depth
   depth := make([]int, n+1)
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   timer := 0
   type stackEntry struct{ v, idx, state int }
   var stack []stackEntry
   for v := 1; v <= n; v++ {
       if parent[v] == 0 {
           stack = append(stack, stackEntry{v: v, idx: 0, state: 0})
           for len(stack) > 0 {
               e := &stack[len(stack)-1]
               if e.state == 0 {
                   timer++
                   tin[e.v] = timer
                   e.state = 1
               }
               if e.idx < len(children[e.v]) {
                   c := children[e.v][e.idx]
                   e.idx++
                   depth[c] = depth[e.v] + 1
                   stack = append(stack, stackEntry{v: c, idx: 0, state: 0})
               } else {
                   tout[e.v] = timer
                   stack = stack[:len(stack)-1]
               }
           }
       }
   }
   // binary lifting
   const LOG = 18
   up := make([][]int, LOG)
   timeUp := make([][]int, LOG)
   for j := 0; j < LOG; j++ {
       up[j] = make([]int, n+1)
       timeUp[j] = make([]int, n+1)
   }
   for v := 1; v <= n; v++ {
       up[0][v] = parent[v]
       timeUp[0][v] = addTime[v]
   }
   for j := 1; j < LOG; j++ {
       for v := 1; v <= n; v++ {
           u := up[j-1][v]
           up[j][v] = up[j-1][u]
           // max time on path of length 2^j
           t1 := timeUp[j-1][v]
           t2 := timeUp[j-1][u]
           if t2 > t1 {
               timeUp[j][v] = t2
           } else {
               timeUp[j][v] = t1
           }
       }
   }
   // answer queries
   answers := make([]string, ansCount)
   for qi, q := range queries {
       x := q.x
       pid := q.pid
       u := s[pid]
       t0 := packetTime[pid]
       // check ancestor in final tree
       if !(tin[x] <= tin[u] && tout[u] <= tout[x]) {
           answers[qi] = "NO"
           continue
       }
       // climb from u to x, track max addTime
       diff := depth[u] - depth[x]
       cur := u
       maxT := 0
       for j := 0; j < LOG; j++ {
           if diff&(1<<j) != 0 {
               if timeUp[j][cur] > maxT {
                   maxT = timeUp[j][cur]
               }
               cur = up[j][cur]
           }
       }
       if maxT <= t0 {
           answers[qi] = "YES"
       } else {
           answers[qi] = "NO"
       }
   }
   // output
   for i := 0; i < ansCount; i++ {
       writer.WriteString(answers[i])
       writer.WriteByte('\n')
   }
}
