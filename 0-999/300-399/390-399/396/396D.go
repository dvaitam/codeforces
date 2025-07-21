package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const inv2 = 500000004  // modular inverse of 2
const inv4 = 250000002  // modular inverse of 4

// BIT implements a Fenwick tree for sum queries
type BIT struct {
   n    int
   tree []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// add v at position i (1-indexed)
func (b *BIT) update(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// sum returns sum of [1..i]
func (b *BIT) sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}


func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }

   // precompute factorials
   fact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }

   bit := NewBIT(n)
   // sum of used values
   var sumYAll int64
   var ans, prefixInv int64

   for i := 1; i <= n; i++ {
       pi := p[i-1]
       // count used elements less than pi
       usedLess := bit.sum(pi - 1)
       c := int64(pi-1-usedLess)
       // total sum of used values
       // number of used elements
       U := int64(i - 1)

       rem := n - i
       factRem := fact[rem]

       // termA: c * prefixInv * factRem
       termA := c * prefixInv % MOD * factRem % MOD
       // termB: factRem * (c*(c-1)/2)
       termB := factRem * (c*(c-1)%MOD * inv2 % MOD) % MOD

       // suffix inv sum for rem > 1: factRem * rem*(rem-1)/4
       var suffixSum int64
       if rem >= 2 {
           suffixSum = factRem * int64(rem) % MOD * int64(rem-1) % MOD * inv4 % MOD
       }
       suffixTotal := c * suffixSum % MOD

       // missingTerm: inv between prefix elements and suffix elements
       // Σ_y T_y = sumYAll - U - U*(U-1)/2
       uC2 := U * (U - 1) % MOD * inv2 % MOD
       temp := (sumYAll - U - uC2) % MOD
       if temp < 0 {
           temp += MOD
       }
       missing := factRem * c % MOD * temp % MOD
       ans = (ans + termA + termB + suffixTotal + missing) % MOD

       // update prefix inversion count with actual p_i
       prefixInv = (prefixInv + int64((i-1)-usedLess)) % MOD
       bit.update(pi, 1)
       sumYAll = (sumYAll + int64(pi)) % MOD
   }
   // include the original permutation's inversions
   ans = (ans + prefixInv) % MOD
   fmt.Fprint(writer, ans)
}
