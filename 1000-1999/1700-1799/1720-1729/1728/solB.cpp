#include<iostream>

#include<cstdio>

using namespace std;

int main()

{

	int t;

	cin >> t;

	while (t--)

	{

		int n;

		scanf("%d", &n);

		if (n & 1)//判断n为奇数

		{

			for (int i = n - 2; i >= 5; i-=2)

			{

				printf("%d %d ", i, i - 1);//输出一个数接这个数-1

			}

			for (int i = 1; i <= 3; i++)

			{

				printf("%d ", i);//输出123

			}

		}

		else

		{//n为偶数

			for (int i = n - 2; i >= 2; i-=2)

			{

				printf("%d %d ", i, i - 1);//输出一个数接这个数-1

			}

		}

		printf("%d %d", n - 1, n);//输出n-1，n

		printf("\n");

	}

	return 0;

}