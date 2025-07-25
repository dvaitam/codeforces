#include <stdio.h>
#include <math.h>
#include <string.h>
#include <iostream>
#include <algorithm>
#include <vector>
#include <map>
#include <utility>
#include <stack>
#include <queue>
#include <set>
#include <list>
#include <bitset>
#include <array>

using namespace std;

#define fi first	
#define se second
#define long long long
typedef pair<int,int> ii;

int main()
{
	ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
	// freopen("input.in", "r", stdin);

	string s; cin >> s;
	s += s;
	int mx = 0;
	string ans = "";
	for(char c : s)
	{
		if(ans.size() == 0 || ans.back() != c) ans += c;
		else ans = c;
		mx = max(min((int)s.size()/2,(int)ans.size()), mx);
	}
	cout << mx << "\n";
}