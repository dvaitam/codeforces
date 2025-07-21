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
   edges := make([][2]int, m+1)
   for i := 1; i <= m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       edges[i][0] = x - 1
       edges[i][1] = y - 1
   }
   // Build prefix DSUs
   prefixParent := make([][]uint16, m+1)
   prefixCount := make([]int, m+1)
   initP := make([]uint16, n)
   for i := 0; i < n; i++ {
       initP[i] = uint16(i)
   }
   prefixParent[0] = initP
   prefixCount[0] = n
   for i := 1; i <= m; i++ {
       prev := prefixParent[i-1]
       p := make([]uint16, n)
       copy(p, prev)
       cnt := prefixCount[i-1]
       a := uint16(edges[i][0])
       b := uint16(edges[i][1])
       ra := find(p, a)
       rb := find(p, b)
       if ra != rb {
           p[rb] = ra
           cnt--
       }
       prefixParent[i] = p
       prefixCount[i] = cnt
   }
   // Build suffix DSUs
   suffixParent := make([][]uint16, m+2)
   suffixCount := make([]int, m+2)
   initS := make([]uint16, n)
   for i := 0; i < n; i++ {
       initS[i] = uint16(i)
   }
   suffixParent[m+1] = initS
   suffixCount[m+1] = n
   for i := m; i >= 1; i-- {
       prev := suffixParent[i+1]
       p := make([]uint16, n)
       copy(p, prev)
       cnt := suffixCount[i+1]
       a := uint16(edges[i][0])
       b := uint16(edges[i][1])
       ra := find(p, a)
       rb := find(p, b)
       if ra != rb {
           p[rb] = ra
           cnt--
       }
       suffixParent[i] = p
       suffixCount[i] = cnt
   }
   // Process queries
   var k int
   fmt.Fscan(reader, &k)
   for qi := 0; qi < k; qi++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // Initialize with prefix state
       curParent := make([]uint16, n)
       copy(curParent, prefixParent[l-1])
       curCount := prefixCount[l-1]
       sp := suffixParent[r+1]
       // Merge suffix components
       seen := make([]int, n)
       for i := 0; i < n; i++ {
           root := find(sp, uint16(i))
           if seen[root] == 0 {
               seen[root] = i + 1
           } else {
               rep := seen[root] - 1
               ri := uint16(i)
               r0 := uint16(rep)
               cr0 := find(curParent, r0)
               cr1 := find(curParent, ri)
               if cr0 != cr1 {
                   curParent[cr1] = cr0
                   curCount--
               }
           }
       }
       fmt.Fprintln(writer, curCount)
   }
}

// find returns the root of x in DSU p without path compression
func find(p []uint16, x uint16) uint16 {
   for p[x] != x {
       x = p[x]
   }
   return x
}
