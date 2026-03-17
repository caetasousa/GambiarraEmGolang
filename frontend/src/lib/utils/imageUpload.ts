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
 */
export async function generateFileHash(file: File): Promise<string> {
  const arrayBuffer: ArrayBuffer = await file.arrayBuffer();
  const hashBuffer: ArrayBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
  const hashArray: number[] = Array.from(new Uint8Array(hashBuffer));
  const hashHex: string = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex.substring(0, 16);
}

/**
 * Valida se o arquivo é uma imagem válida
 */
export function validateImageFile(file: File): { valid: boolean; error?: string } {
  if (!file.type.startsWith('image/')) {
    return {
      valid: false,
      error: "Por favor, selecione apenas arquivos de imagem.",
    };
  }

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
 */
async function uploadImageToSupabaseGeneric(file: File, bucketPath: string): Promise<string | null> {
  const hash = await generateFileHash(file);
  const fileExtension = file.name.split('.').pop()?.toLowerCase() || 'jpg';
  const hashedFilename = `${hash}.${fileExtension}`;

  const bucketParts = bucketPath.split('/');
  const bucketName = bucketParts[0];
  const subPath = bucketParts.slice(1).join('/');
  const fullPath = subPath ? `${subPath}/${hashedFilename}` : hashedFilename;

  const uploadUrl = `${SUPABASE_CONFIG.url}/storage/v1/object/${bucketName}/${fullPath}`;

  const uploadResponse = await fetch(uploadUrl, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${SUPABASE_CONFIG.key}`,
      'apikey': SUPABASE_CONFIG.key,
      'Content-Type': file.type,
    },
    body: file,
  });

  if (!uploadResponse.ok) {
    const errorText = await uploadResponse.text();
    try {
      const errorJson = JSON.parse(errorText);
      throw new Error(`Erro no upload: ${errorJson.message || errorJson.error || errorText}`);
    } catch {
      throw new Error(`Erro no upload (${uploadResponse.status}): ${errorText}`);
    }
  }

  // Supabase retorna body vazio em uploads bem-sucedidos, não tenta parsear JSON
  const publicUrl = `${SUPABASE_CONFIG.url}/storage/v1/object/public/${bucketName}/${fullPath}`;
  return publicUrl;
}

/**
 * Faz upload de imagem de catálogo para o Supabase Storage
 */
export async function uploadImageToSupabase(file: File): Promise<string | null> {
  return uploadImageToSupabaseGeneric(file, SUPABASE_CONFIG.bucket);
}

/**
 * Faz upload de imagem de prestador para o Supabase Storage
 */
export async function uploadPrestadorImageToSupabase(file: File): Promise<string | null> {
  const prestadorBucket = import.meta.env.VITE_SUPABASE_BUCKET_PRESTADORES;
  return uploadImageToSupabaseGeneric(file, prestadorBucket);
}

/**
 * Cria um preview (base64) de uma imagem
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
 */
export async function deleteImageFromSupabase(publicUrl: string): Promise<boolean> {
  try {
    const urlParts = publicUrl.split('/storage/v1/object/public/');
    if (urlParts.length < 2) return false;

    const fullPath = urlParts[1];
    if (!fullPath) return false;

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

    return response.ok;
  } catch {
    return false;
  }
}
