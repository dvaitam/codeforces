#include <stdio.h>



int dat[26];

int gcd(int a, int b)

{

	return b ? gcd(b, a % b) : a;

}



int main()

{

	int n;

	scanf("%d",&n);

	int odd = 0;

	for (int i = 0; i < n; i++)

	{

		scanf("%d",&dat[i]);

		odd += dat[i] & 1;

	}



	if (n == 1)

	{

		printf("%d\n",dat[0]);

		for (int i = 0; i < dat[0]; i++) putchar('a');

		return 0;

	}



	if (odd >= 2)

	{

		printf("0\n");

		for (int i = 0; i < n; i++)

			for (int j = 0; j < dat[i]; j++)

				putchar(i+'a');

	}

	else

	{

		int g = dat[0];

		for (int i = 1; i < n; i++) g = gcd(g, dat[i]);

		printf("%d\n", g);



		if (odd==0)

		for (int tt = 0; tt < g / 2; tt++)

		{

			for (int i = 0; i < n; i++)

				for (int j = 0; j < dat[i] / g; j++)

					putchar(i+'a');

			for (int i = n-1; i >= 0; i--)

				for (int j = 0; j < dat[i] / g; j++)

					putchar(i+'a');

		}

		else

		for (int tt = 0; tt < g; tt++)

		{

			for (int i = 0; i < n; i++)

			{

				if (dat[i] & 1) continue;

				for (int j = 0; j < dat[i] / g / 2; j++)

					putchar(i + 'a');

			}

			for (int i = 0; i < n; i++)

			{

				if (dat[i] & 1)

				{

					for (int j = 0; j < dat[i] / g; j++)

						putchar(i + 'a');

					break;

				}

			}

			for (int i = n - 1; i >= 0; i--)

			{

				if (dat[i] & 1) continue;

				for (int j = 0; j < dat[i] / g / 2; j++)

					putchar(i + 'a');

			}

		}

	}

}





// abab

// abbaabba

// baabbaab

// abbaabba

// baabbaab

//

// abbaabbaabba

// baabbaabbaab

// abbaabbaabba

// baabbaabbaab

// abbaabbaabba

// baabbaabbaab

//

// 그러니까 a1, a2, a3, a4, ..., ak가 있다

// g=gcd(a1, a2, a3, ..., ak)를 구하자

// 팰린드롬을 a1/g, a2/g, ..., ak/g개를 사용해서 만든다음

// 그걸 g개를 붙인다. 그러면 완성!