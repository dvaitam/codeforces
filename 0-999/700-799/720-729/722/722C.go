package main

import (
   "bufio"
   "fmt"
   "os"
)

var parent []int
var sumSegment []int64
var used []bool

func find(x int) int {
   if parent[x] != x {
       parent[x] = find(parent[x])
   }
   return parent[x]
}

func union(x, y int) {
   rx := find(x)
   ry := find(y)
   if rx == ry {
       return
   }
   parent[ry] = rx
   sumSegment[rx] += sumSegment[ry]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       p[i] = x - 1
   }
   parent = make([]int, n)
   sumSegment = make([]int64, n)
   used = make([]bool, n)
   ans := make([]int64, n)
   var currentMax int64
   // Process in reverse: activate elements and record max before each activation
   for i := n - 1; i >= 0; i-- {
       ans[i] = currentMax
       pos := p[i]
       used[pos] = true
       parent[pos] = pos
       sumSegment[pos] = a[pos]
       if pos > 0 && used[pos-1] {
           union(pos, pos-1)
       }
       if pos+1 < n && used[pos+1] {
           union(pos, pos+1)
       }
       root := find(pos)
       if sumSegment[root] > currentMax {
           currentMax = sumSegment[root]
       }
   }
   for i := 0; i < n; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
