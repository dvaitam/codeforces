package main

import (
   "fmt"
)

func main() {
   var n, sx, sy, dx, dy, t int64
   if _, err := fmt.Scan(&n, &sx, &sy, &dx, &dy, &t); err != nil {
       return
   }
   if t == 0 {
       fmt.Printf("%d %d", sx, sy)
       return
   }
   // modulus for sum and diff computations = 2*n
   mod := uint64(n * 2)
   // initial 0-indexed positions
   X0 := uint64((sx - 1) % n)
   Y0 := uint64((sy - 1) % n)
   // initial sum S0 = X0+Y0
   S0 := (X0 + Y0) % mod
   // initial speed sum DS0 = dx+dy
   DS0 := int64(dx + dy)
   // compute S_{-1} = S0 - DS0
   Sneg1 := (int64(S0) - DS0) % int64(mod)
   if Sneg1 < 0 {
       Sneg1 += int64(mod)
   }
   // initial vector V0 = [S0, S_{-1}, 0, 1]
   V0 := [4]uint64{uint64(S0), uint64(Sneg1), 0, 1}
   // build transition matrix M
   // M * [S_{i-1}, S_{i-2}, i-1, 1] = [S_i, S_{i-1}, i, 1]
   M := [4][4]uint64{
       {4 % mod, (mod - 1) % mod, 2 % mod, 4 % mod},
       {1, 0, 0, 0},
       {0, 0, 1, 1},
       {0, 0, 0, 1},
   }
   // exponentiate M^t
   Mt := matPow(M, uint64(t), mod)
   // compute Vt = Mt * V0
   var Vt [4]uint64
   for i := 0; i < 4; i++ {
       var acc uint64
       for j := 0; j < 4; j++ {
           acc += Mt[i][j] * V0[j]
       }
       Vt[i] = acc % mod
   }
   // S_t = Vt[0]
   St := Vt[0]
   // compute diff_t = diff0 + t*delta  (mod 2n)
   // diff0 = X0-Y0
   diff0 := int64(X0) - int64(Y0)
   delta := dx - dy
   // mod diff_t
   tmod := uint64(t % int64(mod))
   delta_mod := (delta % int64(mod) + int64(mod)) % int64(mod)
   diff_t := (int64(diff0)%int64(mod) + (int64(tmod)*delta_mod)%int64(mod)) % int64(mod)
   if diff_t < 0 {
       diff_t += int64(mod)
   }
   d := uint64(diff_t)
   // compute X_t = (St + d)/2 mod n
   sumd := (St + d) % mod
   Xt := (sumd / 2) % uint64(n)
   // compute Y_t = (St - d)/2 mod n
   // ensure non-negative before division
   diff := (St + mod - d) % mod
   Yt := (diff / 2) % uint64(n)
   // back to 1-indexed
   fmt.Printf("%d %d", Xt+1, Yt+1)
}

// matMul multiplies two 4x4 matrices under modulo
func matMul(a, b [4][4]uint64, mod uint64) [4][4]uint64 {
   var c [4][4]uint64
   for i := 0; i < 4; i++ {
       for k := 0; k < 4; k++ {
           if a[i][k] == 0 {
               continue
           }
           for j := 0; j < 4; j++ {
               c[i][j] += a[i][k] * b[k][j]
           }
       }
       for j := 0; j < 4; j++ {
           c[i][j] %= mod
       }
   }
   return c
}

// matPow raises matrix to the power p under modulo
func matPow(mat [4][4]uint64, p, mod uint64) [4][4]uint64 {
   // initialize res = identity
   var res [4][4]uint64
   for i := 0; i < 4; i++ {
       res[i][i] = 1 % mod
   }
   for p > 0 {
       if p&1 != 0 {
           res = matMul(res, mat, mod)
       }
       mat = matMul(mat, mat, mod)
       p >>= 1
   }
   return res
}
