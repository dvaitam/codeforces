package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   x := make([]int64, n)
   y := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   k := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &k[i])
   }

   dis := make([]int64, n)
   pre := make([]int, n)
   vis := make([]bool, n)
   for i := 0; i < n; i++ {
       dis[i] = c[i]
       pre[i] = -1
   }

   var total int64
   var powerStations []int
   var connections [][2]int

   for iter := 0; iter < n; iter++ {
       v := -1
       var useStation bool
       for j := 0; j < n; j++ {
           if !vis[j] && (v == -1 || dis[j] < dis[v]) {
               v = j
               useStation = dis[v] == c[v]
           }
       }
       vis[v] = true
       total += dis[v]
       if useStation {
           powerStations = append(powerStations, v+1)
       } else {
           connections = append(connections, [2]int{pre[v] + 1, v + 1})
       }
       for j := 0; j < n; j++ {
           if vis[j] {
               continue
           }
           cost := (k[v] + k[j]) * (abs64(x[v]-x[j]) + abs64(y[v]-y[j]))
           if cost < dis[j] {
               dis[j] = cost
               pre[j] = v
           }
       }
   }

   fmt.Fprintln(writer, total)
   fmt.Fprintln(writer, len(powerStations))
   for i, v := range powerStations {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
   fmt.Fprintln(writer, len(connections))
   for _, p := range connections {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
