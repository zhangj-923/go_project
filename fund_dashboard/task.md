# Fund Dashboard Implementation Tasks

- [x] 1. Initialize Go module `fund_dashboard`
- [x] 2. Run `go mod tidy` after adding imports
- [x] 3. Create `model/fund.go`, `model/transaction.go`, `model/analytics.go`
- [x] 4. Create `storage/db.go` to initialize SQLite at `data/fund.db`
- [x] 5. Create `storage/fund.go` and `storage/transaction.go` for CRUD operations
- [x] 6. Create `service/nav_fetcher.go` using the East Money APIs
- [x] 7. Create `service/analytics.go` for calculating portfolio summary and PnL
- [x] 8. Create `handler/fund.go`, `handler/transaction.go`, `handler/analytics.go`
- [x] 9. Create `main.go` to wire everything up and start the server on `:8080`
- [x] 10. Initialize Vue 3 + TS + Vite frontend in `web`
- [x] 11. Implement dark hacker terminal theme in `src/style.css`
- [x] 12. Create Pinia stores (`fund`, `transaction`, `analytics`) and `useApi`
- [x] 13. Build left panel components (`BuySellForm`, `TransactionList`)
- [x] 14. Build right panel components (`SummaryCards`, `ValueDistribution`, `PnLCalendar`, `FundDetails`, `PnLLog`)
- [x] 15. Assemble the layout in `App.vue`
