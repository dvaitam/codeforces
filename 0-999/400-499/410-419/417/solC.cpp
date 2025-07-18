#include<cstdio>

using namespace std;

int n,k,t;

int main()

{

    scanf("%d %d",&n,&k);

if (n%2==0)

{

    if (k>n/2-1)

    {

        printf("-1");

        return 0;

    }

}

if (n%2!=0)

{

    if (k>n/2)

    {

        printf("-1");

        return 0;

    }

}

 printf("%d\n",n*k);

  for (int i=1;i<=n;i++)

       for (int j=i+1;j<=i+k;j++)

       {

           if (j>n)

              {

                 t=j-n;

                 printf("%d %d\n",i,t);

              }

     else

           printf("%d %d\n",i,j);

       }



return 0;

}