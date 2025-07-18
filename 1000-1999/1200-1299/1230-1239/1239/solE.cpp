#include <bits/stdc++.h>
using namespace std;

bitset<55000 * 25> chk[55][28];
int a[1100];

int main() {
	ios::sync_with_stdio(0);
	cin.tie(0); cout.tie(0);
	int n; cin >> n;
	for (int i = 0; i < 2 * n; i++) cin >> a[i];
	sort(a, a + 2 * n);
	chk[1][0][0] = 1;
	for (int i = 2; i < 2 * n; i++) {
		for (int j = 0; j <= min(i - 1, n - 1); j++) {
			chk[i][j] |= chk[i - 1][j];
			if (j) chk[i][j] |= (chk[i - 1][j - 1] << a[i]);
		}
	}
	int tot = 0;
	for (int i = 2; i < 2 * n; i++) tot += a[i];
	int sum = -1;
	for (int i = 0; i < 55000 * 25; i++) if (chk[2 * n - 1][n - 1][i]) {
		if (abs(tot - 2 * sum) > abs(tot - 2 * i)) sum = i;
	}
	int gae = n - 1;
	vector<int> v, u;
	for (int i = 2 * n - 1; i >= 2; i--) {
		if (sum - a[i] >= 0 && chk[i - 1][gae - 1][sum - a[i]]) {
			v.push_back(a[i]);
			gae--;
			sum -= a[i];
		} else u.push_back(a[i]);
	}
	v.push_back(a[0]);
	u.push_back(a[1]);
	sort(v.begin(), v.end());
	sort(u.begin(), u.end()); reverse(u.begin(), u.end());
	for (int i : v) cout << i << " "; cout << endl;
	for (int i : u) cout << i << " "; cout << endl;
	return 0;
}