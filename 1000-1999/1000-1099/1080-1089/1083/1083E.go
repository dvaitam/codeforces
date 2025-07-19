package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func readInt64() int64 {
   var x int64
   var c byte
   var err error
   sign := int64(1)
   // skip non-digits
   for {
       c, err = reader.ReadByte()
       if err != nil {
           return x * sign
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   for ; err == nil && c >= '0' && c <= '9'; c, err = reader.ReadByte() {
       x = x*10 + int64(c-'0')
   }
   return x * sign
}

type square struct {
   x, y, w int64
}

func main() {
   defer writer.Flush()
   n := int(readInt64())
   s := make([]square, n)
   for i := 0; i < n; i++ {
       s[i].x = readInt64()
       s[i].y = readInt64()
       s[i].w = readInt64()
   }
   // sort by y descending
   sort.Slice(s, func(i, j int) bool {
       return s[i].y > s[j].y
   })
   f := make([]int64, n)
   q := make([]int, 0, n)
   var ans int64
   eps := 1e-12
   // functions for slope
   getSlope := func(u, v int) float64 {
       dx := float64(s[v].x - s[u].x)
       if math.Abs(dx) < eps {
           return math.Inf(1)
       }
       return float64(f[v]-f[u]) / dx
   }
   l := 0
   // process
   for i := 0; i < n; i++ {
       // find best previous
       for len(q)-l >= 2 && getSlope(q[l], q[l+1]) >= float64(s[i].y) {
           l++
       }
       // base value
       fi := s[i].x*s[i].y - s[i].w
       if len(q)-l > 0 {
           j := q[l]
           val := f[j] + (s[i].x-s[j].x)*s[i].y - s[i].w
           if val > fi {
               fi = val
           }
       }
       f[i] = fi
       if fi > ans {
           ans = fi
       }
       // maintain deque
       for len(q)-l >= 2 && getSlope(q[len(q)-2], q[len(q)-1]) <= getSlope(q[len(q)-1], i) {
           q = q[:len(q)-1]
       }
       q = append(q, i)
   }
   fmt.Fprintln(writer, ans)
}
