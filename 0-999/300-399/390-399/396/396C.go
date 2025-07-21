package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// BIT supports point update and prefix sum query
type BIT struct {
   n   int
   bit []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, bit: make([]int64, n+1)}
}

// Add val at position i
func (b *BIT) Add(i int, val int64) {
   for ; i <= b.n; i += i & -i {
       b.bit[i] = (b.bit[i] + val) % MOD
   }
}

// Sum returns prefix sum [1..i]
func (b *BIT) Sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s = (s + b.bit[i]) % MOD
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       children[p] = append(children[p], i)
   }
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   depth := make([]int, n+1)
   // Euler tour, iterative DFS
   idx := 1
   type frame struct{ u, ci int }
   stack := make([]frame, 0, n)
   stack = append(stack, frame{u: 1, ci: 0})
   depth[1] = 0
   for len(stack) > 0 {
       f := &stack[len(stack)-1]
       u := f.u
       if f.ci == 0 {
           tin[u] = idx
           idx++
       }
       if f.ci < len(children[u]) {
           v := children[u][f.ci]
           f.ci++
           depth[v] = depth[u] + 1
           stack = append(stack, frame{u: v, ci: 0})
       } else {
           tout[u] = idx - 1
           stack = stack[:len(stack)-1]
       }
   }

   bitA := NewBIT(n + 2)
   bitB := NewBIT(n + 2)

   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 1 {
           var v int
           var x, k int64
           fmt.Fscan(reader, &v, &x, &k)
           // Compute A = x + depth[v]*k mod
           A := (x + int64(depth[v])*k) % MOD
           // B = -k mod
           B := (MOD - k%MOD) % MOD
           l := tin[v]
           r := tout[v]
           // range add [l..r]
           bitA.Add(l, A)
           bitA.Add(r+1, (MOD-A)%MOD)
           bitB.Add(l, B)
           bitB.Add(r+1, (MOD-B)%MOD)
       } else {
           var v int
           fmt.Fscan(reader, &v)
           pos := tin[v]
           sumA := bitA.Sum(pos)
           sumB := bitB.Sum(pos)
           res := (sumA + sumB*int64(depth[v])) % MOD
           if res < 0 {
               res += MOD
           }
           fmt.Fprintln(writer, res)
       }
   }
}
