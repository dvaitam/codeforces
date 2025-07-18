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

int book[500];

int main()

{

	ios::sync_with_stdio(false);

	cin.tie(NULL);

	int t;

	cin >> t;

	while(t--)

	{

		string a;

		cin >> a;

		string nex = "";

		memset(book, 0, sizeof(book));

		for (int i = 0; i < a.size(); i++)

		{

			book[a[i]]++;

			if (book[a[i]] == 2)

			{

				nex += a[i];

			}

		}

		cout << nex + nex;

		for (char i = 'a'; i <= 'z'; i++)

		{

			if (book[i] == 1)

			{

				cout << (char)i;

			}

		}

		cout << endl;

	}

	return 0;

}