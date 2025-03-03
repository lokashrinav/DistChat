#include "encryption.h"
#include <cstring>
#include <cstdlib>
extern "C" {
char* en(const char* a, const char* b) {
    size_t i = 0, n = std::strlen(a), m = std::strlen(b);
    char* c = (char*)std::malloc(n + 1);
    for(i = 0; i < n; i++){
        c[i] = a[i] ^ b[i % m];
    }
    c[n] = '\0';
    return c;
}
char* de(const char* a, const char* b) {
    return en(a, b);
}
void fr(char* a) {
    std::free(a);
}
}
