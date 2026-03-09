/**
 * Utilitários para upload de imagens no Supabase
 */

/** Configuração do Supabase */
export const SUPABASE_CONFIG = {
  url: import.meta.env.VITE_SUPABASE_URL,
  key: import.meta.env.VITE_SUPABASE_KEY,
  bucket: import.meta.env.VITE_SUPABASE_BUCKET_CATALOGOS,
} as const;

/**
 * Gera um hash SHA-256 do arquivo para criar um nome único
 * @param file - Arquivo para gerar o hash
 * @returns Promise com os primeiros 16 caracteres do hash em hexadecimal
 */
export async function generateFileHash(file: File): Promise<string> {
  const arrayBuffer: ArrayBuffer = await file.arrayBuffer();
  const hashBuffer: ArrayBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
  const hashArray: number[] = Array.from(new Uint8Array(hashBuffer));
  const hashHex: string = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex.substring(0, 16); // Usa apenas os primeiros 16 caracteres
}

/**
 * Valida se o arquivo é uma imagem válida
 * @param file - Arquivo a ser validado
 * @returns Objeto com resultado da validação e mensagem de erro (se houver)
 */
export function validateImageFile(file: File): { valid: boolean; error?: string } {
  // Valida tipo de arquivo
  if (!file.type.startsWith('image/')) {
    return {
      valid: false,
      error: "Por favor, selecione apenas arquivos de imagem.",
    };
  }

  // Valida tamanho máximo de 5MB
  const MAX_SIZE = 5 * 1024 * 1024; // 5MB em bytes
  if (file.size > MAX_SIZE) {
    return {
      valid: false,
      error: "A imagem deve ter no máximo 5MB.",
    };
  }

  return { valid: true };
}


/**
 * Faz upload de imagem para o Supabase Storage (função genérica)
 * @param file - Arquivo de imagem a ser enviado
 * @param bucketPath - Caminho do bucket (ex: "imagens/catalogos")
 * @returns Promise com a URL pública da imagem ou null em caso de erro
 */
async function uploadImageToSupabaseGeneric(file: File, bucketPath: string): Promise<string | null> {
  try {
    // Gera hash único para o nome do arquivo
    const hash = await generateFileHash(file);
    const fileExtension = file.name.split('.').pop()?.toLowerCase() || 'jpg';
    const hashedFilename = `${hash}.${fileExtension}`;

    console.log('Uploading image:', hashedFilename);
    console.log('File type:', file.type);
    console.log('File size:', file.size);
    console.log('Bucket path:', bucketPath);

    // Extrai o bucket raiz e o caminho (se houver subpasta)
    // Ex: "imagens/catalogos" -> bucket = "imagens", path = "catalogos"
    const bucketParts = bucketPath.split('/');
    const bucketName = bucketParts[0];
    const subPath = bucketParts.slice(1).join('/');
    const fullPath = subPath ? `${subPath}/${hashedFilename}` : hashedFilename;

    // Monta a URL de upload do Supabase Storage
    const uploadUrl = `${SUPABASE_CONFIG.url}/storage/v1/object/${bucketName}/${fullPath}`;

    console.log('Upload URL:', uploadUrl);

    // Faz requisição PUT para enviar o arquivo
    const uploadResponse = await fetch(uploadUrl, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${SUPABASE_CONFIG.key}`,
        'apikey': SUPABASE_CONFIG.key,
        'Content-Type': file.type,
      },
      body: file,
    });

    console.log('Upload response status:', uploadResponse.status);
    console.log('Upload response headers:', Object.fromEntries(uploadResponse.headers.entries()));

    // Verifica se o upload foi bem-sucedido
    if (!uploadResponse.ok) {
      const errorText = await uploadResponse.text();
      console.error('Upload error response:', errorText);

      // Tenta parsear a mensagem de erro como JSON
      try {
        const errorJson = JSON.parse(errorText);
        throw new Error(`Erro no upload: ${errorJson.message || errorJson.error || errorText}`);
      } catch {
        throw new Error(`Erro no upload (${uploadResponse.status}): ${errorText}`);
      }
    }

    const responseData = await uploadResponse.json();
    console.log('Upload success response:', responseData);

    // Constrói a URL pública para acessar a imagem
    const publicUrl = `${SUPABASE_CONFIG.url}/storage/v1/object/public/${bucketName}/${fullPath}`;

    console.log('Public URL:', publicUrl);

    return publicUrl;
  } catch (error) {
    console.error('Upload error:', error);
    throw error;
  }
}

/**
 * Faz upload de imagem de catálogo para o Supabase Storage
 * @param file - Arquivo de imagem a ser enviado
 * @returns Promise com a URL pública da imagem ou null em caso de erro
 */
export async function uploadImageToSupabase(file: File): Promise<string | null> {
  return uploadImageToSupabaseGeneric(file, SUPABASE_CONFIG.bucket);
}

/**
 * Faz upload de imagem de prestador para o Supabase Storage
 * @param file - Arquivo de imagem a ser enviado
 * @returns Promise com a URL pública da imagem ou null em caso de erro
 */
export async function uploadPrestadorImageToSupabase(file: File): Promise<string | null> {
  const prestadorBucket = import.meta.env.VITE_SUPABASE_BUCKET_PRESTADORES;
  return uploadImageToSupabaseGeneric(file, prestadorBucket);
}


/**
 * Cria um preview (base64) de uma imagem
 * @param file - Arquivo de imagem
 * @returns Promise com a URL base64 do preview
 */
export function createImagePreview(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result;
      if (typeof result === 'string') {
        resolve(result);
      } else {
        reject(new Error('Failed to create preview'));
      }
    };
    reader.onerror = () => reject(new Error('Failed to read file'));
    reader.readAsDataURL(file);
  });
}

/**
 * Deleta uma imagem do Supabase Storage
 * @param publicUrl - URL pública da imagem a ser deletada
 * @returns Promise indicando sucesso ou falha
 */
export async function deleteImageFromSupabase(publicUrl: string): Promise<boolean> {
  try {
    // Extrai o caminho completo do arquivo da URL pública
    // Ex: https://.../public/imagens/catalogos/abc.png -> imagens/catalogos/abc.png
    const urlParts = publicUrl.split('/storage/v1/object/public/');
    if (urlParts.length < 2) {
      console.error('Invalid image URL for deletion');
      return false;
    }

    const fullPath = urlParts[1]; // Ex: "imagens/catalogos/abc.png"
    
    if (!fullPath) {
      console.error('Invalid image URL for deletion');
      return false;
    }

    console.log('Deleting image from Supabase:', fullPath);

    // Extrai apenas o bucket raiz para a URL de delete
    const pathParts = fullPath.split('/');
    const bucketName = pathParts[0];
    const filePath = pathParts.slice(1).join('/');

    const deleteUrl = `${SUPABASE_CONFIG.url}/storage/v1/object/${bucketName}/${filePath}`;

    const response = await fetch(deleteUrl, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${SUPABASE_CONFIG.key}`,
        'apikey': SUPABASE_CONFIG.key,
      },
    });

    if (!response.ok) {
      console.error('Failed to delete image:', await response.text());
      return false;
    }

    console.log('Image deleted successfully');
    return true;
  } catch (error) {
    console.error('Error deleting image:', error);
    return false;
  }
}
