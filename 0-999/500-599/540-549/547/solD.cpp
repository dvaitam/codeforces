/*

                ____________    ______________       __

               / _________  /\ /_____   _____/\     / /\

              / /\       / /  \\    /  /\    \ \   / /  \

             / /  \_____/ /   / \__/  /  \____\/  / /   /

            / /   /    / /   /    /  /   /       / /   /

           / /   /    / /   /    /  /   /       / /   /

          / /   /    / /   /    /  /   /       / /   /

         / /___/____/ /   /    /  /   /       / /___/________

        /____________/   /    /__/   /       /______________/\

        \            \  /     \  \  /        \              \ \

         \____________\/       \__\/          \______________\/

           ___       ___               ___    __________

          /  /\     /  /\             /  /\  /_______  /\

         /  /__\___/  /  \           /  /  \ \      /  /  \

        /____    ____/   /          /  /   /  \____/  /   /

        \   /   /\   \  /          /  /   /       /  /   /

         \_/   /  \___\/ ___      /  /   /       /  /   /

          /   /   /     /  /\    /  /   /       /  /   /

         /   /   /     /  /__\__/  /   /       /  /___/____

        /___/   /     /___________/   /       /___________/\

        \   \  /      \           \  /        \           \ \

         \___\/        \___________\/          \___________\/

       

          A FAN OF FIZZYDAVID

           

*/



#include<bits/stdc++.h>



#define HEAP priority_queue

#define rep(i,n) for(int i=0;i<(n);i++)

#define per(i,n) for(int i=(n)-1;i>=0;i--)

#define forn(i,l,r) for(int i=(l);i<=(r);i++)

#define nrof(i,r,l) for(int i=(r);i>=(l);i--)

#define snuke(c,itr) for(__typeof((c).begin()) itr=(c).begin();itr!=(c).end();itr++)

#define mp make_pair

#define pb push_back

#define X first

#define Y second

#define eps 1e-20

#define pi 3.1415926535897932384626433832795

#define orz int

#define yjz main

#define fizzydavid return

#define ak 0

#define la ;

#define debug puts("OK");

#define SZ(x) (int)x.size()

#define ALL(x) x.begin(),x.end()

#define foreach(i,u) for(int i=head[(u)];i;i=e[i].nxt)



using namespace std;



typedef long long LL;

typedef double flt;

typedef vector<int> vi;

typedef pair<int,int> pii;



const int iinf=1e9+7;

const LL linf=1ll<<60;

const flt dinf=1e10;



inline LL read()

{

	bool f=0; LL x=0; char c=getchar();

	while((c<'0' || c>'9') && c!='-') c=getchar();

	if(c=='-') { f=1; c=getchar(); }

	while(c>='0' && c<='9') { x=x*10+c-'0'; c=getchar(); }

	if(f) x=-x; return x;

}



inline void scf(int &x)

{

	bool f=0; x=0; char c=getchar();

	while((c<'0' || c>'9') && c!='-') c=getchar();

	if(c=='-') { f=1; c=getchar(); }

	while(c>='0' && c<='9') { x=x*10+c-'0'; c=getchar(); }

	if(f) x=-x; return;

}



void scf(int &x,int &y) { scf(x); return scf(y); }



void scf(int &x,int &y,int &z) { scf(x); scf(y); return scf(z); }



//---------------------------head----------------------------



const int N = 2e5 + 100 ;



int n ;

int lc [ N ] , lr [ N ] ;

vi g [ N ] ;

int col [ N ] ;



inline void add_edge ( int u , int v ) { g [ u ] . pb ( v ) , g [ v ] . pb ( u ) ; return ; }



inline void dfs ( int u , int c = 0 )

{

	col [ u ] = c ;

	snuke ( g [ u ] , v ) if ( col [ *v ] < 0 ) dfs ( *v , c ^ 1 ) ;

	return ;

}



orz yjz()

{

	scf ( n ) ;

	memset ( col , -1 , sizeof ( col ) ) ;

	forn ( i , 1 , n )

	{

		int x , y ; scf ( x , y ) ;

		if ( lc [ x ] ) add_edge ( i , lc [ x ] ) , lc [ x ] = 0 ; else lc [ x ] = i ;

		if ( lr [ y ] ) add_edge ( i , lr [ y ] ) , lr [ y ] = 0 ; else lr [ y ] = i ;

	}

	forn ( i , 1 , n ) if ( col [ i ] < 0 ) dfs ( i ) ;

	forn ( i , 1 , n ) putchar ( col [ i ] ? 'r' : 'b' ) ;

	fizzydavid ak la

}