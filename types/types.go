package types

type Function func(float64) float64
type Interval []float64
type RootIntervalList []Interval

type Problem struct {
	Function  Function
	GFunction Function
	Start     float64
	End       float64
	Id        uint8
	H         float64
}

type RootMethod func(inst Problem, epsilon float64) (float64, uint16)

type NamedMethod struct {
	Name   string
	Method RootMethod
}

type RaceResult struct {
	Root    float64
	Counter uint16
}

type ResultStr struct {
	Root float64
	Iter uint16
}

// IsolateRoots executa a etapa de isolamento sobre o problema (fase 1 do
// método numérico: localizar subintervalos de [Start, End] que contêm raiz,
// antes de refinar com Bissecção/Falsa Posição/Newton/Secante/Ponto Fixo).
//
// Entrada: nenhuma além do próprio Problem (usa self.Start, self.End,
// self.Function e self.H como tamanho do passo do tabelamento).
//
// Saída: a lista de subintervalos [xmin, xmax] em que foi verificada a
// existência de pelo menos uma raiz isolada (ver isolateRec).
func (self Problem) IsolateRoots() RootIntervalList {
	return isolateRec(self.Start, self.End, self.Function, self.H)
}

// Deriv estima f'(x) no ponto x, usada para aplicar o corolário de raiz
// única em isolateRec e para montar λ no método do Ponto Fixo.
//
// Entrada: x (ponto onde a derivada é avaliada).
// Saída: a estimativa de f'(x).
func (self Problem) Deriv(x float64) float64 {
	return deriv(x, self.Function)
}

// isolateRec varre [a,b] em passos de tamanho "step" (tabelamento) e aplica
// o Teorema de Bolzano: se f é contínua em [xmin,xmax] e f(xmin)*f(xmax) < 0,
// então existe pelo menos uma raiz ξ em (xmin,xmax) tal que f(ξ) = 0. Cada
// troca de sinal encontrada entre dois pontos consecutivos do tabelamento
// (feito por nextInterval) marca um desses subintervalos candidatos.
//
// Em seguida aplica o corolário sobre unicidade da raiz: se f'(x) não muda
// de sinal em (xmin,xmax) (ou seja, f é monótona ali, verificado por
// d1*d2 >= 0 nas pontas do intervalo), a raiz naquele subintervalo é única e
// o intervalo é aceito direto. Caso contrário (f' muda de sinal, f não é
// monótona), pode haver mais de uma raiz ali, então o próprio subintervalo é
// recursivamente reisolado com o mesmo passo até isolar raízes únicas.
//
// Entrada: a, b (extremos do intervalo de busca), f (função a analisar),
// step (amplitude h do tabelamento).
//
// Saída: lista de subintervalos [xmin, xmax], cada um contendo exatamente
// uma raiz isolada de f.
func isolateRec(a, b float64, f Function, step float64) RootIntervalList {
	var collection = RootIntervalList{}
	var curr float64 = a

	for curr < b {
		sequence := nextInterval(curr, b, f, step)
		curr = sequence[len(sequence)-1] // atualiza pro elemento que parou
		n := len(sequence)
		if n == 1 {
			continue
		}
		xmin := sequence[n-2]
		xmax := sequence[n-1]
		// Teorema de Bolzano -> f(xmin)*f(xmax) < 0 garante ao menos 1 raiz em (xmin, xmax)
		// Se não houver troca de sinal, não há garantia de raiz e o trecho é descartado
		if f(xmin)*f(xmax) < 0 {
			d1 := deriv(xmin, f)
			d2 := deriv(xmax, f)
			if d1*d2 >= 0 { // corolário: derivada de sinal constante -> raiz única, boa
				collection = append(collection, Interval{xmin, xmax})
			} else {
				subintervals := isolateRec(xmin, xmax, f, step)
				collection = append(collection, subintervals...)
			}
		}
	}

	return collection
}

// deriv estima a derivada de f no ponto x pelo método das diferenças finitas
// progressivas: f'(x) ≈ (f(x+h) - f(x)) / h, com h pequeno fixo.
//
// Entrada: x (ponto de avaliação), f (função a derivar).
// Saída: a estimativa de f'(x).
func deriv(x float64, f Function) float64 {
	const h float64 = 1e-2
	return (f(x+h) - f(x)) / h
}

// nextInterval faz o tabelamento de f a partir de "a", andando em passos de
// tamanho "step" até encontrar uma troca de sinal em relação a f(a) (condição
// do Teorema de Bolzano) ou até alcançar "b". O último par de pontos do
// vetor retornado é o subintervalo candidato a conter raiz, consumido por
// isolateRec.
//
// Entrada: a, b (extremos do intervalo restante a varrer), f (função
// tabelada), step (amplitude h do passo).
//
// Saída: a sequência de pontos [a, a+step, a+2*step, ...] percorridos até a
// troca de sinal (ou até b), usada tanto para o tabelamento quanto para
// achar o subintervalo [xmin, xmax] onde f muda de sinal.
func nextInterval(a, b float64, f Function, step float64) Interval {
	var slice Interval
	curr := a
	for curr < b && f(a)*f(curr) >= 0 {
		slice = append(slice, curr)
		curr += step
	}

	if curr <= b {
		slice = append(slice, curr)
	} else {
		slice = append(slice, b)
	}

	return slice
}
