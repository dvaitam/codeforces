#include <iostream>

#include <bits/stdc++.h>

using namespace std;

int main()

{

    ios::sync_with_stdio(false); cin.tie(0);

    int n, s, t;

    cin >> n;

    vector<int> a(n);

    int mx = 0;

    for( int i=0 ; i<n ; i++ )

    {

        cin >> a[i];

        mx = max( mx , a[i] );

    }

    cin >> s >> t;

    s--; t--;



    vector<int> minp(mx+1), primes;

    for( int i=2 ; i<=mx ; i++ )

    {

        if( !minp[i] )

        {

            minp[i] = i;

            primes.push_back(i);

        }

        for( auto k : primes )

        {

            if( k * i > mx ) break;

            minp[ k * i ] = k;

            if( k == minp[i] ) break;

        }

    }



    vector< vector<int> > edg(mx+1);

    for( int i=0 ; i<n ; i++ )

    {

        for( int nw = a[i] ; nw > 1 ; nw /= minp[nw] )

            edg[ minp[nw] ].push_back( i );

    }



    vector<int> dist(n+mx+1,-1), nxt(n+mx+1);

    queue<tuple<int,int,int>> q;

    q.emplace(t,0,-1);



    while(!q.empty())

    {

        auto [ w , d , t ] = q.front();

        q.pop();

        if( dist[w] != -1 ) continue;



        dist[w] = d;

        nxt[w] = t;

        if( w < n )

        {

            for( int i=a[w] ; i>1 ; i/=minp[i] ) 

                if( dist[ n+minp[i] ] == -1 ) q.emplace( n+minp[i] , d+1 , w );

        }

        else

        {

            for( auto k : edg[w-n] ) 

                if( dist[k] == -1 )  q.emplace( k , d+1 , w );

        }

    }



    if( dist[s] == -1 ) cout << "-1" << endl;

    else

    {

        cout << dist[s] / 2 + 1 << endl << s + 1;

        for( int i=nxt[s] ; i!=-1 ; i=nxt[i] )

            if( i < n ) cout << " " << i+1;

        cout << endl;

    }



    return 0;

}