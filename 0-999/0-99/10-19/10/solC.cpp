#include <bits/stdc++.h>
#include <cstdio>
using namespace std;

int main() {
	int n;
	long long tot = 0;
	long long c[10] = {0};

	scanf("%d", &n);
	for (int i=1; i<=n; ++i) tot -= n / i;
	for (int i=0; i<9; ++i) c[i] = n / 9 + (n % 9 >= i); --c[0];
	for (int i=0; i<9; ++i) for (int j=0; j<9; ++j) for (int k=0; k<9; ++k) {
		if ((i * j) % 9 == k) tot += c[i] * c[j] * c[k];
	}
	printf("%lld\n", tot);

	return 0;
}