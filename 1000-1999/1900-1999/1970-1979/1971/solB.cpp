#include <bits/stdc++.h>
using namespace std;
#define int long long
signed main(){
	cin.tie(nullptr)->sync_with_stdio(false);
	int t;
	cin >> t;
	while(t--){
		string s;
		cin >> s;
		string a = s;
		sort(a.begin(),a.end());
		string b = a;
		reverse(b.begin(),b.end());
		if(a==s && b==s){
			cout << "NO\n";
			continue;
		}
		cout << "YES\n";
		if(a==s)cout << b << '\n';
		else cout << a << '\n';
	}
	return 0;
}