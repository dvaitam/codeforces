#include <bits/stdc++.h>

using namespace std;



const int N = 1010;



int n;

int a[N];

vector<pair<int,int> > path;



bool dfs( int remain ) {

	bool ok = true;

	for( int i = 1; i <= n; i++ )

		if( a[i] != i ) {

			ok = false;

			break;

		}

	if( ok ) return true;

	if( remain==0 ) return false;

	vector<pair<int,int> > vc;

	vc.push_back( make_pair(1,n) );

	int pos[N];

	for( int i = 1; i <= n; i++ )

		pos[a[i]] = i;

	for( int i = 1, j; i <= n; i = j + 1 ) {

		for( j = i; j+1 <= n && a[j+1]-a[j]==a[i+1]-a[i] && abs(a[i+1]-a[i])==1; j++ );

		vc.push_back( make_pair( i, j ) );

	}

	for( int i = 1; i <= n; i++ ) {

		if( a[i] != 1 ) {

			int j = pos[a[i]-1];

			if( j+1 < i )

				vc.push_back( make_pair( j+1, i ) );

			if( i+1 < j )

				vc.push_back( make_pair( i, j-1 ) );

		}

		if( a[i] != n ) {

			int j = pos[a[i]+1];

			if( j+1 < i )

				vc.push_back( make_pair( j+1, i ) );

			if( i+1 < j )

				vc.push_back( make_pair( i, j-1 ) );

		}

	}

	/*

	for( int i = 1; i <= n; i++ )

		fprintf( stderr, "%d ", a[i] );

	fprintf( stderr, "\n" );

	for( int t = 0; t < vc.size(); t++ )

		fprintf( stderr, "%d %d\n", vc[t].first, vc[t].second );

	*/

	for( int t = 0; t < vc.size(); t++ ) {

		int p = vc[t].first, q = vc[t].second;

		path.push_back( vc[t] );

		reverse( a + p, a + q + 1 );

		if( dfs( remain-1 ) ) return true;

		reverse( a + p, a + q + 1 );

		path.pop_back();

	}

	return false;

}

int main() {

	scanf( "%d", &n );

	for( int i = 1; i <= n; i++ )

		scanf( "%d", a+i );

	dfs( 3 );

	printf( "%d\n", (int)path.size() );

	for( int t = (int)path.size()-1; t >= 0; t-- ) {

		printf( "%d %d\n", path[t].first, path[t].second );

	}

}