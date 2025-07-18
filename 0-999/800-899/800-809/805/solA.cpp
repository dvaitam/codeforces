#define _CRT_SECURE_NO_WARNINGS

#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <string>
#include <cstring>
#include <vector>
#include <set>
#include <map>
#include <algorithm>

using namespace std;

#define mp make_pair;

typedef long long ll;
typedef long double ld;
typedef vector <int> vi;
ll a, b;
ll sa, sb;

int main()
{
	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
	cin >> a >> b;
	if (a == b)
	{
		cout << a << endl;
		return 0;
	}
	if ((max(a, b) - min(a, b)) > 10)
	{
		cout << "2" << endl;
			return 0;
	}
	for (ll i = min(a, b); i <= max(a, b); ++i)
	{
		if (i % 2 == 0)
		{
			sa++;
		}
		if (i % 3 == 0)
		{
			sb++;
		}
	}
	if (sa > sb)
		cout << "2" << endl;
	else
		cout << "3" << endl;

	return 0;
}