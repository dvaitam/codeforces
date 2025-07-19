package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const LOG = 61

// LinearRec computes k-th term of a linear recurrence of order n:
// dp[i] = sum_{j=0..n-1} trans[j] * dp[i-n+j]
// given initial dp[0..n-1] = first[0..n-1].
type LinearRec struct {
   n     int
   first []int64
   trans []int64
   bin   [][]int64
}

// add convolves two polynomials a and b and reduces by characteristic polynomial
func (lr *LinearRec) add(a, b []int64) []int64 {
   n := lr.n
   // result degree up to 2n
   tmp := make([]int64, 2*n+1)
   for i := 0; i <= n; i++ {
       ai := a[i]
       if ai == 0 {
           continue
       }
       for j := 0; j <= n; j++ {
           tmp[i+j] = (tmp[i+j] + ai*b[j]) % MOD
       }
   }
   // reduce terms degree > n
   for i := 2*n; i > n; i-- {
       if tmp[i] == 0 {
           continue
       }
       coef := tmp[i]
       tmp[i] = 0
       for j := 0; j < n; j++ {
           tmp[i-1-j] = (tmp[i-1-j] + coef*lr.trans[j]) % MOD
       }
   }
   // keep 0..n
   res := make([]int64, n+1)
   copy(res, tmp[:n+1])
   return res
}

// NewLinearRec constructs the recurrence with given first and trans
func NewLinearRec(first, trans []int64) *LinearRec {
   n := len(first)
   lr := &LinearRec{n: n, first: make([]int64, n), trans: make([]int64, n)}
   copy(lr.first, first)
   copy(lr.trans, trans)
   // build bin: bin[i] is polynomial for x^(2^i)
   // base polynomial for x is [0,1,0,...,0]
   a := make([]int64, n+1)
   a[1] = 1
   lr.bin = append(lr.bin, a)
   for i := 1; i < LOG; i++ {
       lr.bin = append(lr.bin, lr.add(lr.bin[i-1], lr.bin[i-1]))
   }
   return lr
}

// Calc returns dp[k] (0-indexed) for k>=0
func (lr *LinearRec) Calc(k uint64) int64 {
   n := lr.n
   // a represents polynomial for x^k
   a := make([]int64, n+1)
   a[0] = 1
   for i := 0; i < LOG; i++ {
       if (k>>i)&1 != 0 {
           a = lr.add(a, lr.bin[i])
       }
   }
   var ret int64
   for i := 0; i < n; i++ {
       ret = (ret + a[i+1]*lr.first[i]) % MOD
   }
   return ret
}

func readInt64(r *bufio.Reader) int64 {
   var x int64
   var neg bool
   b, err := r.ReadByte()
   for err == nil && (b < '0' || b > '9') && b != '-' {
       b, err = r.ReadByte()
   }
   if err != nil {
       return 0
   }
   if b == '-' {
       neg = true
       b, _ = r.ReadByte()
   }
   for ; b >= '0' && b <= '9'; b, _ = r.ReadByte() {
       x = x*10 + int64(b-'0')
   }
   if neg {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read N and M
   N := readInt64(reader)
   M := int(readInt64(reader))
   // dp recurrence: dp[i] = dp[i-1] + dp[i-M]
   // initial dp[0..M-1] = 1
   first := make([]int64, M)
   for i := 0; i < M; i++ {
       first[i] = 1
   }
   trans := make([]int64, M)
   trans[0] = 1
   trans[M-1] = 1
   lr := NewLinearRec(first, trans)
   // shift index: our dp[0] corresponds to code dp[1], so compute dp[N+1]
   res := lr.Calc(uint64(N + 1))
   fmt.Println(res)
}
