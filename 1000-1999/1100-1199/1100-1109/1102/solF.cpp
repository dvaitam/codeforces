/**
 *    author:  NeverBeNutella
 *    created: Jan 9, 2019 11:13:59 PM
**/

#include <bits/stdc++.h>

using namespace std;

const int inf = 1 << 30;

int n, m;
int a[20][10010], minimum[20][20];
int dp[16][16][1 << 16];

int calc(const int fi, const int pre, const int d) {
	if (d == (1 << n) - 1) {
    	int r = inf;
		for (int i = 0; i < m - 1; i++) {
			r = min(r, abs(a[pre][i]- a[fi][i + 1]));
		}
    	return r;
	}
	if (dp[fi][pre][d] != -1) {
		return dp[fi][pre][d];
	}
	int bt = 0;
	for (int i = 0; i < 16; i++) {
		if (!(d & (1 << i))) {
    		bt = max(bt, min(minimum[pre][i], calc(fi, i, d | (1 << i))));
		}
	}
	return dp[fi][pre][d] = bt;
}

int main() {
	ios::sync_with_stdio(false);
	cin.tie(nullptr);
	cin >> n >> m;
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < m; j++) {
	    	cin >> a[i][j];
		}
	}
	for (int i = 0; i < n; i++) {
	    for (int j = i + 1; j < n; j++) {
	        int tmp = inf;
	        for (int k = 0; k < m; k++) {
	        	tmp = min(tmp, abs(a[i][k] - a[j][k]));
			}
	        minimum[i][j] = minimum[j][i] = tmp;
	    }
	}
	memset(dp, -1, sizeof(dp));
	int res = 0;
	for (int i = 0; i < n; i++) {
		res = max(res, calc(i, i, 1 << i));
	}
	cout << res << '\n';
	return 0;
}