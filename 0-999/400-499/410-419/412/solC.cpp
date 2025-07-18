#include <iostream>

#include <string>

#include <algorithm>

#include <bits/stdc++.h>

#include <stdio.h>

#include <stdlib.h>

using namespace std ;

const int Max = 10055 ;

long long gcd( long long a , long long b )

{

    if ( b == 0 )

    return a ;

    else

        return  gcd ( b , a % b ) ;

}

 long long lcm ( long long  a , long long b )

{

    return ( a* ( b / gcd ( a , b ) ) ) ;

}

char c[100005] , s [100005] ;

int n , i , k ;

int main()

{

       cin >> n ;

       for ( i = 0 ; i < 100005 ; i++ )

        c[i] = '?' ;

       while ( n -- )

       {

           scanf ( "%s" , s ) ;

           k = strlen (s) ;

        for ( int i = 0 ; i < k ; i++ )

           {

            if ( s[i] == '?')

                continue ;

             if ( c[i] == '?')

                c[i] = s[i] ;

             else if (c[i] != s[i] && s[i] != '?')

                c[i] = 'T' ;

           }

       }

       for (  i = 0 ; i < k ; i++ )

       {

           if ( c[i] == '?' )

            c[i] = 'x' ;

            else

                if ( c[i] == 'T' )

                c[i] = '?' ;

       }

       for (  i = 0 ; i < k ; i++ )

        printf ( "%c" , c[i] ) ;

      return 0 ;

}