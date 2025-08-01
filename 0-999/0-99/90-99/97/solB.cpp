//In the name of God
#include <algorithm>
#include <iostream>
#include <vector>
using namespace std;

typedef pair<int, int> pii;
#define X first
#define Y second

int n;
vector<pii> v;

void solve(int b, int e) {
	if (e - b <= 1)
		return;
	int m = b + e >> 1;
	solve(b, m);
	solve(m, e);
	for (int i = b; i < e; i++)
		v.push_back(pii(v[m].X, v[i].Y));
}
int main() {
	ios_base::sync_with_stdio(false);
	cin >> n;
	pii p;
	for (int i = 0; i < n; i++)
		cin >> p.X >> p.Y, v.push_back(p);
	sort(v.begin(), v.end());
	solve(0, n);
	sort(v.begin(), v.end());
	v.resize(unique(v.begin(), v.end()) - v.begin());
	cout << v.size() << '\n';
	for (int i = 0; i < v.size(); i++)
		cout << v[i].X << ' ' << v[i].Y << '\n';
	return 0;
}