#include <iostream>
#include <cstring>
#include <cstdio>
#include <algorithm>
#include <cmath>
using namespace std;

const int maxn = 100000 + 123;

int n, q;
int h[maxn];

long long query[maxn][3];

int disc[maxn * 2];
int tod;

int nums[maxn * 2];
long long summ[maxn * 2];

int lowbit(int x)
{
	return x&(x^(x-1));    
}
void cal(int n, int &nu, long long &su)
{
    su = nu = 0;
    while(n>0)
    {
		su += summ[n];
		nu += nums[n];
		n=n-lowbit(n);
	}
}

void change(int i,int x, const int& n)
{
	long long tsum = disc[i] * x;
	while(i<=n)
	{
		nums[i] += x;
		summ[i] += tsum; 
        i=i+lowbit(i);
    }
}

int main()
{
	//freopen(data.in, r, stdin);
	//freopen(data.out, w, stdout);

	scanf("%d%d", &n, &q);

	for (int i = 1; i <= n; ++ i)
	{
		scanf("%d", h + i);
		disc[i] = h[i];
	}
	tod = n;
	for (int i = 0 ; i < q; ++ i)
	{
		scanf("%I64d", &query[i][0]);
		if (query[i][0] == 1)
		{
			scanf("%I64d%I64d", &query[i][1], &query[i][2]);
			disc[++ tod] = query[i][2];
		}
		else 
			scanf("%I64d", &query[i][1]);
	}

	sort(disc + 1, disc + tod + 1);
	tod = unique(disc + 1, disc + tod + 1) - (disc + 1);

	memset(nums, 0, sizeof(nums));
	memset(summ, 0, sizeof(summ));

	for (int i = 1; i <= n; ++ i)
	{
		int tid = lower_bound(disc + 1, disc + tod + 1, h[i]) - disc;
		//printf("%d %d %d\n", h[i], tid, disc[tid]);
		change(tid, 1, tod);
	}

	for (int i = 0; i < q; ++ i)
	{
		if (query[i][0] == 1)
		{
			int ni = query[i][1];
			int tid = lower_bound(disc + 1, disc + tod + 1, h[ni]) - disc;
			change(tid, -1, tod);
			tid = lower_bound(disc + 1, disc + tod + 1, query[i][2]) - disc;
			change(tid, 1, tod);
			h[ni] = query[i][2];
		}
		else
		{
			long long v = query[i][1];
			int l = 1, r = tod;
			int cnt;
			long long pref;
			int nl = 1;
			while (l <= r)
			{
				int mid = (l + r) >> 1;
				cal(mid, cnt, pref);
				long long maxv = (v + 0.0) + (pref + 0.0);
				if (maxv > (long long)cnt * disc[mid])
				{
					nl = mid;
					l = mid + 1;
				}
				else r = mid - 1;
			}
			cal(nl, cnt, pref);
			printf("%lf\n",	(v + 0.0) / cnt + (pref + 0.0) / cnt);
		}
	}
	return 0;
}