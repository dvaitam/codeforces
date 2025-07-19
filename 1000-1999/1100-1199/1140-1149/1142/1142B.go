package main

import (
   "bufio"
   "os"
)

// Fast IO
var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func main() {
   defer writer.Flush()
   n := readInt()
   m := readInt()
   q := readInt()
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       v := readInt()
       p[v] = i
   }
   T := make([]int, n+1)
   children := make([][]int, m+1)
   for i := 1; i <= m; i++ {
       v := readInt()
       x := p[v]
       var prev int
       if x == 1 {
           prev = T[n]
       } else {
           prev = T[x-1]
       }
       children[prev] = append(children[prev], i)
       T[x] = i
   }
   left := make([]int, m+1)
   Q := make([]int, m+1)
   type frame struct{ node, dep, next int }
   stack := make([]frame, 0, m+1)
   stack = append(stack, frame{node: 0, dep: 0, next: 0})
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       if top.next == 0 {
           dep := top.dep
           node := top.node
           Q[dep] = node
           if dep >= n {
               left[node] = Q[dep-n+1]
           } else {
               left[node] = -1
           }
       }
       node := top.node
       if top.next < len(children[node]) {
           child := children[node][top.next]
           top.next++
           stack = append(stack, frame{node: child, dep: top.dep + 1, next: 0})
       } else {
           stack = stack[:len(stack)-1]
       }
   }
   for i := 1; i <= m; i++ {
       if left[i] < left[i-1] {
           left[i] = left[i-1]
       }
   }
   for i := 0; i < q; i++ {
       l := readInt()
       r := readInt()
       if left[r] >= l {
           writer.WriteByte('1')
       } else {
           writer.WriteByte('0')
       }
   }
}

// readInt reads next integer from stdin
func readInt() int {
   sign := 1
   b, _ := reader.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = reader.ReadByte()
   }
   if b == '-' {
       sign = -1
       b, _ = reader.ReadByte()
   }
   x := 0
   for ; b >= '0' && b <= '9'; b, _ = reader.ReadByte() {
       x = x*10 + int(b-'0')
   }
   return x * sign
}
