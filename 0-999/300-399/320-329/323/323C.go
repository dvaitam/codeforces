package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

var (
   rdr = bufio.NewReader(os.Stdin)
   wtr = bufio.NewWriter(os.Stdout)
)

// fast read integer (positive)
func readInt() int {
   var c byte
   var err error
   // skip non-digit
   for {
       c, err = rdr.ReadByte()
       if err != nil {
           return 0
       }
       if c >= '0' && c <= '9' {
           break
       }
   }
   x := int(c - '0')
   for {
       c, err = rdr.ReadByte()
       if err != nil || c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   return x
}

func lowerBound(a []int, v int) int {
   // first index i: a[i] >= v
   i := sort.Search(len(a), func(i int) bool { return a[i] >= v })
   return i
}
func upperBound(a []int, v int) int {
   // first index i: a[i] > v
   i := sort.Search(len(a), func(i int) bool { return a[i] > v })
   return i
}

func main() {
   defer wtr.Flush()
   n := readInt()
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = readInt()
   }
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       qv := readInt()
       pos[qv] = i
   }
   // build array A: positions in q of p
   A := make([]int, n+1)
   for i := 1; i <= n; i++ {
       A[i] = pos[p[i]]
   }
   // Fenwick tree of sorted slices
   bit := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       v := A[i]
       for j := i; j <= n; j += j & -j {
           bit[j] = append(bit[j], v)
       }
   }
   for i := 1; i <= n; i++ {
       if len(bit[i]) > 1 {
           sort.Ints(bit[i])
       }
   }
   m := readInt()
   x := 0
   // process queries
   for qi := 0; qi < m; qi++ {
       a := readInt(); b := readInt(); c := readInt(); d := readInt()
       f := func(z int) int { return ((z-1+x)%n) + 1 }
       f1 := f(a); f2 := f(b)
       l1, r1 := f1, f2
       if l1 > r1 { l1, r1 = r1, l1 }
       f3 := f(c); f4 := f(d)
       l2, r2 := f3, f4
       if l2 > r2 { l2, r2 = r2, l2 }
       ans := query(bit, r1, l2, r2) - query(bit, l1-1, l2, r2)
       wtr.WriteString(strconv.Itoa(ans))
       wtr.WriteByte('\n')
       x = ans + 1
   }
}

func query(bit [][]int, idx, l, r int) int {
   res := 0
   for i := idx; i > 0; i -= i & -i {
       s := bit[i]
       lo := lowerBound(s, l)
       hi := upperBound(s, r)
       res += hi - lo
   }
   return res
}
