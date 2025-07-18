#include <bits/stdc++.h>

using namespace std;

#define lli long long int



struct A

{

	int id , val;

	A( int _i , int _v ) { id = _i , val = _v; }

	bool friend operator < ( A x , A y )

	{

		return x.val < y.val;

	}

};



int ans[212345];

int taim[212345];	// time

bool usedN[212345];

bool usedM[212345];



int main()

{

	int n , m;

	scanf( "%d %d" , &n , &m );

	

	vector< pair<int,int> > s;

	for ( int i=1 ; i<=n ; i++ ) 

	{

		int x; scanf( "%d" , &x );

		s.push_back( make_pair( x , i ) );

	}

	sort( s.begin() , s.end() );

	

	vector<A> v;

	for ( int i=1 ; i<=m ; i++ )

	{

		int x; scanf( "%d" , &x );

		v.push_back( A( i , x ) );

	}

	sort( v.begin() , v.end() );

	

	// ----------------------------------------------------------------

	

	int cnt = 0 , u = 0;

	for ( int t=0 ; t<32 ; t++ )

	{

		int j = 0;

		for ( int i=0 ; i<m ; i++ )

		{

			if ( usedM[i] ) continue;

			int val = v[i].val , id = v[i].id;

			

			while ( (j < n) && ( (usedN[j]) || (val>s[j].first) ) ) j++;

			

			if ( j < n )

			{

				int val2 = s[j].first , id2 = s[j].second;

				if ( val2 == val )

				{

					ans[ id2 ] = id;

					cnt++;

					u += t;

					taim[ id ] = t;

					usedN[ j ] = true;

					usedM[ i ] = true;

				}

			}

			

			v[i].val = (v[i].val/2) + (v[i].val%2);

		}

	}

	

	printf( "%d %d\n" , cnt , u );

	for ( int i=1 ; i<=m ; i++ ) printf( "%d " , taim[i] ); printf( "\n" );

	for ( int i=1 ; i<=n ; i++ ) printf( "%d " , ans[i] ); printf( "\n" );

	

	return 0;

}