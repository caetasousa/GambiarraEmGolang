/**
 * Utilitários para validação de formulários
 */

/** Tipo para os erros de validação */
export type ValidationErrors = Record<string, string>;

/** Interface para dados do formulário de serviço */
export interface ServiceFormData {
  nome: string;
  categoria: string;
  duracao_padrao: number;
  preco: string;
  description?: string; // Tornar opcional se não for usado em todo lugar
}

export interface EditServiceFormData {
  nome: string;
  categoria: string;
  duracao: number;
  preco: string;
}

/**
 * Valida os dados do formulário de serviço
 * @param data - Dados do formulário
 * @param hasImage - Se uma imagem foi selecionada
 * @returns Objeto com erros de validação (vazio se tudo válido)
 */
export function validateServiceForm(
  data: ServiceFormData,
  hasImage: boolean
): ValidationErrors {
  const errors: ValidationErrors = {};

  // Valida nome mínimo de 3 caracteres
  if (data.nome.length < 3) {
    errors.nome = "O nome deve ter pelo menos 3 caracteres.";
  }

  // Valida formato do preço (deve usar ponto, não vírgula)
  if (data.preco.includes(",")) {
    errors.preco = "Utilize pontos (.) em vez de vírgulas (,) para decimais.";
  }

  // Valida que o preço não seja negativo
  const numericPrice = parseFloat(data.preco) || 0;
  if (numericPrice < 0) {
    errors.preco = "O preço não pode ser negativo.";
  }

  // Valida que a duração não seja negativa
  if (data.duracao_padrao < 0) {
    errors.duracao = "A duração não pode ser negativa.";
  }

  // Valida que uma imagem foi selecionada (obrigatório)
  if (!hasImage) {
    errors.imagem = "Selecione uma imagem para o serviço.";
  }

  return errors;
}

/**
 * Converte preço de string para centavos
 * @param preco - Preço em formato string (ex: "50.00")
 * @returns Preço em centavos (ex: 5000)
 */
export function convertPriceTocents(preco: string): number {
  const numericPrice = parseFloat(preco) || 0;
  return Math.round(numericPrice * 100);
}

/**
 * Valida nome de categoria
 * @param name - Nome da categoria
 * @param existingCategories - Lista de categorias existentes
 * @returns Objeto com resultado da validação e mensagem de erro (se houver)
 */
export function validateCategoryName(
  name: string,
  existingCategories: string[]
): { valid: boolean; error?: string } {
  if (!name.trim()) {
    return {
      valid: false,
      error: "O nome da categoria é obrigatório.",
    };
  }

  if (existingCategories.includes(name.trim())) {
    return {
      valid: false,
      error: "Esta categoria já existe.",
    };
  }

  return { valid: true };
}

/**
 * Valida o formulário de edição de serviço
 * @param data - Dados do serviço
 * @param hasImage - Se existe imagem definida (nova ou pré-existente)
 * @returns Erros de validação
 */
export function validateEditServiceForm(
  data: EditServiceFormData,
  hasImage: boolean
): ValidationErrors {
  const errors: ValidationErrors = {};

  // Nome: min 3, max 100
  if (!data.nome || data.nome.length < 3) {
    errors.nome = "O nome deve ter pelo menos 3 caracteres.";
  } else if (data.nome.length > 100) {
    errors.nome = "O nome deve ter no máximo 100 caracteres.";
  }

  // Categoria: min 3, max 50
  if (!data.categoria || data.categoria.length < 3) {
    errors.categoria = "A categoria deve ter pelo menos 3 caracteres.";
  } else if (data.categoria.length > 50) {
    errors.categoria = "A categoria deve ter no máximo 50 caracteres.";
  }

  // Duração: > 1
  if (data.duracao <= 1) {
    errors.duracao = "A duração deve ser maior que 1 minuto.";
  }

  // Preço: >= 0 e válido
  const priceFloat = parseFloat(data.preco.replace(",", "."));
  if (isNaN(priceFloat) || priceFloat < 0) {
    errors.preco = "O preço deve ser um valor válido e não negativo.";
  }

  // Imagem obrigatória
  if (!hasImage) {
    errors.imagem = "A imagem é obrigatória.";
  }

  return errors;
}
