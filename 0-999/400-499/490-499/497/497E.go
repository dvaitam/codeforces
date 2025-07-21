package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modAdd(a, b int) int { a += b; if a >= MOD { a -= MOD }; return a }
func modSub(a, b int) int { a -= b; if a < 0 { a += MOD }; return a }
func modMul(a, b int) int { return int((int64(a) * int64(b)) % MOD) }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n uint64
   var k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // prepare base matrices M_base[s] = M(s)
   m := k + 1
   M_base := make([][][]int, k)
   for s := 0; s < k; s++ {
       M := makeMatrix(m)
       // dp': 2*dp - last_contrib[s]
       M[0][0] = 2
       M[0][1+s] = modSub(0, 1)
       // last_contrib'[s] = dp
       M[1+s][0] = 1
       // other last_contrib unchanged
       for i := 0; i < k; i++ {
           if i != s {
               M[1+i][1+i] = 1
           }
       }
       M_base[s] = M
   }
   // compute digits of n-1 in base k
   nn := n - 1
   var D []int
   if n == 0 {
       D = []int{0}
   } else {
       for x := nn; x > 0; x /= uint64(k) {
           D = append(D, int(x % uint64(k)))
       }
       if len(D) == 0 {
           D = append(D, 0)
       }
       // reverse to most significant first
       for i, j := 0, len(D)-1; i < j; i, j = i+1, j-1 {
           D[i], D[j] = D[j], D[i]
       }
   }
   L := len(D)
   // Precompute block matrices Pre[l][s]: product over k^l consecutive j's starting sum-digit mod = s
   // Pre[0][s] = M_base[s]
   Pre := make([][][][]int, L)
   Pre[0] = make([][][]int, k)
   for s := 0; s < k; s++ {
       Pre[0][s] = M_base[s]
   }
   // higher levels
   for lvl := 1; lvl < L; lvl++ {
       Pre[lvl] = make([][][]int, k)
       for s := 0; s < k; s++ {
           C := identityMatrix(m)
           for t := 0; t < k; t++ {
               idx := (s + t) % k
               C = matMul(C, Pre[lvl-1][idx])
           }
           Pre[lvl][s] = C
       }
   }
   // process all j in [0..n-1]
   Mtot := identityMatrix(m)
   sP := 0
   for pos := 0; pos < L; pos++ {
       rem := L - pos - 1
       d := D[pos]
       for t := 0; t < d; t++ {
           sp := (sP + t) % k
           Mtot = matMul(Mtot, Pre[rem][sp])
       }
       sP = (sP + d) % k
   }
   // last element j = n-1
   Mtot = matMul(Mtot, Pre[0][sP])
   // initial v0: v0[0]=1, others 0 => result dp = Mtot[0][0]
   ans := Mtot[0][0]
   fmt.Println(ans)
}

// makeMatrix allocates zero m x m
func makeMatrix(m int) [][]int {
   M := make([][]int, m)
   for i := range M {
       M[i] = make([]int, m)
   }
   return M
}

// identityMatrix returns m x m identity
func identityMatrix(m int) [][]int {
   I := makeMatrix(m)
   for i := 0; i < m; i++ {
       I[i][i] = 1
   }
   return I
}

// matMul multiplies A * B
func matMul(A, B [][]int) [][]int {
   m := len(A)
   C := makeMatrix(m)
   for i := 0; i < m; i++ {
       for k := 0; k < m; k++ {
           if A[i][k] == 0 {
               continue
           }
           aik := A[i][k]
           rowC := C[i]
           rowB := B[k]
           for j := 0; j < m; j++ {
               rowC[j] = (rowC[j] + int64(aik)*int64(rowB[j]) % MOD) % MOD
           }
       }
   }
   return C
}
