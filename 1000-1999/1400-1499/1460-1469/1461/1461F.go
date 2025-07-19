package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

const threshold = 1000000 // 10 * max N (100000)

var (
   a   []int
   mul []bool
)

func solve(l, r int) {
   // skip leading ones
   for l < r && a[l] == 1 {
       l++
   }
   // skip trailing ones
   for l < r && a[r] == 1 {
       r--
   }
   // collect positions of values >1
   where := make([]int, 0, r-l+1)
   var p int64 = 1
   for i := l; i <= r; i++ {
       if a[i] > 1 {
           p *= int64(a[i])
           if p >= threshold {
               // mark all between l+1..r as multiply
               for k := l + 1; k <= r; k++ {
                   mul[k] = true
               }
               return
           }
           where = append(where, i)
       }
   }
   m := len(where)
   // DP arrays
   maxi := make([]int64, m+1)
   prv := make([]int, m+1)
   // DP transitions
   for i := 0; i < m; i++ {
       p = 1
       // g is number of ones between where[i-1] and where[i]
       var g int64
       if i == 0 {
           g = 0
       } else {
           g = int64(where[i] - where[i-1] - 1)
       }
       for j := i; j < m; j++ {
           p *= int64(a[where[j]])
           // candidate score: maxi[i] + p + g
           if maxi[i]+p+g > maxi[j+1] {
               maxi[j+1] = maxi[i] + p + g
               prv[j+1] = i
           }
       }
   }
   // reconstruct choices, mark mul segments
   mm := m
   for mm > 0 {
       start := where[prv[mm]]
       end := where[mm-1]
       for k := start + 1; k <= end; k++ {
           mul[k] = true
       }
       mm = prv[mm]
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // read input
   var n int
   fmt.Fscan(in, &n)
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var ops string
   fmt.Fscan(in, &ops)
   // sort ops
   opsRunes := []rune(ops)
   sort.Slice(opsRunes, func(i, j int) bool { return opsRunes[i] < opsRunes[j] })
   ops = string(opsRunes)
   // handle simple cases
   if len(ops) == 1 || ops == "+-" {
       sep := string(opsRunes[0])
       for i := 0; i < n; i++ {
           if i > 0 {
               out.WriteString(sep)
           }
           out.WriteString(strconv.Itoa(a[i]))
       }
       out.WriteByte('\n')
       return
   }
   if ops == "*-" {
       for i := 0; i < n; i++ {
           if i > 0 {
               if a[i] == 0 {
                   out.WriteByte('-')
               } else {
                   out.WriteByte('*')
               }
           }
           out.WriteString(strconv.Itoa(a[i]))
       }
       out.WriteByte('\n')
       return
   }
   // general case
   mul = make([]bool, n)
   for i := 0; i < n; {
       j := i
       if a[i] == 0 {
           j = i + 1
       } else {
           for j < n && a[j] > 0 {
               j++
           }
           solve(i, j-1)
       }
       i = j
   }
   // output with chosen operators
   for i := 0; i < n; i++ {
       if i > 0 {
           if mul[i] {
               out.WriteByte('*')
           } else {
               out.WriteByte('+')
           }
       }
       out.WriteString(strconv.Itoa(a[i]))
   }
   out.WriteByte('\n')
}
