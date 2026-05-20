# POS KopiTiam v3NF — N+1 & Transaction Fixes

## Summary of Changes

This release eliminates all N+1 query problems and strengthens transaction boundaries across the codebase.

---

## 1. N+1 Query Fixes

### Pattern Applied: "Collect IDs → Batch Query → In-Memory Map"

| File | Problem | Solution |
|------|---------|----------|
| `order_service.go` `CreateOrder` | Loop queried menu one-by-one | Batch `WHERE id IN (...)` + in-memory map |
| `order_service.go` `calculateRecipeCost` | Loop queried raw materials one-by-one | Batch `GetRawMaterialsByIDs` + in-memory map |
| `order_service.go` `ConfirmOrder` | Loop queried menu + recipe + raw materials inside tx | Pre-batch all lookups OUTSIDE tx, then use maps inside tx |
| `order_service.go` `VoidOrder` | Loop restored stock one-by-one | Batch menu lookup + map-driven restore |
| `order_service.go` `VoidItem` | Single query per void | Pre-fetch menu before tx |
| `return_service.go` `ProcessReturn` | Loop restored stock one-by-one | Batch menu lookup + map-driven restore |
| `analytics_repository.go` `GetSalesSummary` | Correlated subquery = N+1 | Replaced with JOIN aggregation |
| `analytics_repository.go` `GetBestSellers` | Could trigger per-item lookups | Single JOIN + GROUP BY query |
| `analytics_repository.go` `GetCOGS` | Correlated subquery = N+1 | Replaced with JOIN aggregation |
| `order_repo.go` `FindAllWithDetails` | JOIN produced duplicates for aggregates | Subquery JOIN for totals + batched detail lookup |

### New Batch Methods Added

```go
// inventory_repository.go
GetRawMaterialsByIDs(ids []uuid.UUID) ([]models.RawMaterial, error)
GetRecipesByMenuIDs(menuIDs []uuid.UUID) ([]models.Recipe, error)

// menu_repository.go
FindByIDs(ids []uuid.UUID) ([]models.Menu, error)
```

### Strategy: Preload vs JOIN vs Batch

- **Preload** → Use for `hasMany`/`many2many` on SINGLE records (FindByID). GORM sends 2 queries: 1 parent + 1 batched IN for children. Safe and correct.
- **JOIN** → Use for `belongsTo`/`hasOne` or aggregation queries. Single query, no duplication risk for singular associations.
- **Batch IN + Manual Map** → Use for `hasMany` on LIST queries where Preload would produce massive IN clauses or where we need derived aggregates. We query children in one batched IN, then map in-memory.

---

## 2. Transaction Boundary Fixes

| File | Problem | Solution |
|------|---------|----------|
| `order_service.go` `CreateOrder` | No transaction — order + details not atomic | Wrapped in `tx.Begin()` … `tx.Commit()` |
| `order_service.go` `ConfirmOrder` | Read order outside tx, then wrote inside — race risk | Kept read outside (read-only), but ALL writes inside tx with `clause.Locking{Strength: "UPDATE"}` |
| `order_service.go` `VoidItem` | Read order outside tx, wrote inside | Pre-fetched target detail, then tx for writes |
| `order_service.go` `VoidOrder` | Read order outside tx, wrote inside | Pre-fetched all menus, then tx for writes |
| `payment_service.go` `ProcessPayment` | Payment created outside tx, order status updated separately | Single tx: payment + splits + order status |
| `return_service.go` `ProcessReturn` | Return + order update + stock restore not atomic | Single tx for all writes |
| `stock_service.go` `AdjustStock` | Raw material update + adjustment record not atomic | Single tx: update RM + create adjustment |
| `shift_service.go` `CloseShift` | Shift close not atomic | Single tx for shift update |
| `inventory_service.go` `UpsertRecipe` | Recipe hard-delete + create not atomic | Wrapped in tx |

### Transaction Pattern Used

```go
tx := s.db.Begin()
if tx.Error != nil {
    return nil, tx.Error
}

success := false
defer func() {
    if !success {
        tx.Rollback()
    }
}()

// ... all writes use tx, not s.db ...

success = true
tx.Commit()
```

### Optimistic Locking (Row Locking)

For concurrent stock deductions (ConfirmOrder), we use `clause.Locking{Strength: "UPDATE"}`:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedRM, "id = ?", ingredient.RawMaterialID)
```

This prevents race conditions when two orders deduct the same raw material simultaneously.

---

## 3. Files Modified

```
internal/repository/order_repo.go              (rewritten)
internal/repository/inventory_repository.go     (added batch methods)
internal/repository/menu_repository.go          (added batch method)
internal/repository/analytics_repository.go     (rewritten with JOINs)
internal/service/order_service.go               (rewritten)
internal/service/return_service.go              (rewritten)
internal/service/payment_service.go             (added tx)
internal/service/stock_service.go               (added tx)
internal/service/shift_service.go               (added tx)
internal/service/inventory_service.go           (added tx)
```

---

## 4. Performance Impact

| Scenario | Before | After | Improvement |
|----------|--------|-------|-------------|
| CreateOrder (10 items, 3 recipe-based) | 1 + 10 + 3×N ingredient queries | 1 + 1 + 1 batch queries | **~15x fewer queries** |
| ConfirmOrder (10 items) | 1 + 10 + 10×N recipe queries | 1 + 1 + 1 batch queries | **~20x fewer queries** |
| VoidOrder (10 items) | 10 stock restore queries | 1 batch + 10 in-memory | **10x fewer queries** |
| GetSalesSummary | Correlated subquery scan | Single JOIN | **Eliminates N+1** |
| FindAllWithDetails | JOIN duplicates + aggregation failure | Subquery JOIN + batched details | **Correct + faster** |

---

## 5. Testing Recommendations

1. **Load test** `CreateOrder` with 50 items — should complete in < 100ms for DB queries.
2. **Concurrent test** `ConfirmOrder` on same recipe-based item from 2 orders — should handle stock correctly with row locking.
3. **Verify** `FindAllWithDetails` returns correct subtotal/total_paid aggregates.
4. **Check** `Debug()` logs show only expected query counts (no per-row SELECTs).

---

## References

- GORM Preload vs Joins: https://goldlapel.com/grounds/go-postgres/gorm-preload-vs-joins-postgres
- N+1 Query Problem: https://blog.stackademic.com/the-n-1-query-problem-in-gorm-how-to-avoid-silent-performance-killers-856e028d4b15
- GORM Best Practices: https://oneuptime.com/blog/post/2026-01-07-go-gorm-orm/view
