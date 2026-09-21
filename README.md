# Cálculo Numérico (DIM0404): Avaliação 01

Implementação computacional das duas etapas de obtenção de raízes de funções
(Isolamento e Refinamento), feita em Go para a disciplina de Cálculo Numérico.

## 1.0. Ambiente computacional

- **Linguagem/Software:** Go `go1.27.1` (linux/amd64)
- **Sistema operacional:** 
- **Hardware (máquina de desenvolvimento):** 

## Estrutura do projeto

```
main.go            # monta os problemas, roda os métodos e imprime os resultados
types/types.go      # Problem, tipos auxiliares e a etapa de Isolamento (Bolzano + corolário)
graph/graph.go      # gera os gráficos das funções (item 1.1)
roots/
  bisection.go      # Método da Bissecção
  falsePosition.go  # Método da Falsa Posição
  newton.go         # Método de Newton
  secant.go         # Método da Secante
  fixedPoint.go     # Método do Ponto Fixo
  chaotic.go         # Estratégia proposta (item 1.3): corrida competitiva em memória compartilhada
```

## Como rodar

```bash
go run .                # executa, imprime os resultados e gera os gráficos em graphs/
go build ./...           # compila
go run -race .            # roda com o detector de data race do Go
```

## 1.1. Isolamento da raiz

Implementado em [`types/types.go`](types/types.go), função `isolateRec`:

- **Tabelamento:** percorre `[a, b]` em passos de tamanho `H` (via `nextInterval`),
  avaliando `f` em cada ponto.
- **Teorema de Bolzano:** se `f` é contínua em `[xmin, xmax]` e
  `f(xmin) * f(xmax) < 0`, existe pelo menos uma raiz nesse subintervalo: é essa
  troca de sinal que o tabelamento busca a cada passo.
- **Corolário de unicidade:** dentro de um subintervalo com raiz garantida, se
  `f'(x)` não muda de sinal nas pontas (`d1 * d2 >= 0`, ou seja, `f` é monótona
  ali), a raiz é única e o intervalo é aceito. Caso contrário, o subintervalo é
  reisolado recursivamente com o mesmo passo até isolar raízes únicas.
- **Gráfico:** [`graph/graph.go`](graph/graph.go) (usa `gonum.org/v1/plot`,
  biblioteca auxiliar, não faz parte do algoritmo em si) plota `f(x)` e marca
  os extremos dos intervalos isolados. Como em algumas funções do enunciado a
  escala do domínio completo esconde a troca de sinal, também é gerado um
  gráfico "zoom" ao redor do(s) intervalo(s) isolado(s). Salvos em `graphs/`
  ao rodar `go run .` (arquivos `f<Id>.png` e `f<Id>_zoom.png`).

## 1.2. Refinamento

Métodos implementados, todos com critério de parada por precisão
(`|f(x)| <= epsilon`, usado com `epsilon = 1e-8`, mais rigoroso que o `10⁻⁶`
mínimo exigido):

| Método         | Arquivo                                     |
|----------------|----------------------------------------------|
| Bissecção      | [`roots/bisection.go`](roots/bisection.go)         |
| Falsa Posição  | [`roots/falsePosition.go`](roots/falsePosition.go) |
| Ponto Fixo     | [`roots/fixedPoint.go`](roots/fixedPoint.go)       |
| Newton         | [`roots/newton.go`](roots/newton.go)               |
| Secante        | [`roots/secant.go`](roots/secant.go)               |

## 1.3. Estratégia proposta

[`roots/chaotic.go`](roots/chaotic.go) implementa uma "corrida caótica": Newton,
Secante e Ponto Fixo rodam em goroutines concorrentes compartilhando a mesma
variável de raiz aproximada em memória, e uma goroutine observadora retorna
assim que a precisão desejada é atingida: o primeiro método a "acertar" o
valor compartilhado decide o resultado.

## 1.4. Funções analisadas

| # | f(x)                                     | Intervalo   | h    |
|---|-------------------------------------------|-------------|------|
| 1 | 2x⁴ + 4x³ + 3x² − 10x − 15                | [0, 3]      | 0,6  |
| 2 | x⁵ − 2x⁴ − 9x³ + 22x² + 4x − 24            | [0, 5]      | 0,7  |
| 3 | 5x³ + x² − e^(1−2x) + cos(x) + 20          | [−5, 5]     | 0,5  |
| 4 | x·sen(x) + 4                                | [1, 5]      | 0,5  |

## Pendências para a entrega

- [ ] Comentar em português `Bisection`, `FalsePosition`, `Newton` e `Secant`
      (`fixedPoint.go`, `types.go` e `chaotic.go` já estão documentados)
- [ ] Preencher Sistema operacional / Hardware na seção 1.0
- [ ] Relatório resumido (entrega ii do descritivo)
- [ ] Lista de atividades por integrante (entrega iii do descritivo)

## Integrantes

- André Gomes
- Andriel Vinicius
- Maria Paz
