package domain

type Prestador struct {
	ID           string
	Nome         string
	Cpf          CPF
	Email        string
	Telefone     string
	Ativo        bool
	ImagemUrl    string
	PasswordHash string // Senha hasheada com bcrypt
	Role         Role   // Sempre RolePrestador
	Catalogo     []Catalogo
	Agenda       []AgendaDiaria
}

func NovoPrestador(id, nome, cpf, email, telefone, passwordHash string, imagem string, catalogos []Catalogo) (*Prestador, error) {
	voCPF, err := NovoCPF(cpf)
	if err != nil {
		return nil, err
	}

	return &Prestador{
		ID:           id,
		Nome:         nome,
		Cpf:          voCPF,
		Email:        email,
		Telefone:     telefone,
		Ativo:        true,
		ImagemUrl:    imagem,
		PasswordHash: passwordHash, // Recebe o hash já pronto
		Role:         RolePrestador,
		Catalogo:     catalogos,
		Agenda:       []AgendaDiaria{},
	}, nil
}

func (p *Prestador) AdicionarAgenda(agenda *AgendaDiaria) error {
	if !p.Ativo {
		return ErrPrestadorInativo
	}

	for _, a := range p.Agenda {
		if a.Data == agenda.Data {
			return ErrAgendaDuplicada
		}
	}

	p.Agenda = append(p.Agenda, *agenda)
	return nil
}

func (p *Prestador) RemoverAgenda(data string) error {
	if !p.Ativo {
		return ErrPrestadorInativo
	}

	// Buscar índice da agenda
	indice := -1
	for i, agenda := range p.Agenda {
		if agenda.Data == data {
			indice = i
			break
		}
	}

	// Se não encontrou, retorna erro
	if indice == -1 {
		return ErrAgendaNaoEncontrada
	}

	// Remove agenda do slice
	p.Agenda = append(p.Agenda[:indice], p.Agenda[indice+1:]...)

	return nil
}