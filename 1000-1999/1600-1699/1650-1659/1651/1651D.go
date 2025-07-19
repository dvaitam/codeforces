package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Solve 1651D: for each point find nearest empty integer grid point by Manhattan distance
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   type Node struct{ x, y, id int }
   nodes := make([]Node, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nodes[i].x, &nodes[i].y)
       nodes[i].id = i
   }
   const maxX = 200000
   mi := make([]int, maxX+2)
   ma := make([]int, maxX+2)
   asw := make([]int, n)
   ansX := make([]int, n)
   ansY := make([]int, n)
   const INF = 1000000000

   // sort by y asc, x asc
   sort.Slice(nodes, func(i, j int) bool {
       if nodes[i].y != nodes[j].y {
           return nodes[i].y < nodes[j].y
       }
       return nodes[i].x < nodes[j].x
   })
   var ASW, xm, X, Y int
   // sweep: horizontal left, vertical down
   for i := 0; i < n; i++ {
       if i == 0 || nodes[i].y != nodes[i-1].y {
           ASW = INF
           xm = nodes[i].x - 1
       } else {
           ASW += nodes[i].x - nodes[i-1].x
           if nodes[i-1].x+1 != nodes[i].x {
               xm = nodes[i].x - 1
           }
       }
       id := nodes[i].id
       // candidate left
       asw[id] = nodes[i].x - xm
       ansX[id], ansY[id] = xm, nodes[i].y
       // candidate down
       if ma[nodes[i].x]+1 != nodes[i].y {
           mi[nodes[i].x] = nodes[i].y - 1
       }
       ma[nodes[i].x] = nodes[i].y
       if nodes[i].y-mi[nodes[i].x] < ASW {
           ASW = nodes[i].y - mi[nodes[i].x]
           X, Y = nodes[i].x, mi[nodes[i].x]
       }
       if ASW < asw[id] {
           asw[id], ansX[id], ansY[id] = ASW, X, Y
       }
   }
   // reset
   for i := 0; i <= maxX; i++ {
       ma[i], mi[i] = 0, 0
   }
   // sweep: horizontal right, vertical up
   for i := n - 1; i >= 0; i-- {
       if i == n-1 || nodes[i].y != nodes[i+1].y {
           ASW = INF
           xm = nodes[i].x + 1
       } else {
           ASW += nodes[i+1].x - nodes[i].x
           if nodes[i+1].x-1 != nodes[i].x {
               xm = nodes[i].x + 1
           }
       }
       id := nodes[i].id
       // candidate right
       if xm-nodes[i].x < asw[id] {
           asw[id], ansX[id], ansY[id] = xm-nodes[i].x, xm, nodes[i].y
       }
       // candidate up
       if ma[nodes[i].x]-1 != nodes[i].y {
           mi[nodes[i].x] = nodes[i].y + 1
       }
       ma[nodes[i].x] = nodes[i].y
       if mi[nodes[i].x]-nodes[i].y < ASW {
           ASW = mi[nodes[i].x] - nodes[i].y
           X, Y = nodes[i].x, mi[nodes[i].x]
       }
       if ASW < asw[id] {
           asw[id], ansX[id], ansY[id] = ASW, X, Y
       }
   }
   // sort by y asc, x desc
   sort.Slice(nodes, func(i, j int) bool {
       if nodes[i].y != nodes[j].y {
           return nodes[i].y < nodes[j].y
       }
       return nodes[i].x > nodes[j].x
   })
   for i := 0; i <= maxX; i++ {
       ma[i], mi[i] = 0, 0
   }
   // sweep: horizontal right (in reversed x order), vertical down
   for i := 0; i < n; i++ {
       if i == 0 || nodes[i].y != nodes[i-1].y {
           ASW = INF
           xm = nodes[i].x + 1
       } else {
           ASW += nodes[i-1].x - nodes[i].x
           if nodes[i-1].x-1 != nodes[i].x {
               xm = nodes[i].x + 1
           }
       }
       id := nodes[i].id
       if xm-nodes[i].x < asw[id] {
           asw[id], ansX[id], ansY[id] = xm-nodes[i].x, xm, nodes[i].y
       }
       if ma[nodes[i].x]+1 != nodes[i].y {
           mi[nodes[i].x] = nodes[i].y - 1
       }
       ma[nodes[i].x] = nodes[i].y
       if nodes[i].y-mi[nodes[i].x] < ASW {
           ASW = nodes[i].y - mi[nodes[i].x]
           X, Y = nodes[i].x, mi[nodes[i].x]
       }
       if ASW < asw[id] {
           asw[id], ansX[id], ansY[id] = ASW, X, Y
       }
   }
   for i := 0; i <= maxX; i++ {
       ma[i], mi[i] = 0, 0
   }
   // sweep: horizontal left, vertical up
   for i := n - 1; i >= 0; i-- {
       if i == n-1 || nodes[i].y != nodes[i+1].y {
           ASW = INF
           xm = nodes[i].x - 1
       } else {
           ASW += nodes[i].x - nodes[i+1].x
           if nodes[i+1].x+1 != nodes[i].x {
               xm = nodes[i].x - 1
           }
       }
       id := nodes[i].id
       if nodes[i].x-xm < asw[id] {
           asw[id], ansX[id], ansY[id] = nodes[i].x-xm, xm, nodes[i].y
       }
       if ma[nodes[i].x]-1 != nodes[i].y {
           mi[nodes[i].x] = nodes[i].y + 1
       }
       ma[nodes[i].x] = nodes[i].y
       if mi[nodes[i].x]-nodes[i].y < ASW {
           ASW = mi[nodes[i].x] - nodes[i].y
           X, Y = nodes[i].x, mi[nodes[i].x]
       }
       if ASW < asw[id] {
           asw[id], ansX[id], ansY[id] = ASW, X, Y
       }
   }
   // output in original order
   for i := 0; i < n; i++ {
       fmt.Fprintln(writer, ansX[i], ansY[i])
   }
}
