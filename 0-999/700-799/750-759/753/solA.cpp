#include <stdio.h>

int largest_discrete_sum(int *a, int j);

int main (void) {

    int i, idx = 0, v;
    scanf("%d",&v);
    int a[v];
    

    idx = largest_discrete_sum (a, v);
    
    printf("%d \n",idx-1);
for (i = 0; i < idx; i++)
        if (a[i])
        	printf("%d ",a[i]);
         
    return 0;
}

int largest_discrete_sum (int *a, int j)
{
    int n, sum = 0;
    for (n = 1;; n++) {
        a[n-1] = n, sum += n;
        if (n * (n + 1) / 2 > j)
            break;
    }
    a[sum - j - 1] = 0; 
    return n;
}