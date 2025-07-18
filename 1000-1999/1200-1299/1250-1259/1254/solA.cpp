#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

#define pb push_back
#define ppp pop_back
#define pii pair<int,int>
#define fi first
#define se second

const int N = 105;

char a[N][N];
int cnt[N];

char get(int x) {
	if (x < 26) return 'a' + x;
	else if (x - 26 < 26) return 'A' + x - 26;
	else return '0' + x - 52;
}

int main() {
    ios_base::sync_with_stdio(0); cin.tie(0);
#ifdef LOCAL
    freopen("input.txt", "r", stdin);
#endif
	int T;
	cin >> T;
	while (T--) {
		int n, m, k;
		cin >> n >> m >> k;
		int r = 0;
		for (int i=1; i<=n; i++) {
			for (int j=1; j<=m; j++) {
				cin >> a[i][j];
				if (a[i][j] == 'R') r++;
			}
		}
		for (int i=0; i<k; i++) cnt[i] = 0;
		for (int i=0; i<r; i++) cnt[i % k]++;
		int chi = 0;
		for (int i=1; i<=n; i++) {
			vector<int> cur;
			if (i & 1) {
				for (int j=1; j<=m; j++) cur.pb(j);
			}
			else {
				for (int j=m; j>=1; j--) cur.pb(j);
			}
			for (int j : cur) {
				while (cnt[chi] == 0 && chi < k - 1) chi++;
				if (a[i][j] == 'R') cnt[chi]--;
				a[i][j] = get(chi);
			}
		}
		for (int i=1; i<=n; i++) {
			for (int j=1; j<=m; j++) {
				cout << a[i][j];
			}
			cout << "\n";
		}
	}
	return 0;
}