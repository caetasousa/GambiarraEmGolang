package domain

type ServicoPrestado struct {
    ID                string
    Nome              string
    Descricao         string
    DuracaoEmMinutos  int // Essencial para calcular DataHoraFim
    PrecoCents        int // Usar inteiro (centavos) para evitar problemas de ponto flutuante
}