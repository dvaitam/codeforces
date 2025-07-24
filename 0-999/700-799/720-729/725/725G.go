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
   parent[0] = -1
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &parent[i])
   }
   depth := make([]int, n+1)
   for i := 1; i <= n; i++ {
       depth[i] = depth[parent[i]] + 1
   }
   // f[u]: time until which node u is busy waiting (finish time)
   f := make([]int64, n+1)
   // lastAcceptTime[u], lastAcceptX[u] for collision handling
   lastAcceptTime := make([]int64, n+1)
   lastAcceptX := make([]int, n+1)
   // answers
   ans := make([]int64, m)
   for i := 0; i < m; i++ {
       var x int
       var t int64
       fmt.Fscan(reader, &x, &t)
       // if initiator busy, reject immediately
       if f[x] > t {
           ans[i] = t
           continue
       }
       // traverse up
       var pathUs []int
       var pathTs []int64
       u := x
       T := t
       // include x itself
       pathUs = append(pathUs, u)
       pathTs = append(pathTs, T)
       // move up until block or root
       blocked := false
       for u != 0 {
           // move to parent
           u = parent[u]
           T++
           if u == 0 {
               break
           }
           // collision: same time and smaller x waited here
           if lastAcceptTime[u] == T && lastAcceptX[u] < x {
               blocked = true
               break
           }
           // busy
           if f[u] > T {
               blocked = true
               break
           }
           // accept at u
           lastAcceptTime[u] = T
           lastAcceptX[u] = x
           pathUs = append(pathUs, u)
           pathTs = append(pathTs, T)
       }
       var tEnd int64
       // if reached root (u==0)
       if u == 0 {
           // arrival to Bob at time T, reply travels depth[x] down
           tEnd = T + int64(depth[x])
       } else {
           // blocked at u at time T, reply travels down from u to x: depth[x]-depth[u]
           tEnd = T + int64(depth[x]-depth[u])
       }
       ans[i] = tEnd
       // update finish times for waiting nodes
       for idx, uu := range pathUs {
           // finish time when reply reaches uu: tEnd - (depth[x]-depth[uu])
           finish := tEnd - int64(depth[x]-depth[uu])
           if f[uu] < finish {
               f[uu] = finish
           }
       }
   }
   // output answers
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
