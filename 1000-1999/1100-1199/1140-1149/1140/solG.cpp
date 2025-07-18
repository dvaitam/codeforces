#include <bits/stdc++.h>

using namespace std;

int n, q, par[300005][19], d[300005];
long long to[300005][19][2][2], at[2][2], bt[2][2], ct[2][2];
long long w[300005];
vector<pair<int, pair<long long, long long>>> v[300005];
priority_queue<pair<long long, int>, vector<pair<long long, int>>, greater<pair<long long, int>>> qq;
int ta, tb, aa, bb;
long long tc, td;

void merge(long long a[2][2], long long b[2][2], long long c[2][2]) {
	for (int i = 0; i < 2; i++)
		for (int j = 0; j < 2; j++)
			c[i][j] = min(a[i][0] + b[0][j], a[i][1] + b[1][j]);
}

void dfs(int a) {
	for (int i = 1; i < 19; i++) {
		merge(to[a][i - 1], to[par[a][i - 1]][i - 1], to[a][i]);
		par[a][i] = par[par[a][i - 1]][i - 1];
	}
	for (auto& i : v[a])
		if (i.first != par[a][0]) {
			par[i.first][0] = a;
			d[i.first] = d[a] + 1;
			// printf("%d->%d %lld\n", a, i.first, i.second.first);
			to[i.first][0][0][0] = min(i.second.first, w[a] + w[i.first] + i.second.second);
			to[i.first][0][1][1] = min(i.second.second, w[a] + w[i.first] + i.second.first);
			to[i.first][0][0][1] = min(i.second.first + w[a], w[i.first] + i.second.second);
			to[i.first][0][1][0] = min(i.second.second + w[a], w[i.first] + i.second.first);
			dfs(i.first);
		}
}

int main() {
	scanf("%d", &n);
	for (int i = 1; i <= n; i++) {
		scanf("%lld", w + i);
		qq.emplace(w[i], i);
	}
	for (int i = 1; i < n; i++) {
		scanf("%d%d%lld%lld", &ta, &tb, &tc, &td);
		v[ta].emplace_back(tb, make_pair(tc, td));
		v[tb].emplace_back(ta, make_pair(tc, td));
	}
	while (!qq.empty()) {
		tie(tc, ta) = qq.top();
		qq.pop();
		if (w[ta] != tc)
			continue;
		for (auto& i : v[ta]) {
			long long nd = w[ta] + i.second.first + i.second.second;
			if (nd < w[i.first]) {
				w[i.first] = nd;
				qq.emplace(nd, i.first);
			}
		}
	}
	dfs(1);
	scanf("%d", &q);
	while (q--) {
		scanf("%d%d", &ta, &tb);
		// printf("ta=%d tb=%d\n", ta, tb);
		if (ta % 2)
			ta = (ta + 1) / 2, aa = 0;
		else
			ta /= 2, aa = 1;
		if (tb % 2)
			tb = (tb + 1) / 2, bb = 0;
		else
			tb /= 2, bb = 1;
		memset(at, 0, sizeof(at));
		memset(bt, 0, sizeof(bt));
		if (d[ta] > d[tb]) {
			swap(ta, tb);
			swap(aa, bb);
		}
		// printf("ta=%d aa=%d tb=%d bb=%d\n", ta, aa, tb, bb);
		at[0][1] = at[1][0] = w[ta];
		bt[0][1] = bt[1][0] = w[tb];
		for (int i = 0; i < 19; i++)
			if ((d[tb] - d[ta]) & (1 << i)) {
				merge(bt, to[tb][i], ct);
				memcpy(bt, ct, sizeof(bt));
				tb = par[tb][i];
			}
		if (ta != tb) {
			for (int i = 18; i >= 0; i--)
				if (par[ta][i] != par[tb][i]) {
					merge(at, to[ta][i], ct);
					memcpy(at, ct, sizeof(at));
					merge(bt, to[tb][i], ct);
					memcpy(bt, ct, sizeof(bt));
					ta = par[ta][i];
					tb = par[tb][i];
				}
			{
				int i = 0;
					merge(at, to[ta][i], ct);
					memcpy(at, ct, sizeof(at));
					merge(bt, to[tb][i], ct);
					memcpy(bt, ct, sizeof(bt));
					ta = par[ta][i];
					tb = par[tb][i];
			}
		}
		// printf("LCA=%d=%d\n", ta, tb);
		// for (int i = 0; i < 2; i++) for (int j = 0; j < 2; j++) printf("%lld ", to[4][0][i][j]); puts("");
		// for (int i = 0; i < 2; i++) for (int j = 0; j < 2; j++) printf("%lld ", at[i][j]); puts("");
		// for (int i = 0; i < 2; i++) for (int j = 0; j < 2; j++) printf("%lld ", bt[i][j]); puts("");
		printf("%lld\n", min(at[aa][0] + bt[bb][0], at[aa][1] + bt[bb][1]));
	}
}