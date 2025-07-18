# include <iostream>

# include <algorithm>

# include <cstring>

# include <string>

# include <vector>

# include <cmath>

# include <map>

# include <unordered_map>

# include <set>

# include <queue>

# include <deque>

# include <bitset>

# include <stack>

# define mod int(1e9 + 7)

# define NoDejareDeQererte ios_base::sync_with_stdio(0), cin.tie(0), cout.tie(0);

using namespace std;

vector<int>SPF(1e3 + 5);

void sieve()

{

    for (int i = 0; i <= 1e3; i++)

    {

        SPF[i] = i;

    }

    for (int i = 2; i * i <= 1e3; i++)

    {

        for (int j = i * 2; j <= 1e3; j += i)

        {

            SPF[j] = min(SPF[j], i);

        }

    }

}

int main()

{

    NoDejareDeQererte;

    sieve();

    int t;

    cin >> t;

    while (t--)

    {

        int n;

        cin >> n;

        vector<int>v(n);

        for (int i = 0; i < n; i++)

        {

            cin >> v[i];

        }

        map<int, int>mp;

        int p = 1;

        for (auto it : v)

        {

            if (!mp[SPF[it]])

            {

                mp[SPF[it]] = p++;

            }

        }

        cout << mp.size() << endl;

        for (auto it : v)

        {

            cout << mp[SPF[it]] << ' ';

        }

        cout << endl;

    }

    return 0;

}