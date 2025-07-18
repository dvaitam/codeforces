using namespace std;

#include <iostream>
#include <cmath>
#include <algorithm>
#include <string>

typedef long long int ll;

ll arr[105];

int main()
{
	ll n;
	cin >> n;

	for(ll i=1; i<=n; i++)
	{
		cin >> arr[i];
	}

	ll min = 10000000000000;
	for(ll j=1; j<=n; j++)
	{
		ll sum = 0;
		// cout << j << endl;
		for(ll k=1; k<=n; k++)
		{
			// cout << k << " " << arr[k] << " " << abs(j-k) << " " << 2*j << endl;
			sum+=(arr[k]*(abs(k-j) + k-1 + j-1 + j-1 + k-1 + abs(k-j)));
			// cout << k << " " << arr[k] << " " << abs(j-k) << " " << j << " " << j << " " << k << " " << abs(j-k) << endl;
			// cout << sum << endl;	
		}

		if(sum<min)
		{
			min = sum;
		}
	}

	cout << min << endl;

	return 0;
}