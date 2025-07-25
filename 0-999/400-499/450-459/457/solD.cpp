#include <iostream>
#include <cstdio>
#include <cmath>
using namespace std;

int n,m,k;
long double fac[100005];

long double C(int a, int b)
{
	return fac[b] - fac[a] - fac[b-a];
}
int main()
{
	cin>>n>>m>>k;
	fac[0] = 0.0;
	for (int i=1; i<=m; i++) fac[i] += fac[i-1] + log(i);
	long double ans = 0.0;
	for (int i=0; i<=n; i++)
		for (int j=0; j<=n; j++) {
			int t=i*n+j*n-i*j;
			if (t>k) break;
			long double tmp = C(k-t, m-t) + C(i,n) + C(j,n) - C(k, m);
			ans += exp(tmp);
			if (ans > 1e99) { printf("1e99\n"); return 0; }
		}
	printf("%.15lf\n", (double)ans);
	
	return 0;
}