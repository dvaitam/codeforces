#include <stdio.h>
char s[1000001];
int main()
{
    int c=0,x;
    scanf("%s",s);
    while(s[1]){
        ++c;
        for(int i=x=0;s[i];++i)
            x+=s[i]-'0';
        sprintf(s,"%d",x);
    }
    printf("%d",c);
    return 0;
}