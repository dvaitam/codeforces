package main

import (
   "bufio"
   "os"
   "strconv"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func readInt() int {
   var c byte
   var err error
   // skip non-numeric
   for {
       c, err = reader.ReadByte()
       if err != nil {
           return 0
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   x := 0
   for ; err == nil && c >= '0' && c <= '9'; c, err = reader.ReadByte() {
       x = x*10 + int(c - '0')
   }
   return x * sign
}

type node struct {
   length, mod, val int
}

func main() {
   defer writer.Flush()
   n := readInt()
   k := readInt()
   // collect unique (length, mod) pairs
   seen := make(map[int]bool)
   uniqs := make([]node, 0)
   // precompute powers of 10 mod k
   pow10 := make([]int, 12)
   pow10[0] = 1 % k
   for i := 1; i < 12; i++ {
       pow10[i] = pow10[i-1] * 10 % k
   }
   for i := 0; i < n; i++ {
       x := readInt()
       v := x
       // compute digit length
       l := 0
       if v == 0 {
           l = 1
       } else {
           for t := v; t != 0; t /= 10 {
               l++
           }
       }
       m := v % k
       key := l*(k+1) + m
       if !seen[key] {
           seen[key] = true
           uniqs = append(uniqs, node{l, m, v})
       }
   }
   // Dijkstra (Dial's algorithm)
   inf := 1e18
   size := 10*k + 5
   dist := make([]int, k+1)
   prevRem := make([]int, k+1)
   prevChoice := make([]int, k+1)
   visited := make([]bool, k+1)
   for i := 1; i <= k; i++ {
       dist[i] = int(inf)
   }
   // buckets for distances
   buckets := make([][]int, size)
   heads := make([]int, size)
   // start from remainder 0, represented by 0
   buckets[0] = append(buckets[0], 0)
   dist[0] = 0
   cnt := 1
   curD := 0
   found := false
   for curD < size {
       if cnt == 0 {
           break
       }
       if heads[curD] >= len(buckets[curD]) {
           curD++
           continue
       }
       u := buckets[curD][heads[curD]]
       heads[curD]++
       cnt--
       // map raw 0 to target representation k
       if u == k {
           found = true
           break
       }
       if visited[u] {
           continue
       }
       visited[u] = true
       // relax edges
       for idx, nd := range uniqs {
           vmod := (u*pow10[nd.length] + nd.mod) % k
           v := vmod
           if vmod == 0 {
               v = k
           }
           if visited[v] {
               continue
           }
           newD := dist[u] + nd.length
           if newD < dist[v] {
               dist[v] = newD
               prevRem[v] = u
               prevChoice[v] = idx
               if newD < size {
                   buckets[newD] = append(buckets[newD], v)
                   cnt++
               }
           }
       }
   }
   if !found {
       writer.WriteString("NO")
       return
   }
   writer.WriteString("YES\n")
   // reconstruct path
   res := make([]int, 0)
   cur := k
   for cur != 0 {
       idx := prevChoice[cur]
       res = append(res, uniqs[idx].val)
       cur = prevRem[cur]
   }
   // reverse and output
   for i := len(res) - 1; i >= 0; i-- {
       writer.WriteString(strconv.Itoa(res[i]))
   }
}
