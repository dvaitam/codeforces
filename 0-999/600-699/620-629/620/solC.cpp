#include <bits/stdc++.h>

using namespace std;

#define endl '\n'

typedef long long ll;
typedef pair<int, int> pii;

const int N = 3e5 + 10;

int a[N], dp[N], coord[N], last[N];

int main() 
{
	//ios_base::sync_with_stdio(0);
	//cin.tie(0);

	int n;
	//cin >> n;
	scanf("%d", &n);
	for(int i = 1; i <= n; i++)
	{
		scanf("%d", &a[i]);
		coord[i - 1] = a[i];
		//cin >> a[i];
	}

	
	sort(coord, coord + n);
	
	int cnt = 0;
	for(int i = 0, j = 0; i < n; i = j, cnt++)
	{
		coord[cnt] = coord[i];
		//cout << coord[cnt] << " ";
		for(j = i; j < n && coord[j] == coord[i]; j++);
	}	
	//cout << endl;
	for(int i = 1; i <= n; i++)
	{
		dp[i] = dp[i - 1];
		int pos = lower_bound(coord, coord + cnt, a[i]) - coord;
		if(last[pos])
		{
			int ptr = last[pos];
			dp[i] = max(dp[ptr - 1] + 1, dp[i]);
		}
		last[pos] = i;

	}

	
	if(dp[n] == 0)
	{
		printf("-1\n");
		//cout << -1 << endl; 
		return 0;
	}

	printf("%d\n", dp[n]);
	//cout << dp[n] << endl;

	int r = n;

	/*for(auto x : dp)
		cout << x << " ";
	cout << endl;*/

	for(int i = n - 2, j = n - 2; i > -1; i = j - 1)
	{	
		if(dp[j] + 1 != dp[r])
		{
			j = i;
			continue;
		}
		for(j = i; j > -1 && dp[j] + 1 == dp[r]; j--);
		j++;
		//ans.push_back({j + 1, r});
		//cout << j + 1 << " " << r << endl;
		printf("%d %d\n", j + 1, r);
		r = j;
	}



	return 0;
}