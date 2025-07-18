#include <cstdio>
#include <algorithm>
using namespace std;

int main()
{
    int n,k;
    scanf("%d%d",&n,&k);
    char str[200000];
    scanf("%s",str);

    int predict=0;
    for(int i=0;i<n;++i) predict+=max('z'-str[i],str[i]-'a');

    if(predict < k) printf("-1");

    else
    {
        int i=0;
        while(true)
        {
            int temp1='z'-str[i];
            int temp2=str[i]-'a';
            if(temp1 >= k)
            {
                str[i]+=k;
                break;
            }

            else if(temp2 >= k)
            {
                str[i]-=k;
                break;
            }

            else
            {
                int temp1='z'-str[i];
                int temp2=str[i]-'a';

                if(temp1 > temp2)
                {
                    str[i]='z';
                    k-=temp1;
                }

                else
                {
                    str[i]='a';
                    k-=temp2;
                }

                ++i;
            }
        }
        printf("%s",str);
    }

    return 0;
}