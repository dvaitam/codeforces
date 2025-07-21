package main

import (
   "bufio"
   "fmt"
   "os"
)

func nextPow2(n int) int {
   p := 1
   for p < n {
       p <<= 1
   }
   return p
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   size := nextPow2(n)
   // segment tree of size size, nodes 1..2*size-1
   const M = 60
   total := 2 * size * M
   f := make([]int32, total)
   // init leaves
   for i := 0; i < size; i++ {
       node := (i + size) * M
       if i < n {
           ai := a[i]
           for t := 0; t < M; t++ {
               if t%ai == 0 {
                   f[node+t] = 2
               } else {
                   f[node+t] = 1
               }
           }
       } else {
           for t := 0; t < M; t++ {
               f[node+t] = 0
           }
       }
   }
   // build
   for i := size - 1; i >= 1; i-- {
       left := 2 * i
       right := left + 1
       base := i * M
       lbase := left * M
       rbase := right * M
       for t := 0; t < M; t++ {
           tL := f[lbase+t]
           tR := f[rbase+((t+int(tL))%M)]
           f[base+t] = tL + tR
       }
   }

   var q int
   fmt.Fscan(reader, &q)
   // buffers for query composition
   tmp := make([]int32, M)
   leftArr := make([]int32, M)
   // process queries
   for qi := 0; qi < q; qi++ {
       var typ string
       var x, y int
       fmt.Fscan(reader, &typ, &x, &y)
       if typ == "A" {
           // query [x-1, y-2]
           l := x - 1 + size
           r := y - 2 + size
           // result f_cur
           for t := 0; t < M; t++ {
               leftArr[t] = 0
           }
           // collect nodes left to right
           rightNodes := make([]int, 0)
           for l <= r {
               if l&1 == 1 {
                   // compose leftArr with f[l]
                   base := l * M
                   for t := 0; t < M; t++ {
                       tL := leftArr[t]
                       tmp[t] = tL + f[base+((t+int(tL))%M)]
                   }
                   copy(leftArr, tmp)
                   l++
               }
               if r&1 == 0 {
                   rightNodes = append(rightNodes, r)
                   r--
               }
               l >>= 1
               r >>= 1
           }
           // apply right nodes in reverse order
           for i := len(rightNodes) - 1; i >= 0; i-- {
               node := rightNodes[i]
               base := node * M
               for t := 0; t < M; t++ {
                   tL := leftArr[t]
                   tmp[t] = tL + f[base+((t+int(tL))%M)]
               }
               copy(leftArr, tmp)
           }
           // answer is leftArr[0]
           fmt.Fprintln(writer, leftArr[0])
       } else {
           // update a[x] = y, segment index x-1
           pos := x - 1
           a[pos] = y
           nodeIdx := pos + size
           base := nodeIdx * M
           for t := 0; t < M; t++ {
               if t%y == 0 {
                   f[base+t] = 2
               } else {
                   f[base+t] = 1
               }
           }
           // update up
           for nodeIdx >>= 1; nodeIdx >= 1; nodeIdx >>= 1 {
               left := 2 * nodeIdx
               right := left + 1
               base := nodeIdx * M
               lbase := left * M
               rbase := right * M
               for t := 0; t < M; t++ {
                   tL := f[lbase+t]
                   tR := f[rbase+((t+int(tL))%M)]
                   f[base+t] = tL + tR
               }
           }
       }
   }
}
