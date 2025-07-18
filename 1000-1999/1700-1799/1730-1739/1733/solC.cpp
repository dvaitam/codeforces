#include <bits/stdc++.h>

using namespace std;



typedef long long ll;



constexpr int N = 1e5 + 10;

constexpr int inf = 0x3f3f3f3f;



int a[N];





void solve()

{	

    int n;

    cin >> n;

    

    for (int i = 1; i <= n; i++)

        cin >> a[i];

    

    vector<pair<int, int>> ans;



    if (a[1] != a[n])

    {

        ans.push_back({1, n});

        if ((a[1] + a[n]) & 1)

            a[n] = a[1];

        else

            a[1] = a[n];

    }



    for (int i = 2; i <= n - 1; i++)

    {

        if (a[i] == a[1])

            continue;

        

        if ((a[i] + a[1]) & 1)

            ans.push_back({1, i});

        else

            ans.push_back({i, n});

    }



    cout << ans.size() << "\n";



    for (auto [x, y] : ans)

        cout << x << " " << y << "\n";



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