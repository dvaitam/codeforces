#include <stdio.h>
#include <math.h>
#include <algorithm>
using namespace std;

double f(double x, double A, double B, double L) {
	double y = sqrt(L*L - x*x);
	return (A*x+B*y-x*y)/L;
}

int main(void) {
	double A, B, L;
	scanf("%lf%lf%lf", &A, &B, &L);
	if(L <= B) {
		printf("%.8f\n", min(A, L));
		return 0;
	} else if (L <= A) {
		printf("%.8f\n", min(B, L));
		return 0;
	} else {
		double left = 0.0, right = L, ma, mb;
		for(int u=0;u<500;u++){
			ma = (left*3+right)/4.0;
			mb = (left+right*3)/4.0;
			double fa = f(ma, A, B, L);
			double fb = f(mb, A, B, L);
			//printf("(%f %f %f %f)\n", ma, mb, fa, fb);
			if(fa < fb)
				right = mb;
			else
				left = ma;
		}
		double ff = f((left+right)/2.0, A, B, L);
		ff = min(ff, L);
		ff = min(ff, A);
		if (ff < 1e-8)
			printf("My poor head =(\n");
		else
			printf("%.8f\n", ff);
	}
	return 0;
}