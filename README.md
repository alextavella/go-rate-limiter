# Rate Limiter

### Getting started

```bash
go install github.com/rakyll/hey@latest
```

### Running

```bash
go run cmd/main.go
```

### Tests

```bash
hey -n 100 -c 20 -q 1 http://localhost:8080/
```

Onde:

```bash
-n 100  = Total de 100 requisições.
-c 20   = Até 20 requisições concorrentes.
-q 1    = Limitando a taxa para 1 requisições por segundo.
```

#### Teste 1 using 5 request per second

```bash
hey -n 10 -c 1 http://localhost:8080/
Status code distribution:
  [200] 5 responses
  [429] 5 responses
```

#### Test 1.1

```bash
hey -n 10 -c 1 -q 1 http://localhost:8080/
Status code distribution:
  [200] 10 responses
  [429] 0 responses
```

#### Test 2 (API-KEY: ABC) using 10 request per second

```bash
hey -n 10 -c 1 -H "API_KEY: ABC" http://localhost:8080/
Status code distribution:
  [200] 10 responses
  [429] 0 responses
```

#### Test 3 (API-KEY: DEF) using 8 request per second

```bash
hey -n 10 -c 1 -H "API_KEY: DEF" http://localhost:8080/
Status code distribution:
  [200] 8 responses
  [429] 2 responses
```
