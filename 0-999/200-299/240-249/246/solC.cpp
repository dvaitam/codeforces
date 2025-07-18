#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef unsigned long long ul;
#define sz(x) ((int)(x).size())
#define rep(i,l,r) for(int i=(l),I=(r);i<I;++i)
//-------head-------
const int N = 55;
int n, k, a[N];
void solve() {
	rep(i, 0, n + 1)
		rep(j, 1, n - i + 1) {
			if (k <= 0)
				return ;
			--k;
			printf("%d", i + 1);
			rep(p, n - i + 1, n + 1) printf(" %d", a[p]);
			printf(" %d\n", a[j]);
		}
}
int main() {
	scanf("%d%d", &n, &k);
	rep(i, 1, n + 1)
		scanf("%d", &a[i]);
	sort(a + 1, a + 1 + n);
	solve();
	return 0;
}