#include<bits/stdc++.h>
using namespace std;
#define ll long long int 
const int sz = 2e5 + 5, mod = 1e9 + 7;
int32_t main() {
    ios_base::sync_with_stdio(0);
	cin.tie(0);
	int n; cin >> n;
	std::vector<pair<int,int>> v(n), ans(n);
	for(int i = 0; i < n; i++) {
		cin >> v[i].first;
		v[i].second = i;
	}
	sort(v.begin(), v.end());
	std::vector<int> dir(n + 1), taken(n + 1), taken_id(n + 1);
	if(v[0].first == 0) {
		ans[v[0].second] = {1, 1};
		dir[v[0].second] = v[0].second + 1;
		taken[1] = 1;
		taken_id[1] = v[0].second;
	} else {
		int i = 0;
		while(v[i].first != v[i + 1].first)
			i++;
		int dis = v[i].first;
		ans[v[i].second] = {1, 1};
		dir[v[i].second] = v[i + 1].second + 1;
		taken[1] = 1;
		taken_id[1] = v[i].second;
		int x = 1 + dis;
		int y = 1;
		if(x > n) {
			y += x - n;
			x = n;
		}
		ans[v[i + 1].second] = {x, y};
		dir[v[i + 1].second] = v[i].second + 1;
		taken[x] = y;	
		taken_id[x] = v[i + 1].second;
	}
	int cur_x = 1;
	for(auto &[a, id]:v) {
		if(dir[id]) continue;
		while(taken[cur_x]) cur_x++;
		int y;
		if(a == 0) {
			y = 1;
			dir[id] = id + 1;
		}
		else {
			if(cur_x - a >= 1) {
				y = taken[cur_x - a];
				dir[id] = taken_id[cur_x - a] + 1;

			} else {
				y = a - cur_x + 2;
				dir[id] = taken_id[1] + 1;
			} 
		}
		taken[cur_x] = y;
		ans[id] = {cur_x, y};
		taken_id[cur_x] = id;
 	} 
 	cout << "YES\n";
 	for(auto &[x, y]:ans) 
 		cout << x << " " << y << "\n";
 	for(int i = 0; i < n; i++) 
 		cout << dir[i] << " \n"[i == n - 1];
	return 0;
}