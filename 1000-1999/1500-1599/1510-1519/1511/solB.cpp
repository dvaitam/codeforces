#define _CRT_SECURE_NO_WARNINGS

using namespace std;

#include <iostream>

#include<cstring>

#include<string>

#include<cstdlib>

#include<cmath>

#include<cstdio>

#include<algorithm>

#include<deque>

#include<vector>

#include<set>

#include<queue>

#include<map>

#include<stack>

typedef long long ll; 

int main()

{

	ios::sync_with_stdio(false);

	cin.tie(NULL);

	int t;

	cin >> t;

	while(t--)

	{

		int a, b, c;

		cin >> a >> b >> c;

		for (int i = 1; i <= a; i++)

		{

			if (i == 1)

			{

				cout << 1;

			}

			else

			{

				cout << 0;

			}

		}

		cout << " ";

		for(int i=1;i<=b-c+1;i++)

		{

			cout << 9;

		}

		for (int i = 1; i < c; i++)

		{

			cout << 0;

		}

		cout << endl;

	}

	return 0;

}