#include<bits/stdc++.h>

using namespace std;

#define all(v)		((v).begin()), ((v).end())

#define sz(v)		((int)((v).size()))

#define clr(v, d)	memset(v, d, sizeof(v))

#define pb					push_back

#define MP					make_pair

#define ff  first

#define ss second

#define mod(a,b)  (a%b+b)%b



typedef long long         ll;

const ll mod = 1e9 + 9;





int vis[301];

int main ()

{

    int n, m, k, ss = 0;

    cin >> n >> m >> k;

    for(int i = 0; i < k; i++)

    {

        int x;

        cin >> x;

        vis[x] = 1;

        if(i == 0)

        {

            vis[x] = 2;

            ss = x;

        }

    }

    vector<int> a, b;

    for(int i = 1; i <= n; i++)

    {

        if(vis[i] != 2) a.push_back(i);

        if(vis[i] == 0) b.push_back(i);

    }

    int tot = (n - k) + ( (n - 1) * (n - 2) ) / 2;

    if(m > tot || k == n)

        cout << -1;

    else

    {

        for(int i = 0; i < n - 1 && m > 1; i++)

            for(int j = i + 1; j < n - 1 && m > 1; j++)

            {

                printf("%d %d\n",a[i] , a[j]);

                m--;

            }

        for(int i = 0; i < m; i++)

            printf("%d %d\n",ss , b[i]);

    }

    cout << endl;

}