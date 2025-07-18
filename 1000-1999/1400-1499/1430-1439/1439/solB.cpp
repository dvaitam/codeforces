#include<bits/stdc++.h>

#define ll int

#define ld long double

#define pll pair<ll,ll>

#define pb push_back

#define mp make_pair

#define fi first

#define se second

using namespace std;

ll n , m , k ;

vector<vector<ll>> T ;

vector<ll> cnt ;

vector<bool> check , used ;

bool rrr = 0 ;

void solve()

{

    cin >> n >> m >> k ;

    map<pll , bool> M ;

    T.clear() ;

    cnt.clear() ;

    check.clear() ;

    used.clear() ;

    cnt.assign(n+5 , 0) ;

    T.assign(n+5 , vector<ll>(0 , 0)) ;

    check.assign(n+5 , 0) ;

    used.assign(n+5 , 0) ;

    for(int i = 0 ; i< m ; i++)

    {

        ll a , b ;

        cin >> a >> b ;

        T[a].push_back(b) ;

        T[b].push_back(a) ;

        M[{min(a,b),max(a,b)}] = 1 ;

        cnt[a]++ ; cnt[b]++ ;

    }

    if(((k-5)>sqrt(m)*2&&k-5>0)||k-5>n)

    {

        cout << "-1\n" ; return ;

    }

    queue<ll> q ;

    for(int i = 1 ; i<= n; i++)

    {

        if(cnt[i]<k)

        {

            q.push(i) ;

        }

    }

    while(!q.empty())

    {

        ll x = q.front() ;

        //cout << x << "\n" ;

        //cout << cnt[x] << "\n" ;

        q.pop() ;

        check[x] = 1 ;

        if(cnt[x]<k-1)

        {

            for(int i = 0 ; i< T[x].size() ; i++)

            {

                if(!check[T[x][i]])

                {

                    cnt[T[x][i]]-- ;

                    if(cnt[T[x][i]]==k-1)

                    {

                        q.push(T[x][i]) ;

                    }

                }

            }

            continue ;

        }

        if(cnt[x]==k-1)

        {

            vector<ll> save;

            save.push_back(x) ;

            used[x] = 1 ;

            bool f = 0 ;

            for(int i = 0 ; i< T[x].size() ; i++)

            {

                if(!check[T[x][i]])

                {

                    used[T[x][i]] = 1 ;

                    cnt[T[x][i]]-- ;

                    if(cnt[T[x][i]]==k-1)

                    {

                        q.push(T[x][i]) ;

                    }

                    save.push_back(T[x][i]) ;

                }

            }

            for(int i = 1 ; i< save.size() ; i++)

            {

                ll sum = 0 ;

                if(k<100)

                {

                    for(int j = 0 ; j< i ; j ++)

                    {

                        if(!M[{min(save[i],save[j]) , max(save[j],save[i])}])

                        {

                            f = 1 ;

                            break ;

                        }

                    }

                    if(f) break ;

                }

                else

                {

                    for(int j = 0 ; j < T[save[i]].size() ; j++)

                    {

                        if(used[T[save[i]][j]]) sum++ ;

                    }

                    if(sum<k-1)

                    {

                        f = 1 ;

                        break;

                    }

                }

            }

            used[x] = 0 ;

            for(int i = 0 ; i< save.size() ; i++) used[save[i]] = 0 ;

            if(f) continue ;

            cout << 2 << "\n" ;

            if(!f)

            {

                for(int i = 0 ; i< save.size() ; i++) cout << save[i] << " " ;

                cout << "\n" ;

                return ;

            }

        }

    }

    bool ans = 0 ;

    ll tong = 0;

    for(int i = 1 ; i<=n ; i++)

    {

        if(!check[i])

        {

            tong ++ ;

        }

    }

    if(tong==0)

    {

        cout << "-1\n" ;

        return ;

    }

    cout << 1 << " " << tong << "\n" ;

    for(int i = 1 ; i<=n ; i++)

    {

        if(!check[i])

        {

            cout << i << " " ;

        }

    }

    cout << "\n" ;

}

int main()

{

    ios_base::sync_with_stdio(NULL) ; cin.tie(nullptr) ; cout.tie(nullptr) ;

    //freopen("tree.inp", "r", stdin);

    //freopen("tree.out", "w", stdout);

    int t ; cin >> t ;

    while(t--) solve() ;

}