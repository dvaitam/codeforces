#include <bits/stdc++.h>

typedef long long ll;

#define forn(i, n) for (int i = 1; i <= n; i++)
#define pb push_back
#define all(x) x.begin(), x.end()
#define y1 qewr1234

using namespace std;

void solve()
{
	string s;
	cin >> s;
	int cnt1 = 0;
	int cnt2 = 0;
	for (int i = 0; i < (int)s.size(); i++)
	{
		if (s[i] >= 'A' && s[i] <= 'Z')
			cnt1++;
		else if (s[i] >= '0' && s[i] <= '9')
			cnt2++;	
	}
	int cnt3 = (int)s.size() - cnt1 - cnt2;
	if (cnt1 > 0 && cnt2 > 0 && cnt3 > 0)
	{
		cout << s << endl;
		return;
	}
	for (int i = 0; i < (int)s.size(); i++)
	{
		cnt1 -= bool(s[i] >= 'A' && s[i] <= 'Z');
		cnt2 -= bool(s[i] >= '0' && s[i] <= '9');
		cnt3 -= bool(s[i] >= 'a' && s[i] <= 'z');
		if (bool(cnt1 == 0) + bool(cnt2 == 0) + bool(cnt3 == 0) <= 1)
		{
			if (cnt1 == 0)
				s[i] = 'A';
			if (cnt2 == 0)
				s[i] = '4';
			if (cnt3 == 0)
				s[i] = 'a';
			cout << s << endl;
			return;
		}
		cnt1 += bool(s[i] >= 'A' && s[i] <= 'Z');
		cnt2 += bool(s[i] >= '0' && s[i] <= '9');
		cnt3 += bool(s[i] >= 'a' && s[i] <= 'z');
	}
	for (int i = 0; i < (int)s.size() - 1; i++)
	{
		cnt1 -= bool(s[i] >= 'A' && s[i] <= 'Z');
		cnt2 -= bool(s[i] >= '0' && s[i] <= '9');
		cnt3 -= bool(s[i] >= 'a' && s[i] <= 'z');
		cnt1 -= bool(s[i + 1] >= 'A' && s[i + 1] <= 'Z');
		cnt2 -= bool(s[i + 1] >= '0' && s[i + 1] <= '9');
		cnt3 -= bool(s[i + 1] >= 'a' && s[i + 1] <= 'z');
		if (bool(cnt1 == 0) + bool(cnt2 == 0) + bool(cnt3 == 0) <= 2)
		{
			int pos = i;
			if (cnt1 == 0)
				s[pos++] = 'A';
			if (cnt2 == 0)
				s[pos++] = '4';
			if (cnt3 == 0)
				s[pos++] = 'a';
			cout << s << endl;
			return;
		}
		cnt1 += bool(s[i] >= 'A' && s[i] <= 'Z');
		cnt2 += bool(s[i] >= '0' && s[i] <= '9');
		cnt3 += bool(s[i] >= 'a' && s[i] <= 'z');
		cnt1 += bool(s[i + 1] >= 'A' && s[i + 1] <= 'Z');
		cnt2 += bool(s[i + 1] >= '0' && s[i + 1] <= '9');
		cnt3 += bool(s[i + 1] >= 'a' && s[i + 1] <= 'z');		
	}
	s[0] = 'a';
	s[1] = 'A';
	s[3] = '4';
	cout << s << endl;
	return;
}

int main()
{
	ios_base::sync_with_stdio(false);
	cin.tie();
	cout.tie();		

	//freopen(".in", "r", stdin);
	//freopen(".out", "w", stdout);

	int T;
	cin >> T;
	while(T--)
		solve();
	
	return 0;
}