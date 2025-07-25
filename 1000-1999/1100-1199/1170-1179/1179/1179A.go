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

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prepare deque in buffer
   bufSize := 3 * n
   buf := make([]int, bufSize)
   head := n
   tail := head + n - 1
   for i := 0; i < n; i++ {
       buf[head+i] = a[i]
   }
   // simulate first n operations
   type pair struct{ x, y int }
   ops := make([]pair, n+1)
   for i := 1; i <= n; i++ {
       A := buf[head]
       B := buf[head+1]
       ops[i] = pair{A, B}
       // pop two
       head += 2
       // push winner to front
       if A > B {
           head--
           buf[head] = A
           // push loser to back
           tail++
           buf[tail] = B
       } else {
           head--
           buf[head] = B
           tail++
           buf[tail] = A
       }
   }
   // after n ops, front is max element
   M := buf[head]
   // collect rest
   restLen := n - 1
   rest := make([]int, restLen)
   for i := 0; i < restLen; i++ {
       rest[i] = buf[head+1+i]
   }

   // answer queries
   for j := 0; j < q; j++ {
       var m int64
       fmt.Fscan(reader, &m)
       if m <= int64(n) {
           p := ops[int(m)]
           fmt.Fprintln(writer, p.x, p.y)
       } else {
           idx := (m - int64(n) - 1) % int64(restLen)
           fmt.Fprintln(writer, M, rest[int(idx)])
       }
   }
}
