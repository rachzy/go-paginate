# Go Paginate v3 Client

O **Go Paginate Client** é uma biblioteca cliente que permite a outras aplicações Go gerarem facilmente query string parameters compatíveis com a biblioteca go-paginate. Ele fornece uma API fluente e type-safe para construir URLs com parâmetros de paginação, busca e filtros complexos.

## 🚀 Instalação

```bash
go get github.com/rachzy/go-paginate/v3
```

## 📖 Uso Básico

```go
package main

import (
    "fmt"
    "github.com/rachzy/go-paginate/v3/client"
)

func main() {
    // Criar um novo cliente
    c := client.New("https://api.example.com/users")
    
    // Construir URL com paginação básica
    url := c.Page(2).Limit(25).BuildURL()
    fmt.Println(url) // https://api.example.com/users?limit=25&page=2
}
```

## 🔧 Funcionalidades

### Paginação Básica

```go
client := client.New("https://api.example.com/users")
url := client.
    Page(2).        // Página 2
    Limit(25).      // 25 itens por página
    BuildURL()
```

### Busca e Campos de Busca

```go
url := client.
    Search("john").                           // Termo de busca
    SearchFields("name", "email", "username"). // Campos para buscar
    BuildURL()
```

### Ordenação

```go
url := client.
    Sort("name", "-created_at").  // Ordenar por nome (asc) e data de criação (desc)
    BuildURL()

// Ou usando métodos separados
url := client.
    SortColumns("name", "created_at").
    SortDirections("asc", "desc").
    BuildURL()
```

### Seleção de Colunas

```go
url := client.
    Columns("id", "name", "email", "created_at"). // Selecionar apenas essas colunas
    BuildURL()
```

## 🎯 Filtros Avançados

### Filtros LIKE

```go
url := client.
    Like("name", "john", "jane").           // name LIKE 'john' OR name LIKE 'jane'
    LikeOr("status", "active", "pending").  // status LIKE 'active' OR status LIKE 'pending'
    LikeAnd("description", "go", "lang").   // description LIKE 'go' AND description LIKE 'lang'
    BuildURL()
```

### Filtros de Igualdade

```go
url := client.
    Eq("department", "IT").                 // department = 'IT'
    EqOr("age", 25, 30, 35).               // age = 25 OR age = 30 OR age = 35
    EqAnd("status", "active", "verified").  // status = 'active' AND status = 'verified'
    BuildURL()
```

### Filtros de Comparação

```go
url := client.
    Gte("age", 18).           // age >= 18
    Gt("score", 80).          // score > 80
    Lte("salary", 100000).    // salary <= 100000
    Lt("experience", 5).      // experience < 5
    BuildURL()
```

### Filtros IN e NOT IN

```go
url := client.
    In("role", "admin", "manager", "user").     // role IN ('admin', 'manager', 'user')
    NotIn("status", "deleted", "banned").       // status NOT IN ('deleted', 'banned')
    BuildURL()
```

### Filtro BETWEEN

```go
url := client.
    Between("salary", 50000, 150000).  // salary BETWEEN 50000 AND 150000
    Between("age", 25, 65).            // age BETWEEN 25 AND 65
    BuildURL()
```

### Filtros NULL

```go
url := client.
    IsNull("deleted_at").              // deleted_at IS NULL
    IsNotNull("email", "phone").       // email IS NOT NULL AND phone IS NOT NULL
    BuildURL()
```

## 🔄 Funcionalidades Avançadas

### Clonagem de Cliente

```go
// Criar um cliente base com parâmetros comuns
baseClient := client.New("https://api.example.com/users")
baseClient.Limit(25).Columns("id", "name", "email")

// Clonar para diferentes casos de uso
activeUsers := baseClient.Clone().Eq("status", "active")
inactiveUsers := baseClient.Clone().Eq("status", "inactive")

activeURL := activeUsers.BuildURL()
inactiveURL := inactiveUsers.BuildURL()
```

### Construção a partir de URL Existente

```go
// Partir de uma URL existente
existingURL := "https://api.example.com/users?page=1&limit=10"
client, err := client.NewFromURL(existingURL)
if err != nil {
    log.Fatal(err)
}

// Adicionar mais parâmetros
newURL := client.Page(2).Eq("status", "active").BuildURL()
```

### Reset e Manipulação de Parâmetros

```go
client := client.New("https://api.example.com/users")
client.Page(1).Limit(10).Search("test")

// Limpar todos os parâmetros
client.Reset()

// Remover parâmetro específico
client.RemoveParam("search")

// Adicionar parâmetros customizados
client.SetCustomParam("custom", "value")
client.AddCustomParam("multi", "value1")
client.AddCustomParam("multi", "value2")
```

### Obter Apenas Query String

```go
client := client.New("") // URL base vazia
queryString := client.
    Page(2).
    Limit(50).
    Search("golang").
    BuildQueryString()

// Usar com diferentes URLs base
url1 := fmt.Sprintf("https://api1.com/posts?%s", queryString)
url2 := fmt.Sprintf("https://api2.com/articles?%s", queryString)
```

## 🌐 Integração com HTTP Clients

### Com net/http

```go
import (
    "net/http"
    "github.com/rachzy/go-paginate/v3/client"
)

func fetchUsers() {
    // Construir URL
    paginateClient := client.New("https://api.example.com/users")
    url := paginateClient.
        Page(1).
        Limit(10).
        Eq("status", "active").
        BuildURL()
    
    // Fazer requisição HTTP
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    // Processar resposta...
}
```

### Com outros HTTP clients

```go
// Obter parâmetros como url.Values
client := client.New("https://api.example.com/users")
params := client.Page(1).Limit(10).GetParams()

// Usar com qualquer biblioteca HTTP que aceite url.Values
// req.URL.RawQuery = params.Encode()
```

## 🎛️ Opções Especiais

```go
url := client.
    Vacuum(true).     // Habilitar modo vacuum
    NoOffset(true).   // Habilitar modo no offset
    BuildURL()
```

## 📝 Exemplo Completo

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/rachzy/go-paginate/v3/client"
)

func main() {
    // Criar cliente com filtros complexos
    c := client.New("https://api.example.com/users")
    
    url := c.
        Page(2).
        Limit(25).
        Search("john").
        SearchFields("name", "email").
        Sort("name", "-created_at").
        LikeOr("status", "active", "pending").
        EqOr("age", 25, 30, 35).
        Gte("created_at", "2023-01-01").
        Lt("score", 100).
        In("department", "IT", "HR").
        IsNotNull("email").
        Vacuum(true).
        BuildURL()
    
    fmt.Println("Generated URL:", url)
    
    // Fazer requisição HTTP
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    fmt.Println("Response Status:", resp.Status)
}
```

## 🧪 Testes

Para executar os testes do cliente:

```bash
cd v3/client
go test -v
```

## 📚 Compatibilidade

O cliente gera query strings totalmente compatíveis com:
- go-paginate v3 `BindQueryParams`
- go-paginate v3 `BindQueryStringToStruct`
- Todos os filtros e operadores suportados pela biblioteca principal

## 🔗 Links Úteis

- [Documentação Principal](../README.md)
- [Exemplos de Bind](../BIND_README.md)
- [Exemplos de Filtros](../FILTER_README.md)
- [Exemplo de Uso do Cliente](../examples/client/main.go)