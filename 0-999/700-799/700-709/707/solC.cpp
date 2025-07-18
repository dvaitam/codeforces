#include <bits/stdc++.h>
#include<cstdio>

#include<cstring>

#include<algorithm>

using namespace std;

#define int64_t long long



int main()

{

	int n;

	scanf("%d",&n);

	int rec = 0;

	if(n == 1)

	{

		printf("-1\n");

		return 0;

	}

	if(n == 2)

	{

		printf("-1\n");

		return 0;

	}

	while(n % 2 == 0)

	{

		n /= 2;

		rec ++;

	}

	if(n == 1)

	{

		long long x = 3,y = 5;

		for(int i = 0;i < rec-2;i ++)

		{

			x *= 2;y *= 2;

		}

		printf("%lld %lld\n",x,y);

		return 0;

	}

	long long x = n / 2;

	long long y = x + 1;

	long long ans1 = 2*x*y;

	long long ans2 = x*x + y*y;

	for(int i = 0;i < rec;i ++)

	{

		ans1 *= 2;ans2 *= 2;

	}

	printf("%lld %lld\n",ans1,ans2);

}