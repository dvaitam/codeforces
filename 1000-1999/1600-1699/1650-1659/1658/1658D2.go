package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   bitMax   = 17
   maxNodes = 3000005
)

var (
   nodes [maxNodes][2]int
   tot   int
)

func insert(x int) {
   p := 1
   for i := bitMax - 1; i >= 0; i-- {
       k := (x >> i) & 1
       if nodes[p][k] == 0 {
           tot++
           nodes[p][k] = tot
           // clear new node
           nodes[tot][0], nodes[tot][1] = 0, 0
       }
       p = nodes[p][k]
   }
}

// query1 returns the minimum xor of x with inserted values
func query1(x int) int {
   p, ans := 1, 0
   for i := bitMax - 1; i >= 0; i-- {
       k := (x >> i) & 1
       if nodes[p][k] != 0 {
           p = nodes[p][k]
       } else {
           p = nodes[p][k^1]
           ans |= 1 << i
       }
   }
   return ans
}

// query2 returns the maximum xor of x with inserted values
func query2(x int) int {
   p, ans := 1, 0
   for i := bitMax - 1; i >= 0; i-- {
       k := (x >> i) & 1
       if nodes[p][k^1] != 0 {
           p = nodes[p][k^1]
           ans |= 1 << i
       } else {
           p = nodes[p][k]
       }
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       n := r - l + 1
       // reset trie
       tot = 1
       nodes[1][0], nodes[1][1] = 0, 0
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
           insert(a[i])
       }
       // find x
       for i := 0; i < n; i++ {
           x := a[i] ^ l
           if query1(x) == l && query2(x) == r {
               fmt.Fprintln(writer, x)
               break
           }
       }
   }
}
