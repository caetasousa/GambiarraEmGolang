-- Cria o usuário administrador padrão do sistema como prestador
-- Email: admin@admin.com | Senha: admin123
INSERT INTO prestadores (id, nome, cpf, email, telefone, ativo, imagem_url, senha_hash, role)
VALUES (
    'admin000000000000001',
    'Administrador',
    '00000000000',
    'admin@admin.com',
    '00000000000',
    TRUE,
    '',
    '$2a$10$26/ZzhA.NoQn/QlhKNQ2oOd.BroZlTflrP8fru1ESnl68sTQi87R.',
    'admin'
)
ON CONFLICT DO NOTHING;
