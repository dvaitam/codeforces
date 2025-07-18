#include<bits/stdc++.h>
using namespace std;
int main()
{
	int n, m;
	string s;
	cin >> n >> m;
	cin >> s;
	bool f = true;
	for (int i = 0; i < s.size() - m; i++)
	{
		if (s[i] == s[i + m] && s[i] != '.')
		{
			f &= true;
			continue;
		}
		f &= false;
		if (s[i] == '.' || s[i + m] == '.')
		{
			if (s[i] == '1')
				s[i + m] = '0';
			else if (s[i] == '0')
				s[i + m] = '1';
			else if (s[i + m] == '1')
				s[i] = '0';
			else if (s[i + m] == '0')
				s[i] = '1';
			else
			{
				s[i] = '0';
				s[i + m] = '1';
			}
		}
	}
	if (f)
	{
		cout << "No" << endl;
		return 0;
	}
	for (int i = s.size() - m; i < n; i++)
	{
		if (s[i] == '.')
			s[i] = '0';
	}
	cout << s << endl;
	return 0;
}