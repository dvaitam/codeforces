#include <iostream>

using namespace std;

void Solve();

int main()

{

	Solve();

	//system("pause");

	return 0;

}

void Solve()

{

	long long a, b, n, count = 0,cpy;

	cin >> a >> b >> n;

	bool bol = false;

	for (int i = 0; i < 10; i++)

	{

		if (((a * 10) + i) % b == 0)

		{

			bol = true;

			cpy = i;

			break;

			

		}

	}

	if (bol){

		cout << a;

		cout << cpy;

		for (int i = 0; i < n - 1; i++)

			cout << 0;

	}

	else

		cout << -1 << endl;

}