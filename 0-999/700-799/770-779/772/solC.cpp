#include <iostream>
#include <vector>
#include <algorithm>
#include <map>
using namespace std;

bool bad[200000];
map<int, vector<int>> mp;

int gcd(int a, int b)
{
	while(b)
		a %= b, swap(a, b);
	return a;
}

int gcdex(int a, int b, int & x, int & y)
{
	if(a == 0)
	{
		x = 0; y = 1;
		return b;
	}
	int x1, y1;
	int d = gcdex(b%a, a, x1, y1);
	x = y1 - (b / a) * x1;
	y = x1;
	return d;
}

int inv(int a, int m)
{
	int x, y;
	gcdex(a, m, x, y);
	return (x % m + m) % m;
}

int main()
{
	int n, m; // m - модуль
	scanf("%d%d", &n, &m);

	for(int i=0; i<n; i++)
	{
		int k;
		scanf("%d", &k);
		bad[k] = true;
	}

	for(int j=0; j<m; j++)
		if(!bad[j])
			mp[gcd(m, j)].push_back(j);

	int size = mp.size();

	vector<int> gcds;
	for(auto z : mp)
		gcds.push_back(z.first);

	int cnt = gcds.size();
	vector<int> ans(cnt);
	vector<int> best(cnt, -1);

	for(int i = cnt-1; i>=0; i--)
	{
		for(int j = i+1; j < cnt; j++)
			if(gcds[j] % gcds[i] == 0)
				if(ans[j] > ans[i])
				{
					ans[i] = ans[j];
					best[i] = j;
				}
		ans[i] += mp[gcds[i]].size();
	}

	int best_i = 0;
	for(int i=1; i<cnt; i++)
		if(ans[best_i] < ans[i])
			best_i = i;

	printf("%d\n", ans[best_i]);

	int prev_k = 1, prev_gcd = 1;

	while(best_i != -1)
	{
		int gcdi = gcds[best_i];
		for(int k : mp[gcdi])
		{
			printf("%d ", 1ll * (k/prev_gcd) * inv(prev_k/prev_gcd, m/prev_gcd) % m);

			prev_k = k;
			prev_gcd = gcdi;
		}
		best_i = best[best_i];
	}
}