package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int { a += b; if a >= MOD { a -= MOD }; return a }
func mul(a, b int) int { return int((int64(a) * int64(b)) % MOD) }
func pow2init(n int) []int {
   p := make([]int, n+1)
   p[0] = 1
   for i := 1; i <= n; i++ {
       p[i] = add(p[i-1], p[i-1])
   }
   return p
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M int
   fmt.Fscan(in, &N, &M)
   adj := make([][]int, N+1)
   for i := 0; i < M; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // if no cycles (tree), only one configuration (all gray)
   cExp := M - (N - 1)
   if cExp <= 0 {
       fmt.Println(1)
       return
   }
   // find cycles through node 1
   used := make(map[int]bool)
   cycles := []int{}
   for _, u := range adj[1] {
       if used[u] {
           continue
       }
       // traverse cycle starting from u
       prev := 1
       cur := u
       used[cur] = true
       length := 1
       for cur != 1 {
           // find next neighbor; skip previous unless cycle length is 2
           nxt := 0
           for _, w := range adj[cur] {
               if w != prev {
                   nxt = w
                   break
               }
           }
           if nxt == 0 {
               // only neighbor is prev (cycle length 2)
               nxt = prev
           }
           prev, cur = cur, nxt
           if cur != 1 {
               used[cur] = true
           }
           length++
       }
       cycles = append(cycles, length)
   }
   c := len(cycles)
   if c == 0 {
       fmt.Println(1)
       return
   }
   // count odd and even cycles
   o, e := 0, 0
   for _, L := range cycles {
       if L&1 == 1 {
           o++
       } else {
           e++
       }
   }
   // precompute factorials and invFactorials up to c
   maxN := c
   fact := make([]int, maxN+1)
   invFact := make([]int, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invFact[maxN] = 1
   // Fermat inverse of fact[maxN]
   invFact[maxN] = modInv(fact[maxN])
   for i := maxN; i > 0; i-- {
       invFact[i-1] = mul(invFact[i], i)
   }
   pow2 := pow2init(c)
   // precompute S[m] = sum_{s=0..m/2} m!/( (m-2s)! * s! )
   S := make([]int, c+1)
   for m := 0; m <= c; m++ {
       sum := 0
       for s := 0; 2*s <= m; s++ {
           // term = m!/( (m-2s)! * s! )
           term := fact[m]
           term = mul(term, invFact[m-2*s])
           term = mul(term, invFact[s])
           sum = add(sum, term)
       }
       S[m] = sum
   }
   // compute full coverage
   full := 0
   for t := 0; t <= o; t++ {
       Mrem := c - t
       if Mrem&1 == 1 {
           continue
       }
       k := Mrem / 2
       // C(o, t) * 2^t * (Mrem)!/k!
       ways := mul(comb(o, t, fact, invFact), pow2[t])
       ways = mul(ways, mul(fact[Mrem], invFact[k]))
       full = add(full, ways)
   }
   // compute early stop
   early := 0
   if e > 0 {
       Nrem := c - 1
       // o' = o, cycles excluding one even
       for t := 0; t <= o; t++ {
           if t > o {
               break
           }
           // C(o, t) * 2^t * S[Nrem - t]
           if Nrem-t < 0 {
               break
           }
           ways := mul(comb(o, t, fact, invFact), pow2[t])
           ways = mul(ways, S[Nrem-t])
           early = add(early, ways)
       }
       early = mul(early, mul(2*e%MOD, 1))
   }
   ans := add(full, early)
   fmt.Println(ans)
}

// comb computes C(n, k)
func comb(n, k int, fact, invFact []int) int {
   if k < 0 || k > n {
       return 0
   }
   return mul(fact[n], mul(invFact[k], invFact[n-k]))
}

// modInv computes modular inverse via Fermat's little theorem
func modInv(a int) int {
   return modPow(a, MOD-2)
}

func modPow(a, e int) int {
   res := 1
   base := a % MOD
   for e > 0 {
       if e&1 == 1 {
           res = mul(res, base)
       }
       base = mul(base, base)
       e >>= 1
   }
   return res
}
