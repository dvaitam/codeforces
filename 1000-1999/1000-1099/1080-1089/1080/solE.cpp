#include <iostream>
#include <algorithm>
#include <cstring>
const int max = 250;
char ch[max];
bool cnt[26];
int H[max][max][max],base[26];
const int mod = 1E9 + 9;
int idx = mod,S[max << 1 | 3],p[(max | 1) << 1];
main()
{
	* base = 233333;
	for (int i = 1;i < 26;++ i)
		base[i] = static_cast <long long>(base[i - 1]) * * base % mod;
	std::ios_base::sync_with_stdio(0),std::cin.tie(0);
	int n,m;
	std::cin >> n >> m;
	for (int i = 0;i < n;++ i)
	{
		std::for_each(ch,ch + m,[](char & ch){std::cin >> ch,ch -= 97;});
		for (int L = 0;L < m;++ L)
		{
			std::memset(cnt,0,sizeof cnt);
			int odd = 0,cur = 0,* H = ::H[i][L];
			for (int R = L;R < m;++ R)
			{
				if (cnt[ch[R]] ^= 1)
					++ odd;
				else
					-- odd;
				(cur += base[ch[R]]) %= mod;
				H[R] = odd <= 1 ? cur : idx ++;
			}
		}
	}
	* S = idx ++,S[n << 1 | 1] = S[n << 1 | 2] = - 1;
	unsigned ans = 0;
	for (int L = 0;L < m;++ L)
		for (int R = L;R < m;++ R)
		{
			for (int i = 0,j = 0;i < n;++ i)
				S[++ j] = - 1,S[++ j] = H[i][L][R];
			std::memset(p,0,sizeof p);
			int C,M = 0;
			for (int i = 2;i <= n << 1;++ i)
				if (S[i] < mod)
				{
					int & p = ::p[i];
					if (i < M)
						p = std::min(::p[(C << 1) - i],M - i);
					while (S[i - p] == S[i + p])
						++ p;
					if (i + p > M)
						M = (C = i) + p;
					ans += p >> 1;
				}
		}
	std::cout << ans << '\n';
}