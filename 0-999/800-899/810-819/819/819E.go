package main

import (
   "bufio"
   "fmt"
   "os"
)

type node struct { a, b, c, d int }

var ans []node

func add(n int) {
   if n&1 == 1 {
       ans = append(ans, node{n, 0, n + 1, -1})
       ans = append(ans, node{n, 0, n + 1, -1})
       for i := 1; i < n; i += 2 {
           ans = append(ans, node{n, i, n + 1, i + 1})
           ans = append(ans, node{n, i, n + 1, i + 1})
       }
   } else {
       ans = append(ans, node{n, 0, n + 1, -1})
       ans = append(ans, node{n, 1, n + 1, -1})
       ans = append(ans, node{n, 0, n + 1, 1})
       for i := 2; i < n; i += 2 {
           ans = append(ans, node{n, i, n + 1, i + 1})
           ans = append(ans, node{n, i, n + 1, i + 1})
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)

   if n&1 == 1 {
       ans = append(ans, node{0, 1, 2, -1})
       ans = append(ans, node{0, 1, 2, -1})
       for now := 3; now < n; now += 2 {
           add(now)
       }
   } else {
       ans = append(ans, node{0, 1, 2, -1})
       ans = append(ans, node{1, 2, 3, -1})
       ans = append(ans, node{2, 3, 0, -1})
       ans = append(ans, node{3, 0, 1, -1})
       for now := 4; now < n; now += 2 {
           add(now)
       }
   }

   fmt.Fprintln(writer, len(ans))
   for _, i := range ans {
       if i.d == -1 {
           fmt.Fprintln(writer, 3, i.a+1, i.b+1, i.c+1)
       } else {
           fmt.Fprintln(writer, 4, i.a+1, i.b+1, i.c+1, i.d+1)
       }
   }
