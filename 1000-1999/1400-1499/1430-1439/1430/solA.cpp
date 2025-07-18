#include <bits/stdc++.h>



using namespace std;



void solve()

{

    int n;

    cin >> n;

    for (int x = 0; x < 334; x++)

    {

        int left = n - 3 * x;

        int y = 0, z = 0;

        while (left > 0 && left % 5 != 0)

        {

            ++z;

            left -= 7;

        }

        if (left >= 0 && left % 5 == 0)

        {

            y = left / 5;

        }

        if (x * 3 + y * 5 + z * 7 == n)

        {

            cout << x << ' ' << y << ' ' << z << '\n';

            return;

        }

    }

    cout << "-1\n";

}



int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(0);

    int t = 1;

    cin >> t;

    while (t--)

    {

        solve();

    }

}