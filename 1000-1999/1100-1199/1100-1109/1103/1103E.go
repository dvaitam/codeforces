package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   M        = 10
   MaxRange = 100000
   Inv5     = 57646075230342349 // inverse of 5 mod 2^58
   Mod      = 1 << 58
)

// SumUnits holds counts in 10 units for convolution
type SumUnits struct {
   value [M]uint64
}

// Add adds other into su
func (su *SumUnits) Add(other *SumUnits) {
   for i := 0; i < M; i++ {
       su.value[i] += other.value[i]
   }
}

// Mul returns the convolution (mod M) of su and other
func (su SumUnits) Mul(other SumUnits) SumUnits {
   var ans SumUnits
   for i := 0; i < M; i++ {
       for j := 0; j < M; j++ {
           ans.value[(i+j)%M] += su.value[i] * other.value[j]
       }
   }
   return ans
}

// PowMod raises su to the n-th power under convolution
func (su SumUnits) PowMod(n int) SumUnits {
   var ans SumUnits
   ans.value[0] = 1
   val := su
   for n > 0 {
       if n&1 != 0 {
           ans = ans.Mul(val)
       }
       n >>= 1
       val = val.Mul(val)
   }
   return ans
}

// Shift rotates the values by x mod M
func (su SumUnits) Shift(x int) SumUnits {
   var ans SumUnits
   for i := 0; i < M; i++ {
       ans.value[(i+x)%M] = su.value[i]
   }
   return ans
}

// GetInt extracts the integer count from SumUnits
func (su SumUnits) GetInt() uint64 {
   cpy := su
   // eliminate contributions
   v1 := cpy.value[1]
   for i := 0; i < M; i++ {
       cpy.value[i] -= v1
   }
   v2 := cpy.value[2]
   for i := 0; i < M; i += 2 {
       cpy.value[i] -= v2
   }
   v5 := cpy.value[5]
   for i := 0; i < M; i += 5 {
       cpy.value[i] -= v5
   }
   return cpy.value[0]
}

var values [MaxRange]SumUnits

// DoDFT performs digit-based DFT with base coef, invert flag
func DoDFT(coef int, invert bool) {
   dir := 1
   if invert {
       dir = M - 1
   }
   for base := 0; base < MaxRange; base++ {
       if base%(coef*10) >= coef {
           continue
       }
       var ids [M]int
       for i := 0; i < M; i++ {
           ids[i] = base + i*coef
       }
       var su [M]SumUnits
       for i := 0; i < M; i++ {
           su[i] = SumUnits{}
       }
       for i := 0; i < M; i++ {
           for j := 0; j < M; j++ {
               n := ids[j]
               shifted := values[n].Shift((i*j*dir)%M)
               su[i].Add(&shifted)
           }
       }
       for i := 0; i < M; i++ {
           values[ids[i]] = su[i]
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   fmt.Fscan(in, &N)
   for i := 0; i < N; i++ {
       var v int
       fmt.Fscan(in, &v)
       values[v].value[0]++
   }
   // forward DFT
   for coef := 1; coef < MaxRange; coef *= 10 {
       DoDFT(coef, false)
   }
   // pointwise power
   for i := 0; i < MaxRange; i++ {
       values[i] = values[i].PowMod(N)
   }
   // inverse DFT
   for coef := 1; coef < MaxRange; coef *= 10 {
       DoDFT(coef, true)
   }
   // extract answers
   answers := make([]uint64, N)
   for i := 0; i < N; i++ {
       answers[i] = values[i].GetInt()
   }
   // normalize
   for coef := 1; coef < MaxRange; coef *= 10 {
       for i := range answers {
           answers[i] /= 2
           answers[i] *= Inv5
       }
   }
   // output
   for _, x := range answers {
       fmt.Fprintln(out, x%Mod)
   }
}
