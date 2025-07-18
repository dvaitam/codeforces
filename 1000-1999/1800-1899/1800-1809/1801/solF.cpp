#include <cstdio>
#include <cmath>
#include <vector>

const int N = 1e7 + 10, M = 1e4;
typedef double db;

int n, k, S;
int st[M], cnt, st2[M];
std::vector <std::pair <int,int>> go[M];
// dp[i] * [x/j] / x -> dp[i/j] 
db dp[M];
int id[N];
void cmax(db &a, const db &b){ ((a < b) && (a = b)); }
db rem[N];

int main() {
	scanf("%d%d", &n, &k), k--;
	for(int i = 1, j; i <= k; i = j + 1) st[++cnt] = j = k / (k / i), id[j] = cnt, st2[cnt] = i;
	dp[cnt] = k + 1;
	for(int u = 1; u <= cnt; ++u) {
		int v = st[u], c = 0; go[u].resize((sqrt(v) + 1) * 2);
		for(int i = 2, j; i <= v; i = j + 1) {
			j = v / (v / i), go[u][c++] = std::make_pair(i, id[v / i]);
		}
		while(!go[u].back().first) go[u].pop_back();
	}
	for(int t = 1, x; t <= n; ++t) {
		scanf("%d", &x);
		for(int i = 1; i <= cnt; ++i) rem[st2[i]] = 1. * (x / st2[i]) / x;
		for(int u = 1; u <= cnt; ++u) {
			int v = st[u];
			if(dp[u] < 1e-8) continue;
			for(const auto &[p, to] : go[u]) {
				if(p > x) break;
				cmax(dp[to], dp[u] * rem[p]);
			}
			if(x > v) cmax(dp[0], dp[u] * (x / (v + 1)) / x);
		}
	}
	printf("%.10lf\n", dp[0]);
}