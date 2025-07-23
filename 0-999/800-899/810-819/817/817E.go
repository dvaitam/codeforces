package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxB = 30

type node struct {
   child [2]*node
   cnt   int
}

var root node

func add(x int, delta int) {
   cur := &root
   for i := maxB - 1; i >= 0; i-- {
      b := (x >> i) & 1
      if cur.child[b] == nil {
         cur.child[b] = &node{}
      }
      cur = cur.child[b]
      cur.cnt += delta
   }
}

func query(x int, l int) int {
   cur := &root
   res := 0
   for i := maxB - 1; i >= 0; i-- {
      if cur == nil {
         break
      }
      xb := (x >> i) & 1
      lb := (l >> i) & 1
      if lb == 1 {
         if cur.child[xb] != nil {
            res += cur.child[xb].cnt
         }
         cur = cur.child[xb^1]
      } else {
         cur = cur.child[xb]
      }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
      return
   }
   for i := 0; i < q; i++ {
      var t int
      fmt.Fscan(reader, &t)
      if t == 1 {
         var x int
         fmt.Fscan(reader, &x)
         add(x, 1)
      } else if t == 2 {
         var x int
         fmt.Fscan(reader, &x)
         add(x, -1)
      } else if t == 3 {
         var x, l int
         fmt.Fscan(reader, &x, &l)
         fmt.Fprintln(writer, query(x, l))
      }
   }
}

