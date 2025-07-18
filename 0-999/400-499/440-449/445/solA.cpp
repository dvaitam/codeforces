#include <bits/stdc++.h>

#define ll long long

#define mp make_pair

#define pb push_back

#define print(v) for(ll i=0;i<v.size();i++)cout<<v[i]<<" "

using namespace std;

char mat[125][125];

char f(pair<ll,ll> x)

{

    if(mat[x.first][x.second]=='B')

        return 'W';

    else

        return 'B';

}

int main()

{

    //freopen("input.txt","r",stdin);

    // freopen("output.txt","w",stdout);

    ll n,m;



    queue< pair<ll,ll> >myqueue;

    cin>>n>>m;

    for(ll i=1; i<=n; i++)

    {

        for(ll j=1; j<=m; j++)

            cin>>mat[i][j];

    }

    // vector< vector<bool> > visited(125 , vector<bool>(125,false));

    for(ll i=1; i<=n; i++)

    {

        for(ll j=1; j<=m; j++)

        {

            if(mat[i][j]=='.')//&& mat[i][j] =='.')

            {

                mat[i][j]='B';

                myqueue.push(mp(i,j));

               // visited[i][j]=false;

                while(!myqueue.empty())

                {

                    pair<ll,ll> top = myqueue.front();

                    myqueue.pop();

                    ll a=top.first;

                    ll b=top.second;

                    if(mat[a-1][b]=='.' )// && visited[i-1][j])

                        {

                            mat[a-1][b]=f(top);

                            myqueue.push(mp(a-1,b));

                        }

                    if(mat[a+1][b]=='.')

                       {

                           mat[a+1][b]=f(top);

                           myqueue.push(mp(a+1,b));

                       }

                    if(mat[a][b-1]=='.')

                        {

                            mat[a][b-1]=f(top);

                            myqueue.push(mp(a,b-1));

                        }

                    if(mat[a][b+1]=='.')

                        {

                            mat[a][b+1]=f(top);

                            myqueue.push(mp(a,b+1));

                        }



                }

            }

        }

    }

    for(ll i=1;i<=n;i++)

    {

        for(ll j=1;j<=m;j++)

            cout<<mat[i][j];

        cout<<endl;

    }





    return 0;

}