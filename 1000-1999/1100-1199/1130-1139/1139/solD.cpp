#include <bits/stdc++.h>
#define lld I64d
using namespace std ;
inline long long Readin() {
	long long K = 0 , F = 1 ; char C = ' ' ;
	while( C < '0' or C > '9' ) F = C == '-' ? -1 : 1 , C = getchar() ;
	while( C <= '9' and C >= '0' ) K = ( K << 1 ) + ( K << 3 ) + C - '0' , C = getchar() ;
	return F * K ;
}
const int Mod = 1e9 + 7 ;
const int MaxN = 200000 + 10 ;
inline int Pow( int A , int B ) {
	register int Ans = 1 ;
	while( B ) {
		if( B & 1 ) Ans = 1ll * Ans * A % Mod ;
		A = 1ll * A * A % Mod ;
		B >>= 1 ;
	}
	return Ans ;
}
int N ;
int Miu[MaxN] , Prime[MaxN] , Pris ;
bool P[MaxN] ;
int main() {
	N = Readin() ;
	Miu[1] = 1 ;
	for(register int i = 1 ; ++i <= N ; ) {
		if( not P[i] ) {
			Prime[++Pris] = i ;
			Miu[i] = 1 ;
		}
		for(register int j = 0 ; ++j <= Pris ; ) {
			register int S = i * Prime[j] ;
			if( S > N ) break ;
			P[S] = 1 ;
			if( i % Prime[j] ) Miu[S] = -Miu[i] ;
			else {
				Miu[S] = 0 ;
				break ;
			}
		}
	}
	int Ansa = 1 , Ansb = 1 ;
	for(register int i = 1 ; ++i <= N ; ) {
		if( Miu[i] == 0 ) continue ;
		register int Thea = N / i , Theb = N - Thea ;
		if( Miu[i] == -1 ) Thea = Mod - Thea ;
		Ansa = ( 1ll * Ansa * Theb + 1ll * Ansb * Thea ) % Mod ;
		Ansb = 1ll * Ansb * Theb % Mod ;
	}
	return not printf( "%lld\n" , 1ll * Ansa * Pow( Ansb , Mod - 2 ) % Mod ) ;
}