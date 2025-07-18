#include <iostream>
#include <algorithm>
#include <cmath>
#include <vector>
#include <string>
#include <queue>
#define MOD 1000000007
#define LL long long
using namespace std;

int n,a[15],sum;

int main()
{
	cin >> n;
	for (int i = 0; i < n; i++)
	{
		cin >> a[i];
		sum += a[i];
	}
	if (n == 1 || (n == 2 && a[0] == a[1]))
		cout << "-1" << endl;
	else
	{
		if (sum - a[0] != a[0])
			cout << 1 << endl << 1 << endl;
		else
			cout << 2 << endl << 1 <<" "<< 2 << endl;
	}
	return 0;
}