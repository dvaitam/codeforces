#include <iostream>
#include<algorithm>
#include<vector>
#include<string>
#include<set>
#include<deque>
#include<map>
#include<stack>
using namespace std;

#define endl '\n';
#define yes std::cout<<"YES\n";
#define no std::cout<<"NO\n";

using ll = long long;
using ull = unsigned long long;
using db = double;

const int inf = 1e7;
const int maxn = 2000;
const int mod = 1e9 + 7;



void solve() {
	int n; cin >> n;
	vector<bool> s(n);
	vector<int>direct(n);
	vector<int>deg(n, 0);
	stack<int> lst;
	vector<int> ans;
	vector<bool> used(n, false);

	for (int i = 0; i < n; ++i) {
		char c; cin >> c;
		if (c == '1') s[i] = true;
		else s[i] = false;
	}

	for (int i = 0; i < n; ++i) {
		int k; cin >> k; --k;
		direct[i] = k;
		++deg[k];
	}

	for (int i = 0; i < n; ++i) {
		if (deg[i] == 0)lst.push(i);
	}

	while (!lst.empty()) {
		int f = lst.top(); lst.pop();
		int x = direct[f];
		used[f] = true;
		if (s[f]) {
			ans.push_back(f);
			s[f] = 0;
			s[x] = !s[x];
		}
		--deg[x];
		if (deg[x] == 0)lst.push(x);
	}


	auto s1 = s;
	for (int i = 0; i < n; ++i) {
		if (used[i]) continue;

		int f = i, t = 0, cur = i, e = -1;
		while (true) {
			used[cur] = true;
			int x = direct[cur];
			if (s[cur] and e == -1) e = cur;
			if (s[cur])++t;
			if (x == f) break;
			cur = x;
		}
		f = e;


		if (t % 2 == 1) { std::cout << "-1\n"; return; }
		if (f == -1) continue;


		vector<int> ans1;
		cur = f;
		while (true) {
			int x = direct[cur];
			if (s[cur]) {
				s[cur] = !s[cur]; s[x] = !s[x];
				ans1.push_back(cur);
			}
			cur = x;
			if (cur == f) break;
		}


		f = direct[f]; cur = f;
		vector<int> ans2;
		while (true) {
			int x = direct[cur];
			if (s1[cur]) {
				s1[cur] = !s1[cur]; s1[x] = !s1[x];
				ans2.push_back(cur);
			}
			cur = x;
			if (cur == f) break;
		}


		if (ans1.size() < ans2.size()) for (int i : ans1) ans.push_back(i);
		else for (int i : ans2) ans.push_back(i);

	}
	std::cout << ans.size() << endl;
	for (auto& i : ans) std::cout << i+1 << ' ';
	std::cout << endl;
}


int main()
{
	ios::sync_with_stdio(0);
	cin.tie(0);
	int t; cin >> t; while (t--) solve();
}