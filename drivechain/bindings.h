#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct Block {
  const char *data;
  int64_t time;
  const char *main_block_hash;
} Block;

void test_function(void);

void init(const char *db_path,
          uintptr_t this_sidechain,
          const char *rpcuser,
          const char *rpcpassword);

void flush(void);

void attempt_bmm(const char *critical_hash, const char *block_data, uint64_t amount);

const struct Block *confirm_bmm(void);

bool verify_bmm(const char *main_block_hash, const char *critical_hash);

const char *get_prev_main_block_hash(const char *main_block_hash);
