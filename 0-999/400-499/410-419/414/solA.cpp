//By Don4ick

#include <bits/stdc++.h>



using namespace std;

int n, k;



int main()

{

	scanf("%d%d", &n, &k);

	

	if (k == 0 && n == 1)

	{

		cout << 1 << endl;

		return 0;

	}

	

	if (k < (n / 2) || n == 1)

	{

		cout << -1 << endl;

		return 0;

	}



	printf("%d %d ", (k - ((n / 2) - 1)), (k - ((n / 2) - 1)) * 2);

	int x = (k - ((n / 2) - 1)) * 2;

	for (int i = 4; i <= (int)n; i += 2)

	{

		printf("%d %d ", x + 1, x + 2);

		x += 2;

	}

	

	if (n % 2)

		printf("%d", x + 1);

	

	

	return 0;

}