ghz -i ../../third_party/googleapis/ --insecure -n 1 -c 1 \
  --proto ../../proto/ledger/ledger.proto \
  --call ledger.LedgerService.CreateTransaction \
  -d '[{"id":"{{.UUID}}", "entries":[{"id":"{{newUUID}}", "account_id":"assets:aaa:bbb:{{newUUID}}", "expected_version":"0", "operation":"OPERATION_DEBIT", "amount": 123}, {"id":"{{newUUID}}", "account_id":"assets:aaa:bbb:{{newUUID}}", "expected_version":"0", "operation":"OPERATION_CREDIT", "amount": 123}]}]' \
  0.0.0.0:3000
