#include <bits/stdc++.h>

#define int long long

using namespace std;

int a[1005][1005];

signed main()

{

	ios::sync_with_stdio(false);

	cin.tie(0);

	int T;

	cin >> T;

	while(T--)

	{

		int n,m;

		cin >> n >> m;

		int s=0;

		for(int i=1;i<=n;i++)

		{

			for(int j=1;j<=m;j++)

			{

				char c;

				cin >> c;

				a[i][j]=c-'0';

				s+=a[i][j];

			}

		}

		if(a[1][1])

		{

			cout << "-1\n";

			continue;

		}

		cout << s << "\n";

		for(int i=n;i>=2;i--)

		{

			for(int j=1;j<=m;j++)

			{

				if(a[i][j])

				{

					cout << i-1 << " " << j << " " << i << " " << j << "\n";

				}

			}

		}

		for(int i=1;i<=1;i++)

		{

			for(int j=m;j>=1;j--)

			{

				if(a[i][j])

				{

					cout << i << " " << j-1 << " " << i << " " << j << "\n";

				}

			}

		}

	}

	return 0;

}