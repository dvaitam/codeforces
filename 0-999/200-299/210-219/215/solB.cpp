#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
#define sz(x) ((int)(x).size())
#define rep(i,l,r) for(int i=(l);i<(r);++i)
//-------head-------
const int INF = 1e9 + 7;
int n;
int main() {
	int A, B, r1 = -INF, p1 = -INF, p2 = INF;
	scanf("%d", &n);
	rep(i, 0, n) {
		int x;
		scanf("%d", &x);
		r1 = max(r1, x);
	}
	scanf("%d", &n);
	rep(i, 0, n) {
		int x;
		scanf("%d", &x);
		p1 = max(p1, x);
	}
	scanf("%d", &n);
	rep(i, 0, n) {
		int x;
		scanf("%d", &x);
		p2 = min(p2, x);
	}
	scanf("%d%d", &A, &B);
//	printf("%d %d %d %d %d\n", r1, p1, p2, A, B);
	double r2 = 1.0 * B * p1 * r1 * r1 / (1.0 * A * p2 + 1.0 * B * p1);
	r2 = sqrt(r2);
	printf("%.10f", r2);

	return 0;
}