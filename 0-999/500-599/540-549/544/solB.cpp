#include <bits/stdc++.h>

using namespace std;



int main() 

{

    int n,k,i,j,z=0;

    char s[101][101];

    scanf("%d%d",&n,&k);

    if(((n*n)+1)/2<k)

    {

        printf("NO");

    }

    else

    {

        printf("YES\n");

        for(i=0;i<n;i++)

        {

            for(j=0;j<n;j++)

            {

                if((z+j)%2==0&&k>0)

                {

                    printf("L");

                    k--;

                }

                else

                {

                    printf("S");

                }

            }

            z++;

            printf("\n");

        }

    }

	return 0;

}