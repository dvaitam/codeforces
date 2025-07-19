package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node holds depth (d), move range (m), and relay range (r)
type Node struct {
   d, m, r int
}

var (
   a  [4]Node
   ff [47][47][47][8][8]bool
   an int
)

// dfs explores reachable states, tracking visited in ff and updating an
func dfs(n1, n2 int) {
   d1, d2, d3 := a[1].d, a[2].d, a[3].d
   if ff[d1][d2][d3][n1][n2] {
       return
   }
   // update max depth
   for i := 1; i <= 3; i++ {
       if a[i].d <= 42 && a[i].d > an {
           an = a[i].d
       }
   }
   ff[d1][d2][d3][n1][n2] = true
   // minimum current depth
   mm := a[1].d
   if a[2].d < mm {
       mm = a[2].d
   }
   if a[3].d < mm {
       mm = a[3].d
   }
   // identify relayed nodes
   var ju [4]int
   for i := 1; i <= 3; i++ {
       if a[i].d > 42 {
           idx := a[i].d - 42
           if idx >= 1 && idx <= 3 {
               ju[idx] = i
           }
       }
   }
   // move self
   for i := 1; i <= 3; i++ {
       if a[i].d <= 42 && (n1&(1<<(i-1))) == 0 && ju[i] == 0 {
           old := a[i].d
           from := old - a[i].m
           if from < mm {
               from = mm
           }
           to := old + a[i].m
           for j := from; j <= to; j++ {
               a[i].d = j
               dfs(n1|(1<<(i-1)), n2)
           }
           a[i].d = old
       }
   }
   // initiate relay
   for i := 1; i <= 3; i++ {
       if a[i].d <= 42 && (n2&(1<<(i-1))) == 0 && ju[i] == 0 {
           for j := 1; j <= 3; j++ {
               if abs(a[j].d-a[i].d) == 1 {
                   oldj := a[j].d
                   a[j].d = 42 + i
                   dfs(n1, n2|(1<<(i-1)))
                   a[j].d = oldj
               }
           }
       }
   }
   // relay moves
   for i := 1; i <= 3; i++ {
       if a[i].d <= 42 && ju[i] != 0 {
           y := ju[i]
           oldy := a[y].d
           from := a[i].d - a[i].r
           if from < mm {
               from = mm
           }
           to := a[i].d + a[i].r
           for j := from; j <= to; j++ {
               a[y].d = j
               dfs(n1, n2)
           }
           a[y].d = oldy
       }
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       // read initial node
       if _, err := fmt.Fscan(reader, &a[1].d, &a[1].m, &a[1].r); err != nil {
           break
       }
       // read remaining
       for i := 2; i <= 3; i++ {
           fmt.Fscan(reader, &a[i].d, &a[i].m, &a[i].r)
       }
       // reset visited and answer
       ff = [47][47][47][8][8]bool{}
       an = 0
       dfs(0, 0)
       fmt.Fprintln(writer, an)
   }
}
