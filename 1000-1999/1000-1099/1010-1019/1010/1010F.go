package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353
const maxLog = 18

var (
   n    int
   mVal int64
   adj  [][]int
   sz, son, anc []int
   fac, invf, upw []int
   w, iw [][]int
   a, f [][]int
)

func addMod(x, y int) int { x += y; if x >= mod { x -= mod }; return x }
func subMod(x, y int) int { x -= y; if x < 0 { x += mod }; return x }

func fpow(a, b int) int {
   res := 1
   base := a
   for b > 0 {
       if b&1 == 1 {
           res = int(int64(res) * base % mod)
       }
       base = int(int64(base) * base % mod)
       b >>= 1
   }
   return res
}

func prework() {
   w = make([][]int, maxLog+1)
   iw = make([][]int, maxLog+1)
   for i := 1; i <= maxLog; i++ {
       size := 1 << i
       half := size >> 1
       w[i] = make([]int, half)
       iw[i] = make([]int, half)
       root := fpow(3, (mod-1)/size)
       iroot := fpow(root, mod-2)
       w[i][0] = 1
       iw[i][0] = 1
       for j := 1; j < half; j++ {
           w[i][j] = int(int64(w[i][j-1]) * root % mod)
           iw[i][j] = int(int64(iw[i][j-1]) * iroot % mod)
       }
   }
}

// convolution of two polynomials
func mulPoly(x, y []int) []int {
   n1, n2 := len(x), len(y)
   lim := n1 + n2 - 2
   size := 1
   for size <= lim {
       size <<= 1
   }
   // bit-reversal
   pos := make([]int, size)
   for i := 1; i < size; i++ {
       pos[i] = pos[i>>1] >> 1
       if i&1 == 1 {
           pos[i] |= size >> 1
       }
   }
   // prepare buffers
   fa := make([]int, size)
   fb := make([]int, size)
   copy(fa, x)
   copy(fb, y)
   // NTT forward on fa
   for i := 1; i < size; i++ {
       if i < pos[i] {
           fa[i], fa[pos[i]] = fa[pos[i]], fa[i]
       }
   }
   for length, lvl := 2, 1; length <= size; length, lvl = length<<1, lvl+1 {
       half := length >> 1
       for i := 0; i < size; i += length {
           for j := 0; j < half; j++ {
               u := fa[i+j]
               v := int(int64(fa[i+j+half]) * w[lvl][j] % mod)
               fa[i+j] = addMod(u, v)
               fa[i+j+half] = subMod(u, v)
           }
       }
   }
   // NTT forward on fb
   for i := 1; i < size; i++ {
       if i < pos[i] {
           fb[i], fb[pos[i]] = fb[pos[i]], fb[i]
       }
   }
   for length, lvl := 2, 1; length <= size; length, lvl = length<<1, lvl+1 {
       half := length >> 1
       for i := 0; i < size; i += length {
           for j := 0; j < half; j++ {
               u := fb[i+j]
               v := int(int64(fb[i+j+half]) * w[lvl][j] % mod)
               fb[i+j] = addMod(u, v)
               fb[i+j+half] = subMod(u, v)
           }
       }
   }
   // point-wise multiply
   for i := 0; i < size; i++ {
       fa[i] = int(int64(fa[i]) * fb[i] % mod)
   }
   // NTT inverse on fa
   // bit-reversal again
   for i := 1; i < size; i++ {
       if i < pos[i] {
           fa[i], fa[pos[i]] = fa[pos[i]], fa[i]
       }
   }
   for length, lvl := 2, 1; length <= size; length, lvl = length<<1, lvl+1 {
       half := length >> 1
       for i := 0; i < size; i += length {
           for j := 0; j < half; j++ {
               u := fa[i+j]
               v := int(int64(fa[i+j+half]) * iw[lvl][j] % mod)
               fa[i+j] = addMod(u, v)
               fa[i+j+half] = subMod(u, v)
           }
       }
   }
   // divide by size
   invSize := fpow(size, mod-2)
   for i := 0; i < size; i++ {
       fa[i] = int(int64(fa[i]) * int64(invSize) % mod)
   }
   // trim to limit
   res := make([]int, lim+1)
   copy(res, fa[:lim+1])
   return res
}

func addPoly(x, y []int) []int {
   nx, ny := len(x), len(y)
   nn := nx
   if ny > nn { nn = ny }
   res := make([]int, nn)
   for i := 0; i < nx; i++ {
       res[i] = addMod(res[i], x[i])
   }
   for i := 0; i < ny; i++ {
       res[i] = addMod(res[i], y[i])
   }
   return res
}

func upPoly(x []int) []int {
   res := make([]int, len(x)+1)
   for i, v := range x {
       res[i+1] = v
   }
   return res
}

func dfs(u, p int) {
   sz[u] = 1
   anc[u] = p
   for _, v := range adj[u] {
       if v == p { continue }
       dfs(v, u)
       if sz[v] > sz[son[u]] {
           son[u] = v
       }
       sz[u] += sz[v]
   }
}

// returns (fi, se)
func calc(l, r int) ([]int, []int) {
   if l == r {
       return a[l], a[l]
   }
   mid := (l + r) >> 1
   f1, s1 := calc(l, mid)
   f2, s2 := calc(mid+1, r)
   fi := addPoly(f1, mulPoly(s1, f2))
   se := mulPoly(s1, s2)
   return fi, se
}

func solve(u int) {
   if son[u] == 0 {
       f[u] = []int{1, 1}
       return
   }
   // solve all side branches along heavy path
   for i := u; i != 0; i = son[i] {
       for _, v := range adj[i] {
           if v != son[i] && v != anc[i] {
               solve(v)
           }
       }
   }
   // build polynomials for segments
   cnt := 0
   for i := u; i != 0; i = son[i] {
       cnt++
       a[cnt] = []int{0, 1}
       for _, v := range adj[i] {
           if v != son[i] && v != anc[i] {
               a[cnt] = upPoly(f[v])
               f[v] = nil
           }
       }
   }
   fi, se := calc(1, cnt-1)
   tmp := addPoly(fi, upPoly(se))
   f[u] = addPoly([]int{1}, tmp)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &mVal)
   // adjust m
   mVal = (mVal + 1) % mod
   adj = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   sz = make([]int, n+1)
   son = make([]int, n+1)
   anc = make([]int, n+1)
   fac = make([]int, n+1)
   invf = make([]int, n+1)
   upw = make([]int, n+1)
   a = make([][]int, n+2)
   f = make([][]int, n+2)
   prework()
   fac[0] = 1
   for i := 1; i <= n; i++ {
       fac[i] = int(int64(fac[i-1]) * int64(i) % mod)
   }
   invf[n] = fpow(fac[n], mod-2)
   for i := n; i >= 1; i-- {
       invf[i-1] = int(int64(invf[i]) * int64(i) % mod)
   }
   upw[0] = 1
   mInt := int(mVal)
   for i := 1; i <= n; i++ {
       upw[i] = int(int64(upw[i-1]) * int64(mInt+i-1) % mod)
   }
   dfs(1, 0)
   solve(1)
   res := f[1]
   ans := subMod(res[0], 1)
   for i := 1; i <= n && i < len(res); i++ {
       coef := int(int64(upw[i-1]) * int64(invf[i-1]) % mod)
       ans = addMod(ans, int(int64(coef) * int64(res[i]) % mod))
   }
   fmt.Println(ans)
}
