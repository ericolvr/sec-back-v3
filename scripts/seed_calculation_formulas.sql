-- Seed para calculation_formulas
-- Fórmulas de cálculo para questionários NR-1

-- Fórmula padrão: média simples das respostas
INSERT INTO calculation_formulas (partner_id, name, description, formula_type, weight_config, created_at, updated_at)
VALUES 
(1, 'Média Simples', 'Calcula a média aritmética simples de todas as respostas', 'simple_average', '{}', NOW(), NOW()),
(1, 'Média Ponderada NR-1', 'Calcula média ponderada considerando peso das perguntas', 'weighted_average', '{"default_weight": 1.0}', NOW(), NOW()),
(1, 'Score Normalizado', 'Normaliza o score para escala 0-100', 'normalized', '{"min": 0, "max": 100}', NOW(), NOW());

-- Verificar inserção
SELECT * FROM calculation_formulas;
