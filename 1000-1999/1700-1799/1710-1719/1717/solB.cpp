#include <bits/stdc++.h>

using namespace std;



typedef long long ll;



constexpr int N = 3e5 + 10;

constexpr int inf = 0x3f3f3f3f;



void solve()

{	

    int n, k, r, c;

    cin >> n >> k >> r >> c;



    r--, c--;

    int dif = r % k - c % k;



    if (dif < 0)

        dif += k;

    

    for (int i = 0; i < n; i++)

    {

        for (int j = 0; j < n; j++)

        {

            int t = i % k - j % k;

            if (t < 0)

                t += k;

            if (t == dif)

                cout << "X";

            else

                cout << ".";

        }



        cout << "\n";

    }

}



int main()

{

	ios::sync_with_stdio(false);

	cin.tie(nullptr);

	cout.tie(nullptr);





	int T;

	cin >> T;



	while (T--)

	    solve();



	return 0;

}