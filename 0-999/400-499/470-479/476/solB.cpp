using namespace std;

#include <bits/stdc++.h>

#define FOR(i,a,b) for(int i=a,_b=(b);i<=_b;i++)

#define op freopen("in.inp","r",stdin);freopen("in.out","w",stdout);

#define m(s) memset(s,0,sizeof s);

#define fd(i,a,b) for(int i=a,_b=(b);i>=_b;i--)

#define st first

#define nd second

#define PB push_back

#define TR(c, it) for(__typeof((c).begin()) it=(c).begin(); it!=(c).end(); it++)

typedef long long ll;

typedef pair<int,int> II;

const int MAXN=1e5+10;

const int INF = 1e9 + 9;

const int MOD = 1e9 + 7;



int ans = 0, p1 = 0;

string s, a;



void Btr(int x)

{

    if(x == s.size())

    {

        //cout << a << endl;

        int p = 0;

        FOR(i, 0, x - 1)

        if(a[i] == '+') p++;

        else p--;

        //cout << p << endl;

        if(p == p1) ans++;

        //cout <<"End case.\n";

        return;

    }



    if(s[x] != '?')

    {

        a = a + s[x];

        Btr(x + 1);

        a.erase(a.size() - 1, 1);

    }

    else

    {

        a = a + '+';

        Btr(x + 1);

        a.erase(a.size() - 1, 1);

        a = a + '-';

        Btr(x + 1);

        a.erase(a.size() - 1, 1);



    }

}

int main()

{

    cin >> s;



    FOR(i, 0, s.size() -1)

    if(s[i] == '+') p1++;

    else p1--;



    //cout << p1 << endl;

    cin >> s;

    int z = 0;

    FOR(i, 0, s.size() - 1)

    if(s[i] == '?') z++;





    Btr(0);

    //cout << ans << " " << z << endl;

    cout << setprecision(9) << fixed;

    cout << double(ans) / double(1 << z);

}