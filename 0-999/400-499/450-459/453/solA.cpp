#include <bits/stdc++.h>



using namespace std;



#define D double



int n, m;

D p, ans;



int main() {

	cin >> m >> n;

	ans = m;

	for (int i = 1; i <= m; i++) {

		p = D(i - 1) / m;

		ans -= pow(p, n);

	}

	cout << ans;

	return 0;

}