/*
 * 	Author : Pallab
 * 	Prog: B.cxx
 *
 * 	Created on: 09:59:37 PM - 2013:3:13
 * "I have not failed, I have just found 10000 ways that won't work."
 */
#include <iostream>
#include <algorithm>
#include <vector>
#include <sstream>
#include <fstream>
#include <string>
#include <list>
#include <map>
#include <set>
#include <queue>
#include <deque>
#include <numeric>
#include <stack>
#include <functional>
#include <bitset>
#include <iomanip>

#include <ctime>
#include <cassert>
#include <cstdio>
#include <cmath>
#include <cstring>
#include <climits>
#include <cstring>
#include <cstdlib>

using namespace std;


#define foR(i,frm,upto) for( int i=(frm),_i=(upto);(i) < (_i) ; ++i )
#define foD(i,frm,upto) for( int i=(frm),_i=(upto);(i) >= (_i) ; --i )
#define foit(i, x) for (typeof((x).begin()) i = (x).begin(); i != (x).end(); i++)
#define Int long long
#define pb push_back
#define SZ(X) (int)(X).size()
#define LN(X) (int)(X).length()
#define mk make_pair
#define dbg(...)       printf(__VA_ARGS__)
#define SET( ARRAY , VALUE ) memset( ARRAY , VALUE , sizeof(ARRAY) )
#define line puts("")

inline void wait ( double seconds ) {
    double endtime = clock() + ( seconds * CLOCKS_PER_SEC );
    while ( clock() < endtime ) {
        ;
    }
}
template<class T>
    inline T fastIn() {
    register char c=0;
    register T a=0;
    bool neg=false;
    while ( c<33 ) c=getchar();
    while ( c>33 ) {
        if ( c=='-' ) {
            neg=true;
        } else {
            a= ( a*10 ) + ( c-'0' );
        }
        c=getchar();
    }
    return neg?-a:a;
}
int const mxn= ( int ) 1e6+5;
int N;
Int A[mxn],B[mxn];
inline void read() {
    N=fastIn<int>();
    foR ( i,0,N ) {
        A[i]=fastIn<Int>(),B[i]=fastIn<Int>();
    }
}
inline bool isPossible() {
    Int AA = accumulate ( A,A+N,0LL );
    Int BB = accumulate ( B,B+N,0LL );
    return abs ( AA-BB ) <= 500;
}
bool res[mxn];
inline void proc() {
    Int totalA= accumulate ( A,A+N,0LL );
    Int totalB= 0;
    SET ( res,true );
    for ( int i=0; i<N; ++i ) {
        if ( totalA>totalB ) {
            if ( totalA-totalB>500 ) {
                totalA-=A[i];
                totalB+=B[i];
                res[i]=false;
            }
        } else {
            if ( totalB-totalA>500 ) {
                puts ( "-1" ); return ;
            }
        }
    }
    foR(i,0,N){
        putchar( res[i]?'A':'G' );
    }line;
}
int main() {
    int kase = 1;
#if defined( xerxes_pc )
    if ( !freopen ( "in1", "r", stdin ) )
        puts ( "error opening file in " ), assert ( 0 );
    kase = 1;
#endif

    foR ( i,0,kase ) {
        read();
        proc();
    }
    return 0;
}
// kate: indent-mode cstyle; indent-width 4; replace-tabs on;