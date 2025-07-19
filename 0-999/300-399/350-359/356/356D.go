package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Bitset for sums up to s (max 70010)
type Bitset []uint64

func newBitset(words int) Bitset {
   return make(Bitset, words)
}
// clone b
func (b Bitset) clone() Bitset {
   nb := make(Bitset, len(b))
   copy(nb, b)
   return nb
}
// Set bit i
func (b Bitset) set(i int) {
   b[i>>6] |= 1 << uint(i&63)
}
// Test bit i
func (b Bitset) get(i int) bool {
   return (b[i>>6]&(1<<uint(i&63))) != 0
}
// OrShift: b |= b << shift
func (b Bitset) orShift(shift, words int) {
   w := shift >> 6
   offset := uint(shift & 63)
   // traverse from high to low to avoid overwrite
   for i := words - 1; i >= 0; i-- {
       var v uint64
       j := i - w
       if j < 0 {
           continue
       }
       v = b[j] << offset
       if offset != 0 && j-1 >= 0 {
           v |= b[j-1] >> (64 - offset)
       }
       b[i] |= v
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, s int
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // find max index p
   p := 0
   for i := 1; i < n; i++ {
       if a[i] > a[p] {
           p = i
       }
   }
   if s < a[p] {
       fmt.Fprintln(writer, -1)
       return
   }
   // prepare DP
   const S = 7
   m := (n + S - 1) / S
   // bitset length in bits = s+1, words in uint64
   bits := s + 1
   words := (bits + 63) >> 6
   f := make([]Bitset, m+1)
   f[0] = newBitset(words)
   f[0].set(0)
   // flags for chosen in DP
   flag := make([]bool, n)
   // build DP blocks
   for i := 0; i < m; i++ {
       cur := f[i].clone()
       // items in block
       for j := 0; j < S; j++ {
           k := i*S + j
           if k >= n || k == p {
               continue
           }
           shift := a[k]
           if shift <= s {
               cur.orShift(shift, words)
           }
       }
       f[i+1] = cur
   }
   t := s - a[p]
   if t < 0 || t > s || !f[m].get(t) {
       fmt.Fprintln(writer, -1)
       return
   }
   // backtrack
   for i := m; i > 0; i-- {
       // try masks
       base := i*S
       for mask := 0; mask < 1<<S; mask++ {
           val := 0
           ok := true
           for j := 0; j < S; j++ {
               if mask>>j&1 == 1 {
                   k := base - S + j
                   if k < 0 || k >= n || k == p {
                       ok = false
                       break
                   }
                   val += a[k]
               }
           }
           if !ok || val > t {
               continue
           }
           if f[i-1].get(t - val) {
               // select these
               for j := 0; j < S; j++ {
                   if mask>>j&1 == 1 {
                       k := base - S + j
                       flag[k] = true
                   }
               }
               t -= val
               break
           }
       }
   }
   flag[p] = true
   // build tree
   c := make([]int, n)
   adj := make([][]int, n)
   pool := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if !flag[i] {
           pool = append(pool, i)
       } else {
           c[i] = a[i]
       }
   }
   sort.Slice(pool, func(i, j int) bool { return a[pool[i]] < a[pool[j]] })
   pool = append(pool, p)
   for i, x := range pool {
       c[x] = a[x]
       if i > 0 {
           par := pool[i-1]
           adj[x] = append(adj[x], par)
           c[x] -= a[par]
       }
   }
   // verify sum
   sumC := 0
   for i := 0; i < n; i++ {
       sumC += c[i]
   }
   if sumC != s {
       fmt.Fprintln(writer, -1)
       return
   }
   // output
   for i := 0; i < n; i++ {
       fmt.Fprint(writer, c[i], " ", len(adj[i]))
       for _, v := range adj[i] {
           fmt.Fprint(writer, " ", v+1)
       }
       fmt.Fprintln(writer)
   }
}
