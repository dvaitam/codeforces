#include <bits/stdc++.h>
using namespace std;

typedef long long ll;

const int maxn = 505;

int n, cnt, tot, now;
int d[maxn], used[maxn];
pair<int, int> p[maxn], q[maxn];
std::vector<pair<int, int> > ans;

void ins(int u, int v)
{
	ans.push_back(make_pair(u, v));
	++used[u], ++used[v];
}

int main()
{
	scanf("%d", &n);
	for(int i = 1; i <= n; ++i)
	{
		scanf("%d", &d[i]);
		p[i].first = d[i];
		p[i].second = i;
	 	if(p[i].first == 1) ++cnt;
	 	else tot += p[i].first;
	}
	if(cnt == n)
	{
		if(n <= 2)
		{
			printf("YES %d\n%d\n", n - 1, n - 1);
			if(n > 1) puts("1 2");
		}
		else puts("NO");
		return 0;
	}
	tot -= (n - cnt - 1) * 2;
	sort(p + 1, p + 1 + n);
	int len = n - cnt - 1;
	if(tot < cnt) {puts("NO"); return 0;}

	for(int i = cnt + 1; i < n; ++i)
		ins(p[i].second, p[i + 1].second);
	if(cnt > 0)
		ins(p[1].second, p[cnt + 1].second);
	if(cnt > 1)
		ins(p[2].second, p[n].second);
	for(int i = 3, j = n; i <= cnt; ++i)
	{
		while(j > cnt && used[p[j].second] == d[p[j].second]) --j;
		ins(p[i].second, p[j].second);
	}
	if(used[p[n].second] >= 2) ++len;
	if(used[p[cnt + 1].second] >= 2) ++len;
	printf("YES %d\n", len);
	printf("%d\n", (int)ans.size());
	for(auto x: ans)
		printf("%d %d\n", x.first, x.second);
	return 0;
}