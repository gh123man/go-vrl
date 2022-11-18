#ifndef VRL_H_INCLUDED
#define VRL_H_INCLUDED

typedef struct CResult {
    void* value;
    char* error;
} CResult;

char* run_vrl_c(char* str, void* program);

void* kind_bytes();
void* kind_object();
// TODO: add the other kinds


void* external_env_default();
void* external_env(void* target, void* metadata);

CResult compile_vrl(char* str);
CResult compile_vrl_with_external(char* str, void* externalEnv);

void* new_runtime();
CResult runtime_resolve(void* runtime, void* program, char* input);
void runtime_clear(void* runtime);
int runtime_is_empty(void* runtime);


#endif