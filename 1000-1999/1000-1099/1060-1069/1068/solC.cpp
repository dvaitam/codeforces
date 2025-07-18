#include <bits/stdc++.h>
#pragma GCC optimize("03")
#define ll long long
#define ld long double
#define fi first
#define se second
#define mod 100

using namespace std;

int n, m, x, y, ok[110][110], vf[110], sol[110][110], cnt[110];

int main() {
//	ifstream cin("tst.in");
//	ofstream out("tst.out");
	ios_base::sync_with_stdio(0);
	cin.tie(NULL);
	cin >> n >> m;
	int cnt = 0;
	for (int i = 1; i <= m; i++) {
		cin >> x >> y;
		cnt++;
		vf[x]++;
		vf[y]++;
		sol[x][vf[x]] = cnt;
		sol[y][vf[y]] = cnt;
	}
	for (int i = 1; i <= n; i++) {
		if (vf[i] == 0) {
			vf[i]++;
			sol[i][vf[i]] = ++cnt;
		}
	}
	int rs = 0;
	for (int i = 1; i <= n; i++) {
		cout << vf[i] << '\n';
		for (int j = 1; j <= vf[i]; j++)
			cout << i << ' ' << sol[i][j] << '\n';
	}
	return 0;
}