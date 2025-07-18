#include <bits/stdc++.h>

using namespace std;

int main()
{
	int i = 0;
	bool ans = 1;
	int n;
	cin >> n;
	char in;
	for (i; i < n; i++)
	{
		cin >> in;
		if (in == '1'){
			ans >>= 1;
			break;
		}
	}
	for (++i; i < n; ++i) cin >> in;
	
	cout << (ans ? "EASY" : "HARD") << "\n";
	return 0;
}