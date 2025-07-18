#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

#define trace(x)    cerr << #x << ": " << x << endl;
#define bitcount    __builtin_popcountll
#define MOD 1000000007
#define pb push_back
#define pi pair<int,int>
#define pii pair<pi,int>
#define mp make_pair



vector <int> a[205];

int b[205];

vector <int> c, a1, a2;

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    cout.tie(NULL);
    // freopen("A-small-attempt0.in", "r", stdin);
    // freopen("a1.txt", "w", stdout);
    int t=1, i, j, x, y, n, m, k, l;
    //cin>>t;
    while(t--)
    {

        cin>>n>>m;

        for(i=1; i<=m; i++)

        {

            cin>>x;

            for(j=0; j<x; j++)

            {

                cin>>y;

                a[i].push_back(y);

                b[y]=i;

            }

        }

        for(i=1; i<=n; i++)

            if(!b[i])

                c.push_back(i);

        k=1;

        for(i=1; i<=m; i++)

        {

            for(j=0; j<a[i].size(); j++)

            {

                if(a[i][j]==k)

                {

                    k++;

                    continue;

                }

                if(!b[k])

                {

                    b[k]=i;

                    b[a[i][j]]=0;

                    c.push_back(a[i][j]);

                    //cout<<a[i][j]<<" "<<k<<endl;

                    a1.pb(a[i][j]);

                    a2.pb(k);

                    a[i][j]=k;

                }

                else

                {

                    x=c.back();

                    c.pop_back();

                    c.push_back(a[i][j]);

                    for(l=0; l<a[b[k]].size(); l++)

                    {

                        if(a[b[k]][l]==k)

                        {

                            a[b[k]][l]=x;

                            break;

                        }

                    }

                    b[x]=b[k];

                    b[k]=i;

                    b[a[i][j]]=0;

                    //cout<<k<<" "<<x<<endl;

                    //cout<<a[i][j]<<" "<<k<<endl;

                    a1.pb(k);

                    a2.pb(x);

                    a1.pb(a[i][j]);

                    a2.pb(k);

                    a[i][j]=k;

                }

                k++;

            }

        }

        cout<<a1.size()<<endl;

        for(i=0; i<a1.size(); i++)

            cout<<a1[i]<<" "<<a2[i]<<endl;



    }
    return 0;
}