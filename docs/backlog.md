# Backlog

Lista de tarefas e melhorias do projeto goapi.

## Done

- [x] Refatorar entidades para usar `time.Time` em campos de data
- [x] Aplicar naming convention do banco (`tb_`, `pk_`, `tx_`, `nr_`, `ts_`, `fk_`)
- [x] Corrigir `Save*` para usar o ID da entidade em vez de gerar um novo UUID
- [x] Adicionar operações Update e Delete no Category e Product
- [x] Criar handlers CRUD para Category e Product
- [x] Adicionar rota de produtos por categoria (`GET /api/categories/{id}/products`)
- [x] Conectar database ao server via `lib/pq`

## Todo

- [x] Migrar `entity.go` para arquivos separados (`category.go`, `product.go`)
- [ ] Criar schema SQL de inicialização do banco
- [ ] Adicionar validação de input nos handlers
