#include <bits/stdc++.h>
#define LL long long
using namespace std;
const int N = 2e6 + 10;
int n, K, cnt[N]; char S[N];
vector<int> opt[N]; bool DP[2][N * 2];
int main() {


	ios :: sync_with_stdio(0); cin.tie(0); cout.tie(0);
	int _; cin >> _;
	while (_ --) {
		cin >> n >> K >> (S + 1);
		for (int i = 0; i <= K; i ++) cnt[i] = 0, opt[i].resize(n * 2 + 1), opt[i].shrink_to_fit();
		for (int i = 1; i <= n; i ++) for (int j = 1; j <= K; j ++) {
			char c; cin >> c; if (c == '1') cnt[K - j + 1] ++;
		}
		for (int i = 0; i <= n * 2; i ++) DP[0][i] = !i, DP[1][i] = 0;
		for (int i = 0; i < K; i ++) {
			for (int j = 0; j <= n * 2; j ++) if (DP[0][j]) {
				DP[1][j / 2 + n - cnt[i + 1]] = 1, opt[i + 1][j / 2 + n - cnt[i + 1]] = j;
				DP[1][j / 2 + cnt[i + 1]] = 1, opt[i + 1][j / 2 + cnt[i + 1]] = j;
			}
			for (int j = 0; j <= n * 2; j ++) DP[0][j] = DP[1][j], DP[1][j] = 0;
			for (int j = 0; j <= n * 2; j ++) if ((j & 1) != (S[K - i] & 1)) DP[0][j] = 0;
		} 
		if (!DP[0][S[1] & 1]) { cout << "-1\n"; continue; }
		for (int i = K, j = (S[1] & 1); i >= 1; i --) {
			int nxt = opt[i][j]; cout << (nxt / 2 + n - cnt[i] == j); j = nxt;
		} cout << "\n";
	}
	return 0;
}