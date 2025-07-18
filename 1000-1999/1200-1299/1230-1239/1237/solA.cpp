#include <iostream>
#include <vector>
#include <random>
#include <set>
#include <cmath>
#include <queue>

using namespace std;

random_device rd;
mt19937 rnd(rd());
uniform_int_distribution<int> nextRand(0, (int)1e9);

vector <int> A;
int n;

int main()
{
	ios_base::sync_with_stdio(false);
	cin.tie(nullptr), cout.tie(nullptr);
	cin >> n;
	int a;
	bool f1 = true, f2 = true;
	for (int i = 0; i < n; i++)
	{
		cin >> a;
		if ((a > 0) && (a % 2 != 0))
		{
			
			if (f1)
				a++;
			f1 = !(f1);
		}
		if ((a < 0) && (a % 2 != 0))
		{
			if (f2)
				a--;
			f2 = !(f2);
		}
		cout << a / 2 << '\n';
	}
	ios_base::sync_with_stdio(false);
	cin.tie(nullptr), cout.tie(nullptr);
	return 0;
}