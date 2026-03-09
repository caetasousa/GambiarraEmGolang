export interface Service {
    id: string;
    nome: string;
    preco: number;
    duracao_padrao: number;
    categoria: string;
    image_url?: string;
}

export interface AgendaInterval {
    id: string;
    hora_inicio: string;
    hora_fim: string;
}

export interface AgendaDay {
    id: string;
    data: string;
    intervalos: AgendaInterval[];
}

export interface ProviderProfile {
    name: string;
    cpf: string;
    phone: string;
    email: string;
    avatarUrl: string;
    ativo?: boolean;
    role?: string;
}

export interface ProviderData {
    id: string;
    nome: string;
    email: string;
    telefone: string;
    cpf: string;
    ativo: boolean;
    image_url: string;
    catalogo: Service[];
    agenda: AgendaDay[];
    catalogo_ids?: string[];
}
