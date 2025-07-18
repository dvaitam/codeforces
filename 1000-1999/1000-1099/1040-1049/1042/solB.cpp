#include<bits/stdc++.h>
#define ms(a, b) memset(a, b, sizeof(a));
#define lowbit(x) (x&(-x))
#define rep(i,a,b) for(int i = a; i <= b; i++)
#define _rep(i,a,b) for(int i = a; i >= b; i--)
using namespace std;
typedef long long ll;
typedef unsigned long long ull;
typedef pair<int, int> pii;
const int maxn = 2e5 + 10;
const int maxm = 1e6 + 10;
const int maxb = 26;
const int sigma_size = 62;
const double eps = 1e-8;
const ll mod =  1e9 + 7;
const double pi = acos(-1.0);
const int Ha = 2333;
const ll Linf = 3e12 + 10;
const int inf = 3e5 + 7;
const double Dinf = 1000000007.0;

map<string,int> mp;
char s[5];
string model[] = {"A","B","C","AB","AC","BC","ABC"};
string sl[3][3] = {{"A","AB","AC"},{"B","AB","BC"},{"C","AC","BC"}};

int main()
{
#ifdef LOCAL
    freopen("input.in","r",stdin);
//    freopen("out.in","w",stdout);
#endif
    int n;
    while(~scanf("%d",&n)){
        for(int i = 0; i < 7; i++) mp[model[i]] = inf;
        int c;
        for(int i = 1; i <= n; i++){
            scanf("%d%s",&c,s);
            sort(s,s+strlen(s));
            mp[(string)s] = min(mp[(string)s],c);
        }
        int ans = mp["ABC"];
        for(int i = 0; i < 3; i++){
            for(int j = 0; j < 3; j++){
                for(int k = 0; k < 3; k++){
                    ans = min(ans,mp[sl[0][i]]+mp[sl[1][j]]+mp[sl[2][k]]);
                }
            }
        }
        ans = min(ans,mp["A"]+mp["BC"]);

        ans = min(ans,mp["AB"]+mp["BC"]);
        ans = min(ans,mp["AB"]+mp["AC"]);
        ans = min(ans,mp["AB"]+mp["C"]);

        ans = min(ans,mp["AC"]+mp["B"]);
        ans = min(ans,mp["AC"]+mp["AB"]);
        ans = min(ans,mp["AC"]+mp["BC"]);

        ans = min(ans,mp["B"]+mp["AC"]);

        ans = min(ans,mp["BC"]+mp["A"]);
        ans = min(ans,mp["BC"]+mp["AB"]);
        ans = min(ans,mp["BC"]+mp["AC"]);

        if(ans == inf) ans = -1;
        printf("%d\n",ans);
    }
    return 0;
}