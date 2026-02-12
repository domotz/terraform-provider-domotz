# Terraform Provider Domotz - Stato Lavori

## Completati

### Tier 1 (PR #1 - fix-tier1-pre-release)
- [x] 1.1 Pagination sulle chiamate List
- [x] 1.2 Gestione 404 in Delete/Read (NotFoundError)
- [x] 1.3 Versione nel Makefile (1.0.0)

### Tier 2 (PR #2 - fix-tier2-improvements)
- [x] 2.1 Retry con backoff esponenziale (429, 502, 503, 504)
- [x] 2.2 Context propagation in tutto il client
- [x] 2.3 Validazione input (hex color, port 1-65535, OneOf categories)
- [x] 2.4 User-Agent header

## Da Fare (Tier 3 - dopo merge PR)

| Item | Descrizione | Note |
|------|-------------|------|
| 3.1 | Race condition nel pattern list-and-find dopo Create | Richiede API che restituisca ID in POST response |
| 3.2 | Test di acceptance | Mock server o sandbox necessari |
| 3.3 | ID int32 → int64 refactoring | Breaking change, valutare per v2 |
| 3.4 | UpdateDevice non transazionale | Limitazione API, non risolvibile |
| 3.5 | Pubblicazione Terraform Registry | GPG + GitHub Actions |
| 3.6 | Badges reali nel README | Sostituire placeholder |

## PR in attesa di review

- PR #1: Tier 1 fixes (pagination, 404 handling, version)
- PR #2: Tier 2 improvements (retry, context, validation, user-agent)
- Reviewer: atondo
