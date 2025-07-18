#include<bits/stdc++.h>
using namespace std;
char s[200005];
int main()
{
	ios::sync_with_stdio(false);
	cin.tie(0);
	int n;
	cin >> n;
	cin >> s;
	int ans = 0;
	for(int i = 1;i < n;i++)
	{
		if(s[i] == s[i - 1])
		{
			ans++;
			if(i == n - 1)
			{
				if(s[i - 1] == 'R')
					s[i] = 'G';
				else if(s[i - 1] == 'G')
					s[i] = 'B';
				else 
					s[i] = 'R';
			} 
			else
			{
			if(s[i-1] != 'R' && s[i+1] != 'R')
				s[i] = 'R';
			else if(s[i-1] != 'G' && s[i+1] != 'G')
				s[i] = 'G';
			else if(s[i-1] != 'B' && s[i+1] != 'B')
				s[i] = 'B';
			}
		}
	}
	cout << ans << endl;
	cout << s << endl;
	return 0;
}