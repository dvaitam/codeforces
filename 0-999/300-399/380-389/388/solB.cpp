#include <algorithm>

#include <iostream>

#include <memory.h>

#include <iomanip>

#include <utility>

#include <cstdio>

#include <string>

#include <vector>

#include <deque>

#include <stack>

#include <queue>

#include <cmath>

#include <ctime>

#include <list>

#include <map>

#include <set>



#define fr first

#define sc second

#define mk make_pair

#define pb push_back

#define ll unsigned long long

#define all(s) s.begin(), s.end()



using namespace std;



const int N = 1*1e2+1;



bool edge[N][N];

int k;



inline void draw(int a, int b)

{

        edge[a][b] = edge[b][a] = 1;

}



main()

{

        scanf("%d", &k);



        for(int i = 0 ; i < 30 ; i ++)

                for(int j = 2 * i + 2 ; j < 2 * i + 4 ; j ++)

                        for(int l = 2 * i + 4 ; l < 2 * i + 6 ; l ++)

                                draw(j, l);



        for(int i = 64 ; i < 94 ; i ++)

                draw(i, i + 1);

        draw(0, 2);

        draw(0, 3);

        draw(94, 1);



        for(int i = 0 ; i < 30 ; i ++){

                if(k & (1 << i))

                        draw(2 * i + 2, 64 + i + 1);

        }



        puts("95");



        for(int i = 0 ; i < 95 ; i ++)

        {

                for(int j = 0 ; j < 95 ; j ++)

                        printf("%c", (edge[i][j]) ? 'Y' : 'N');

                puts("");

        }

}