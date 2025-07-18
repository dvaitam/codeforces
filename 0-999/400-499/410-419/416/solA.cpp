#include <vector>

#include <list>

#include <map>

#include <math.h>

#include <cmath>

#include <set>

#include <queue>

#include <deque>

#include <string>

#include <stack>

#include <bitset>

#include <algorithm>

#include <functional>

#include <numeric>

#include <utility>

#include <string.h>

#include <sstream>

#include <iostream>

#include <iomanip>

#include <cstdio>

#include <cstdlib>

#include <ctime>

#include <unordered_map>



using namespace std;



#define ll  long long int

#define ld long double



int main()

{

	ios::sync_with_stdio(false);

	ios_base::sync_with_stdio(false);

	cin.tie(nullptr), cout.tie(nullptr);

	//freopen("input.txt", "r", stdin);

	//freopen("output.txt", "w", stdout);

	int n;

	ll number, l, r;

	char verify;

	string op;

	cin >> n;

	l = -2 * 1000000000;

	r = +2 * 1000000000;

	while (n--)

	{

		cin >> op>>number>>verify;

		if (op == ">")

		{

			if (verify == 'Y')

			{

				l = max(l,number + 1);

			}

			else

			{

				r = min(r,number);

			}

		}

		else if (op == "<")

		{

			if (verify == 'Y')

			{

				r = min(r,number -1);

			}

			else

			{

				l = max(l,number);

			}

		}

		else if (op == ">=")

		{

			if (verify == 'Y')

			{

				l = max(l,number);

			}

			else

			{

				r = min(number - 1,r);

			}

		}

		else if (op == "<=")

		{

			if (verify == 'Y')

			{

				r = min(r,number);

			}

			else

			{

				l = max(l,number + 1);

			}

		}

		//cout << l << " " << r << endl;

		if (l > r)

		{

			cout << "Impossible" << endl;

			return 0;

		}

	}

	cout << l;// << " " << r << endl;

	return 0;

}