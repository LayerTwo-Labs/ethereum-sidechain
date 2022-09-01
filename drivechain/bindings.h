#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct Deposit {
  const char *address;
  uint64_t amount;
} Deposit;

typedef struct Deposits {
  struct Deposit *ptr;
  uintptr_t len;
} Deposits;

void init(const char *db_path,
          uintptr_t this_sidechain,
          const char *rpcuser,
          const char *rpcpassword);

void flush(void);

void attempt_bmm(const char *critical_hash, uint64_t amount);

uint32_t confirm_bmm(void);

bool verify_bmm(const char *main_block_hash, const char *critical_hash);

const char *get_prev_main_block_hash(const char *main_block_hash);

const char *format_deposit_address(const char *address);

struct Deposits get_deposit_outputs(void);

bool connect_block(struct Deposits deposits, bool just_check);

bool disconnect_block(struct Deposits deposits, bool just_check);

void free_string(const char *string);

void free_deposits(struct Deposits deposits);
