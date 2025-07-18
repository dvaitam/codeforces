#include<bits/stdc++.h>
using namespace std;
typedef long long LL;
typedef pair<int, int> PII;
typedef pair<LL, LL> PLL;

const int MAXN = 2e5 + 10;
const int INF = 0x3f3f3f3f;
const LL INFLL = 0x3f3f3f3f3f3f3f3f;

mt19937 rnd(114514);

void solve()
{
	LL n, s1, s2;
	cin >> n >> s1 >> s2;

	vector<PLL> r(n);
	for(LL i = 0; i < n; ++i)
	{
		auto &[v, id] = r[i];
		cin >> v;
		id = i + 1;
	}

	auto cmp = [&](PLL a, PLL b)
	{
		return a.first > b.first;
	};

	sort(r.begin(), r.end(), cmp);

	vector<LL> res1, res2;
	for(LL i = 0; i < n; ++i)
	{
		auto [v, id] = r[i];
		LL cost1 = v * (res1.size() + 1) * s1;
		LL cost2 = v * (res2.size() + 1) * s2;

		if(cost1 < cost2) res1.push_back(id);
		else if(cost1 > cost2) res2.push_back(id);
		else
		{
			if(s1 > s2) res2.push_back(id);
			else res1.push_back(id);
		}
	}

	cout << res1.size() << ' ';
	for(LL i : res1) cout << i << ' ';
	cout << '\n';

	cout << res2.size() << ' ';
	for(LL i : res2) cout << i << ' ';
	cout << '\n';
}

int main()
{
	ios::sync_with_stdio(false);
	cin.tie(0); cout.tie(0);

	int Task = 1;
	cin >> Task;
	while(Task--)
	{
		solve();
	}

	return 0;
}