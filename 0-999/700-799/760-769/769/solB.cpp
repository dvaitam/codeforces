#define _CRT_SECURE_NO_WARNINGS
#include <bits/stdc++.h>

#define pii pair<long long, long long>
#define F first
#define S second
#define pb push_back
#define mk make_pair
#define all(a) a.begin(), a.end()
#define sz(a) (long long)a.size()

using namespace std;

typedef vector<int> vint;
typedef string str;

int u[1000];
vector<pii> a, ans;
queue<pii> q;

main()
{
	ios_base::sync_with_stdio(0); cin.tie(0);
	int n;
	cin >> n;
	pii tmp;
	cin >> tmp.F, tmp.S = 1;
	if (tmp.F == 0)
		return cout << -1, 0;
	q.push(tmp);
	n--;
	a.resize(n);
	for (int i = 0; i < n; i++)
		cin >> a[i].F, a[i].S = i + 2;
	sort(a.rbegin(), a.rend());
	int id = 0;
	while (!q.empty())
	{
		pii v = q.front(); q.pop();
		while (v.F > 0 && id < sz(a))
		{
			if (!u[a[id].S])
			{
				v.F--;
				ans.pb(mk(v.S, a[id].S));
				q.push(a[id]);
				u[v.S] = u[a[id].S] = 1;
				id++;
			}
			else
			    break;
		}
	}
	for (int i = 1; i <= n + 1; i++)
	    if (!u[i])
	        return cout << -1, 0;
	cout << sz(ans) << "\n";
	for (pii x : ans)
		cout << x.F << " " << x.S << "\n";
}