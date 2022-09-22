#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct WithdrawalAddress {
  uint8_t address[20];
} WithdrawalAddress;

typedef struct Withdrawal {
  const char *id;
  uint8_t address[20];
  uint64_t amount;
  uint64_t fee;
} Withdrawal;

typedef struct Withdrawals {
  struct Withdrawal *ptr;
  uintptr_t len;
} Withdrawals;

typedef struct Deposit {
  const char *address;
  uint64_t amount;
} Deposit;

typedef struct Deposits {
  struct Deposit *ptr;
  uintptr_t len;
} Deposits;

typedef struct Refund {
  const char *id;
  uint64_t amount;
} Refund;

typedef struct Refunds {
  struct Refund *ptr;
  uintptr_t len;
} Refunds;

void init(const char *db_path,
          uintptr_t this_sidechain,
          const char *rpcuser,
          const char *rpcpassword);

void flush(void);

void attempt_bmm(const char *critical_hash, const char *prev_main_block_hash, uint64_t amount);

uint32_t confirm_bmm(void);

bool verify_bmm(const char *main_block_hash, const char *critical_hash);

const char *get_prev_main_block_hash(const char *main_block_hash);

const char *get_mainchain_tip(void);

const char *format_deposit_address(const char *address);

bool create_deposit(const char *address, uint64_t amount, uint64_t fee);

struct WithdrawalAddress get_new_mainchain_address(void);

const char *format_mainchain_address(struct WithdrawalAddress dest);

bool attempt_bundle_broadcast(void);

struct Withdrawals get_unspent_withdrawals(void);

struct Deposits get_deposit_outputs(void);

bool connect_block(struct Deposits deposits,
                   struct Withdrawals withdrawals,
                   struct Refunds refunds,
                   bool just_check);

bool disconnect_block(struct Deposits deposits,
                      struct Withdrawals withdrawals,
                      struct Refunds refunds,
                      bool just_check);

bool is_outpoint_spent(const char *outpoint);

void free_string(const char *string);

void free_deposits(struct Deposits deposits);

void free_withdrawals(struct Withdrawals withdrawals);
