#include <bits/stdc++.h>

#define x first

#define y second



using namespace std;



typedef pair<int, int> PII;

typedef long long LL;

const int N = 100010;

int n, m1, m2;

int fa1[N], fa2[N];



int find1(int x)

{

    return fa1[x]==x ? x : fa1[x]=find1(fa1[x]);

}



int find2(int x)

{

    return fa2[x]==x ? x : fa2[x]=find2(fa2[x]);

}



void solve()

{

    cin >> n >> m1 >> m2;

    for(int i = 1; i <= n; i++)

        fa1[i] = fa2[i] = i;

    while(m1--)

    {

        int u, v;

        cin >> u >> v;

        int fu = find1(u), fv = find1(v);

        if(fu != fv)

            fa1[fu] = fv;

    }

    while(m2--)

    {

        int u, v;

        cin >> u >> v;

        int fu = find2(u), fv = find2(v);

        if(fu != fv)

            fa2[fu] = fv;

    }

    vector<PII> v;

    for(int i = 2; i <= n; i++)

    {

        int fi1 = find1(i), fj1 = find1(1);

        int fi2 = find2(i), fj2 = find2(1);

        if(fi1 != fj1 && fi2 != fj2)

        {

            v.push_back({i, 1});

            fa1[fi1] = fj1;

            fa2[fi2] = fj2;

        }

    }

    vector<int> s1, s2;

    for(int i = 2; i <= n; i++)

    {

        if(find1(i) == i && find1(i) != find1(1))

            s1.push_back(i);

        if(find2(i) == i && find2(i) != find2(1))

            s2.push_back(i);

    }

    for(int i = 0, j = 0; i < s1.size()&&j < s2.size(); i++, j++)

        v.push_back({s1[i], s2[j]});

    cout << v.size() << '\n';

    for(auto p : v)

        cout << p.x << ' ' << p.y << '\n';

}



int main()

{

    ios::sync_with_stdio(false);

    cin.tie(0);

    cout.tie(0);

    int T = 1;

    //cin >> T;

    while(T--)

        solve();

}