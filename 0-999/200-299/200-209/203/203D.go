package main

import (
	"fmt"
	"math"
)

type point struct {
	x, y, z float64
}

func abs(a float64) float64 {
	return math.Abs(a)
}

func main() {
	var a, b, m int
	var vx, vy, vz float64
	if _, err := fmt.Scan(&a, &b, &m, &vx, &vy, &vz); err != nil {
		return
	}
	aF := float64(a)
	bF := float64(b)
	eps := 1e-9
	INF := math.Inf(1)

	hitTime := make([]float64, 5)
	hitPoints := make([]point, 5)
	hitVectors := make([]point, 5)

	curPoint := point{x: aF * 0.5, y: float64(m), z: 0}
	curVector := point{x: vx, y: vy, z: vz}

	for {
		// time to hit ground y=0
		if abs(curVector.y) < eps {
			hitTime[0] = INF
		} else {
			t := -curPoint.y / curVector.y
			hitTime[0] = t
			hitPoints[0] = point{
				x: curPoint.x + t*curVector.x,
				y: 0,
				z: curPoint.z + t*curVector.z,
			}
		}
		// time to hit z planes
		if abs(curVector.z) < eps {
			hitTime[1], hitTime[2] = INF, INF
		} else {
			t1 := (bF - curPoint.z) / curVector.z
			t2 := -curPoint.z / curVector.z
			hitTime[1], hitTime[2] = t1, t2
			// reflect z
			v1 := curVector
			v2 := curVector
			v1.z = -v1.z
			v2.z = -v2.z
			hitVectors[1], hitVectors[2] = v1, v2
			hitPoints[1] = point{
				x: curPoint.x + t1*curVector.x,
				y: curPoint.y + t1*curVector.y,
				z: bF,
			}
			hitPoints[2] = point{
				x: curPoint.x + t2*curVector.x,
				y: curPoint.y + t2*curVector.y,
				z: 0,
			}
		}
		// time to hit x planes
		if abs(curVector.x) < eps {
			hitTime[3], hitTime[4] = INF, INF
		} else {
			t3 := (aF - curPoint.x) / curVector.x
			t4 := -curPoint.x / curVector.x
			hitTime[3], hitTime[4] = t3, t4
			// reflect x
			v3 := curVector
			v4 := curVector
			v3.x = -v3.x
			v4.x = -v4.x
			hitVectors[3], hitVectors[4] = v3, v4
			hitPoints[3] = point{
				x: aF,
				y: curPoint.y + t3*curVector.y,
				z: curPoint.z + t3*curVector.z,
			}
			hitPoints[4] = point{
				x: 0,
				y: curPoint.y + t4*curVector.y,
				z: curPoint.z + t4*curVector.z,
			}
		}
		// find next hit
		minind := 0
		for i := 1; i < 5; i++ {
			if hitTime[i] < hitTime[minind] && hitTime[i] > eps {
				minind = i
			}
		}
		if minind == 0 {
			curPoint = hitPoints[0]
			break
		}
		curPoint = hitPoints[minind]
		curVector = hitVectors[minind]
	}
	fmt.Printf("%.6f %.6f", curPoint.x, curPoint.z)
}
