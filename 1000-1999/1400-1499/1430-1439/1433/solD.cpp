#include <bits/stdc++.h>

using namespace std;

int main()
{
	ios_base::sync_with_stdio(false);
	cout << setprecision(10) << fixed;

	int t;
	cin >> t;
	while (t--)
	{
		int n;
		cin >> n;
		vector<int> nums(n);
		int mx = 0, mn = INT_MAX;
		int idmx = -1, idmn = -1;
		for (int i = 0; i < n; ++i)
		{
			cin >> nums[i];
			if (nums[i] > mx)
			{
				mx = nums[i];
				idmx = i;
			}
			if (nums[i] < mn)
			{
				mn = nums[i];
				idmn = i;
			}
		}
		if (mx == mn)
			cout << "NO" << '\n';
		else
		{
			cout << "YES\n";
			bool chk = false;
			for (int i = 0; i < n; ++i)
			{
				if (!chk && idmx == i)
				{
					chk = true;
					continue;
				}
				if (chk && nums[i] == mx)
					cout << i + 1 << ' ' << idmn + 1 << '\n';
				else
					cout << idmx + 1 << ' ' << i + 1 << '\n';
			}
		}
	}
}