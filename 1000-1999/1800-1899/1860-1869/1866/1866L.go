package main

import (
	"fmt"
	"math/big"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	var nVal, mVal int64
	if _, err := fmt.Scan(&nVal, &mVal); err != nil {
		return
	}

	maxBalls := new(big.Int).SetInt64(-1)
	bestK := int64(1)

	bigN := big.NewInt(nVal)
	
	// Reusable big ints to avoid allocation in loops
	currentSum := new(big.Int)
	term := new(big.Int)
	temp1 := new(big.Int)
	temp2 := new(big.Int)
	temp3 := new(big.Int)
	bigKPrime := new(big.Int)
	bigNPrime := new(big.Int)
	bigD := new(big.Int)

	for K := int64(1); K <= mVal; K++ {
		d := gcd(K, nVal)
		nPrime := nVal / d
		kPrime := K / d
		
		currentSum.SetInt64(0)
		// If d > 1, Lihmuf gets box N at turn N' (since N > N')
		// This corresponds to the y=0 (mod N') case in the cycle
		if d > 1 {
			currentSum.Set(bigN)
		}

		bigKPrime.SetInt64(kPrime)
		bigNPrime.SetInt64(nPrime)
		bigD.SetInt64(d)

		// Iterate through segments where floor(y * K' / N') is constant k
		for k := int64(0); k < kPrime; k++ {
			// Determine valid range for y in [1, N'-1] such that floor(y * K' / N') == k
			// k <= y * K' / N' < k + 1
			// ceil(k * N' / K') <= y <= floor(((k + 1) * N' - 1) / K')
			
			minY := (k*nPrime + kPrime - 1) / kPrime
			maxY := ((k+1)*nPrime - 1) / kPrime
			
			if minY < 1 {
				minY = 1
			}
			if maxY >= nPrime {
				maxY = nPrime - 1
			}
			
			if minY > maxY {
				continue
			}
			
			// Condition for Lihmuf to get the balls: BoxValue > TurnNumber
			// BoxValue = (y * K' - k * N') * d
			// TurnNumber = y
			// (y * K' - k * N') * d > y
			// y * (K' * d - 1) > k * N' * d
			
			denom := kPrime*d - 1
			var threshold int64
			if denom <= 0 {
				// If K'd - 1 <= 0, then K=1, d=1 => 0 > k*N' (impossible for k>=0, N'>=1)
				threshold = maxY + 2
			} else {
				// y > (k * N' * d) / denom
				// y >= floor(...) + 1
				num := k * nPrime * d
				threshold = num/denom + 1
			}
			
			startY := minY
			if threshold > startY {
				startY = threshold
			}
			endY := maxY
			
			if startY <= endY {
				// Sum the box values for y in [startY, endY]
				// Sum = sum over y of ((y * K' - k * N') * d)
				//     = d * (K' * sum(y) - count * k * N')
				
				count := endY - startY + 1
				sumYNum := startY + endY
				
				// temp1 = sumYNum * count / 2  (Sum of arithmetic progression startY...endY)
				temp1.SetInt64(sumYNum)
				temp2.SetInt64(count)
				temp1.Mul(temp1, temp2)
				temp1.Rsh(temp1, 1) // Divide by 2
				
				// temp1 = temp1 * kPrime
				temp1.Mul(temp1, bigKPrime)
				
				// temp2 = count * k * nPrime
				temp2.SetInt64(count)
				temp3.SetInt64(k)
				temp2.Mul(temp2, temp3)
				temp2.Mul(temp2, bigNPrime)
				
				// temp1 = temp1 - temp2
				temp1.Sub(temp1, temp2)
				
				// term = temp1 * d
				term.Mul(temp1, bigD)
				
				currentSum.Add(currentSum, term)
			}
		}
		
		// Update maxBalls. Since we iterate K from 1 to M, strictly greater check
		// ensures we keep the minimum K for the same maximum score.
		if currentSum.Cmp(maxBalls) > 0 {
			maxBalls.Set(currentSum)
			bestK = K
		}
	}
	
	fmt.Println(bestK)
}