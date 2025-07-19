package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// DSU with parity
type DSU struct {
   parent []int
   w      []int
}

func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   w := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       w[i] = 0
   }
   return &DSU{parent: parent, w: w}
}

func (d *DSU) find(x int) int {
   if d.parent[x] != x {
       root := d.find(d.parent[x])
       d.w[x] ^= d.w[d.parent[x]]
       d.parent[x] = root
   }
   return d.parent[x]
}

func (d *DSU) getw(x int) int {
   d.find(x)
   return d.w[x]
}

// merge u and v with parity s (0 if even, 1 if odd)
// returns false if conflict
func (d *DSU) merge(u, v, s int) bool {
   ru := d.find(u)
   rv := d.find(v)
   wu := d.w[u]
   wv := d.w[v]
   if ru != rv {
       // attach rv under ru
       d.parent[rv] = ru
       // set w[rv] so that wu ^ w[rv] ^ wv == s
       d.w[rv] = s ^ wu ^ wv
       return true
   }
   // same component, check consistency
   return (wu ^ wv ^ s) == 0
}

func readInt(reader *bufio.Reader) (int, error) {
   var sb []byte
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return 0, err
       }
       if b == '-' || (b >= '0' && b <= '9') {
           sb = append(sb, b)
           break
       }
   }
   for {
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           break
       }
       sb = append(sb, b)
   }
   return strconv.Atoi(string(sb))
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   t, err := readInt(reader)
   if err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       n, _ := readInt(reader)
       q, _ := readInt(reader)
       u := make([]int, n)
       v := make([]int, n)
       w := make([]int, n)
       dsu := NewDSU(n)
       for i := 1; i < n; i++ {
           ui, _ := readInt(reader)
           vi, _ := readInt(reader)
           wi, _ := readInt(reader)
           u[i] = ui
           v[i] = vi
           w[i] = wi
           if wi >= 0 {
               // parity of bits
               dsu.merge(ui, vi, bitsParity(wi))
           }
       }
       ok := true
       for i := 0; i < q; i++ {
           ui, _ := readInt(reader)
           vi, _ := readInt(reader)
           wi, _ := readInt(reader)
           if ok {
               if !dsu.merge(ui, vi, wi) {
                   ok = false
               }
           }
       }
       if ok {
           writer.WriteString("YES\n")
           for i := 1; i < n; i++ {
               wi := w[i]
               if wi < 0 {
                   // compute parity
                   wi = dsu.getw(u[i]) ^ dsu.getw(v[i])
               }
               writer.WriteString(
                   strconv.Itoa(u[i]) + " " + strconv.Itoa(v[i]) + " " + strconv.Itoa(wi) + "\n")
           }
       } else {
           writer.WriteString("NO\n")
       }
   }
}

func bitsParity(x int) int {
   // parity of set bits
   p := 0
   for x != 0 {
       p ^= x & 1
       x >>= 1
   }
   return p
}
