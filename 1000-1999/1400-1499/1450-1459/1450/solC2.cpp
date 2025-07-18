#include <cstdio>



const int MAXN = 300;

int T , n , k; char s[ MAXN + 5 ][ MAXN + 5 ];



bool Solve( char c0 , char c1 , char c2 ) {

    int cnt = 0;

    for( int i = 1 ; i <= n ; i ++ )

        for( int j = 1 ; j <= n ; j ++ ) if( s[ i ][ j ] != '.' ) {

            int d = ( i + j ) % 3;

            if( d == 0 ) cnt += c0 != '*' && s[ i ][ j ] != c0;

            if( d == 1 ) cnt += c1 != '*' && s[ i ][ j ] != c1;

            if( d == 2 ) cnt += c2 != '*' && s[ i ][ j ] != c2;

        }

    // printf("%d\n", cnt );

    if( cnt > k / 3 ) return 0;

    for( int i = 1 ; i <= n ; i ++ ) {

        for( int j = 1 ; j <= n ; j ++ ) {

            if( s[ i ][ j ] != '.' ) {

                int d = ( i + j ) % 3;

                if( d == 0 ) putchar( c0 == '*' ? s[ i ][ j ] : c0 );

                if( d == 1 ) putchar( c1 == '*' ? s[ i ][ j ] : c1 );

                if( d == 2 ) putchar( c2 == '*' ? s[ i ][ j ] : c2 );

            }

            else putchar('.');

        }

        puts("");

    }

    return 1;

}



int main( ) {

    scanf("%d",&T);

    while( T -- ) {

        scanf("%d",&n); k = 0;

        for( int i = 1 ; i <= n ; i ++ ) {

            scanf("%s", s[ i ] + 1 );

            for( int j = 1 ; j <= n ; j ++ ) if( s[ i ][ j ] != '.' ) k ++;

        }

        if( Solve( '*' , 'O' , 'X' ) ) continue;

        if( Solve( '*' , 'X' , 'O' ) ) continue;

        if( Solve( 'O' , '*' , 'X' ) ) continue;

        if( Solve( 'X' , '*' , 'O' ) ) continue;

        if( Solve( 'O' , 'X' , '*' ) ) continue;

        if( Solve( 'X' , 'O' , '*' ) ) continue;

    }

    return 0;

}



/*

黑白灰染色

1.黑为x,白为x，灰为 o

2.黑为x,白为o，灰为 x

3.黑为x,白为o，灰为 o

...



六种情况取最小，这样做是 k/2 的



选一个点最多的颜色，令它不变

那么剩下至多 2/3 k 个点

枚举颜色做到 1/3 k

*/