#include<bits/stdc++.h>

using namespace std;

#define FOR(i, s, e) for(int i = (s); i < (e); i++)

#define FOE(i, s, e) for(int i = (s); i <= (e); i++)

#define FOD(i, s, e) for(int i = (s); i >= (e); i--)

#define ll long long

#define mp make_pair

#define pb push_back

#define pii pair<int, int>

#define pff pair<double, double>


int dp[400][400][400];
int A[405], inf = 1000000005;
int n, m;

int main () {
	scanf("%d %d", &n, &m);
	FOR(i, 0, n) scanf("%d", &A[i]);

	FOR(i, 0, n) FOR(j, i, n) {	
		int s = i;
		FOE(k, 0, (j - i)) {
			if (k == 0) dp[i][j][0] = A[j] - A[i];
			else {
				while (s < j - 1 && dp[i][s][k - 1] < A[j] - A[s]) s++;  

				dp[i][j][k] = inf;

				if (s != i) {
					dp[i][j][k] = max(dp[i][s - 1][k - 1], A[j] - A[s - 1]);
				}
				dp[i][j][k] = min(dp[i][j][k], max(dp[i][s][k - 1], A[j] - A[s]));
			}

			//printf("dp[%d][%d][%d] = %d\n", i, j, k, dp[i][j][k]);
		}
	}	

	ll res = 0ll;

	while (m--) {
		int s, t, c, r;
		scanf("%d %d %d %d", &s, &t, &c, &r);
		s--; t--;
		if (r > t - s) r = t - s;
		ll tmp = 1ll * dp[s][t][r] * c;
		res = max(res, tmp);
	}

	printf("%lld\n", res);
	return 0;
}