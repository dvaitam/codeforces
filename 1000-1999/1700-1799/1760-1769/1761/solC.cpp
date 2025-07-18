//#pragma GCC optimize(2)

#include <bits/stdc++.h>

using namespace std;

#define endl '\n'

typedef pair<int,int> PII;

typedef long long LL;

const int INF=0x3f3f3f3f,mod=1e9 + 7;

const double eps=1e-6;

const double pai=acos(-1.0);

//3.14159

int T;

const int N = 101;

vector<int> e[N];

int in[N];

int top[N],cnt = 0;

void solve()

{

    int n;

    cin >> n;

    vector<int> res[N];

    for(int i = 1; i <= n; i++) {

        string s;

        cin >> s;

        res[i].push_back(i);

        for (int j = 0; j < s.size(); j++) {

            if (s[j] == '1') {

                res[j + 1].push_back(i);

            }

        }

    }

    for(int i = 1; i <= n; i++)

    {

        cout << res[i].size() << " ";

        for(auto t : res[i])

        {

            cout << t << " " ;

        }

        cout << endl;

    }

}

int main()

{

    ios::sync_with_stdio(false);

    cin.tie(0);

    cout.tie(0);

    cin >> T;

    while(T--)

        solve();

    return 0;

}