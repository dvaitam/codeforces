package main

import (
   "bufio"
   "fmt"
   "os"
)

// getInv computes the inverse as defined in the problem.
func getInv(v int) int {
   // compute floor(log2(v)) by finding highest bit
   lg := 0
   for (1 << (lg + 1)) <= v {
       lg++
   }
   inv := (1 << (lg + 1)) - v
   if inv == v {
       return 0
   }
   return inv
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // build paths from root (0) to each node
   aP := make([][]int, n)
   for i := 0; i < n; i++ {
       path := []int{a[i]}
       for path[len(path)-1] != 0 {
           path = append(path, getInv(path[len(path)-1]))
       }
       // reverse path
       for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
           path[l], path[r] = path[r], path[l]
       }
       aP[i] = path
   }

   // find farthest from node 0
   from, maxDist, maxV := 0, -1, -1
   for i := 0; i < n; i++ {
       if i == from {
           continue
       }
       // compute shared prefix length
       pf := 0
       p1, p2 := aP[from], aP[i]
       m := len(p1)
       if len(p2) < m {
           m = len(p2)
       }
       for pf < m && p1[pf] == p2[pf] {
           pf++
       }
       dist := len(p1) + len(p2) - 2*pf
       if dist > maxDist {
           maxDist = dist
           maxV = i
       }
   }
   // run again from the farthest
   from = maxV
   maxV, maxDist = -1, -1
   for i := 0; i < n; i++ {
       if i == from {
           continue
       }
       pf := 0
       p1, p2 := aP[from], aP[i]
       m := len(p1)
       if len(p2) < m {
           m = len(p2)
       }
       for pf < m && p1[pf] == p2[pf] {
           pf++
       }
       dist := len(p1) + len(p2) - 2*pf
       if dist > maxDist {
           maxDist = dist
           maxV = i
       }
   }
   // output 1-based indices and distance
   u, v := maxV+1, from+1
   fmt.Printf("%d %d %d\n", u, v, maxDist)
}
