#include<stdio.h>

int main(){
    int a,b,c;
    double d;
    scanf("%d%d%d",&a,&b,&c);
    d=(double)(a+b+c);
    printf("%.10lf\n",d/3.0);
    return 0;
}