#include <cstdio>
#include <algorithm>
#include <cstring>
#define long long long
using namespace std;
long a[1001];
long n;
long gcd(long a, long b) {
	return b == 0 ? a : gcd(b, a % b);
}
int main() {
	scanf("%lld", &n);
	long t;
	for(int i=1; i<=n; i++) {
		scanf("%lld", &a[i]);
	}
	t = a[1];
	for(int i=2; i<=n; i++) {
		t = gcd(a[i], t);
	}
	if(t != a[1]) {
		printf("-1\n");
		return 0;
	}
	printf("%lld\n", (n * 2 - 1));
	printf("%lld", a[1]);
	for(int i=2; i<=n; i++) {
		printf(" %lld %lld", a[i], a[1]);
	}
	return 0;
}