#include<bits/stdc++.h>

using namespace std;

void solve(){
	int n, k; cin >> n >> k;
	int cn = 0;
	int d = k - 1;
	vector<int> ans;
	while(d){
		ans.push_back((d + 1) / 2);
		d -= (d + 1) / 2;
	}
	d = k + 1;
	while(ans.size() < 24){
		ans.push_back(d);
		d *= 2;
	}
	ans.push_back(2 * k + 1);
	cout << ans.size() << '\n';
	for(auto x : ans)cout << x << ' ';
	cout << '\n';
}

int main(){
	ios_base::sync_with_stdio(0); cin.tie(0);
	int t; cin >> t;
	while(t--)solve();
}