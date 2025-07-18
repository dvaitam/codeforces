//#pragma linker("/STACK:134217728")



#include <cstdio>  // need for the freopen func.

#include <cstddef> // for ::std::size_T



#include <iostream>  // for ::std::cout

#include <iomanip>

#include <vector>    // for ::std::vector

#include <string>    // for ::std::string

#include <algorithm> // for algorithms

#include <queue>     // for  ::std::queue

#include <set>       // for  ::std::set

#include <utility>

#include <numeric>

#include <cmath>



int gcd(int x, int y)

{

	while (y > 0)

	{

		int m = x % y;

		x = y;

		y = m;

	}

	return x;

}

 

int main(int argc, char* argv[])

{



#ifndef ONLINE_JUDGE

	freopen("input.txt", "r", stdin);

	freopen("output.txt", "w", stdout);

#endif

	

	::std::ios::sync_with_stdio(false);

	using  ::std::cout;

	using  ::std::cin;

	using  ::std::vector;

	using  ::std::string;

	using  ::std::pair;



	::std::size_t n;

	double a[32][32];

	vector<double> d;

	vector<double>ans;

	

	d.resize(1 << 18);

	int s[32] = { 0 };

	int ss = 0;





	cin >> n;

 

	ans.resize(n);



	for (size_t i = 0; i < n; ++i)

		for (size_t j = 0; j < n; ++j)

			cin >> a[i][j];





	d[(1 << n) - 1] = 1;

	for (size_t mask = (1 << n) - 1; mask > 0; --mask)

	{

		ss = 0;

		for (size_t i = 0; i < n;++i)

			if (mask & (1 << i))

				s[ss++] = i;



		if (ss == 1)

		{

			ans[s[0]] = d[mask];

		}

		else

		{

			double prob = 2.0 / (ss * (ss - 1) );

			for (size_t i = 0; i < ss - 1; ++i)

				for (size_t j = i + 1; j < ss; ++j)

				{

					d[mask ^ (1 << s[i])] += d[mask] * prob * a[s[j]][s[i]];

					d[mask ^ (1 << s[j])] += d[mask] * prob * a[s[i]][s[j]];

				}

		}

	}

	for (size_t i = 0; i < n;++i)

		cout << std::setprecision(12) << std::fixed << ans[i] << ' ';

	cout << '\n';



	return 0;

}