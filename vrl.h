#ifndef VRL_H_INCLUDED
#define VRL_H_INCLUDED

typedef struct CResult {
    void* value;
    char* error;
} CResult;

char* run_vrl_c(char* str, void* program);

CResult compile_vrl(char* str);
void* new_runtime();
CResult runtime_resolve(void* runtime, void* program, char* input);


char* run_vrl_c(char* str, void* program);
void* compile_vrl_c(char* str);


#endif