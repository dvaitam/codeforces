//In the name  of God
#include <iostream>
using namespace std;

int n, k;
bool adj[101][101];

int main() {
	cin >> k;
	for (n = 1; k; n++) {
		int m = n;
		while (m * (m - 1) / 2 > k)
			m--;
		for (int i = 0; i < m; i++)
			adj[i][n] = adj[n][i] = true;
		k -= m * (m - 1) / 2;
	}
	cout.sync_with_stdio(false);
	cout << n << '\n';
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < n; j++)
			cout << adj[i][j];
		cout << '\n';
	}
	return 0;
}