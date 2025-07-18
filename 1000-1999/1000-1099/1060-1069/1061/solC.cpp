#include <bits/stdc++.h>
using namespace std;

const int M = 1000000007;

int a[100001];
int ff[1000001];
int g[1000001] = {0};

int main()
{
	for (int i = 2; i <= 1000000; i++)
		if (!ff[i])
			for (int j = i; j <= 1000000; j += i)
				if (!ff[j])
					ff[j] = i;
	int n;
	scanf("%d", &n);
	for (int i = 1; i <= n; i++)
		scanf("%d", &a[i]);

	g[0] = 1;
	for (int j = 1; j <= n; j++) {
		vector<int> divisor = {1};
		int x = a[j];
		while (ff[x]) {
			long long gg;
			size_t oldsize = divisor.size();
			for (gg = ff[x]; x % gg == 0; gg *= ff[x])
				for (int i = 0; i < oldsize; i++)
					divisor.push_back(divisor[i] * gg);
			x /= (gg / ff[x]);
		}
		vector<int> value;
		value.resize(divisor.size());
		for (size_t i = 0; i < divisor.size(); i++)
			value[i] = g[divisor[i] - 1];
		for (size_t i = 0; i < divisor.size(); i++) {
			g[divisor[i]] += value[i];
			if (g[divisor[i]] >= M)
				g[divisor[i]] -= M;
		}
	}
	long long ans = 0;
	for (int i = 1; i <= 1000000; i++)
		ans = ans + g[i];
	ans %= M;
	cout << ans << '\n';
	return 0;
}