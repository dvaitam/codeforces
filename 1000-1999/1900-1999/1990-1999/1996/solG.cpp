#include <bits/stdc++.h>
#define l long long
using namespace std;

mt19937_64 g(chrono::steady_clock::now().time_since_epoch().count());

int main() {
	l n,m,a,b,k,t; 
	cin >> t;
	while (t--) {
		cin >> n >> m;
		vector<l> v(n); 
		while (m--) {
			k=g(); cin >> a >> b;
			v[a-1]^=k;
			v[b-1]^=k;
		}
		map<l,l> c;
		for (l r:v) m=max(m,++c[a^=r]); 
		cout << n-m << "\n"; 
	}
}