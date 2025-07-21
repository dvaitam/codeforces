package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU structure
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := 0; i <= n; i++ {
       p[i] = i
   }
   return &DSU{p: p}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx != ry {
       d.p[ry] = rx
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // sieve primes
   isPrime := make([]bool, n+1)
   for i := 2; i <= n; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= n; i++ {
       if isPrime[i] {
           for j := i * i; j <= n; j += i {
               isPrime[j] = false
           }
       }
   }
   // build DSU
   dsu := NewDSU(n)
   for p := 2; p <= n; p++ {
       if !isPrime[p] {
           continue
       }
       for k := 2; p*k <= n; k++ {
           dsu.Union(p, p*k)
       }
   }
   // group positions
   groups := make(map[int][]int)
   for i := 1; i <= n; i++ {
       r := dsu.Find(i)
       groups[r] = append(groups[r], i)
   }
   // collect groups in slice
   grpList := make([][]int, 0, len(groups))
   for _, g := range groups {
       grpList = append(grpList, g)
   }
   // sort groups by decreasing size
   // simple insertion sort for small n
   for i := 1; i < len(grpList); i++ {
       j := i
       for j > 0 && len(grpList[j-1]) < len(grpList[j]) {
           grpList[j-1], grpList[j] = grpList[j], grpList[j-1]
           j--
       }
   }
   // count letter frequencies
   freq := make([]int, 26)
   for _, ch := range s {
       freq[ch-'a']++
   }
   // result array
   res := make([]byte, n)
   // assign groups
   for _, g := range grpList {
       size := len(g)
       // find a letter with enough freq
       idx := -1
       for j := 0; j < 26; j++ {
           if freq[j] >= size {
               idx = j
               break
           }
       }
       if idx == -1 {
           fmt.Println("NO")
           return
       }
       // assign
       for _, pos := range g {
           res[pos-1] = byte('a' + idx)
       }
       freq[idx] -= size
   }
   fmt.Println("YES")
   fmt.Println(string(res))
}
