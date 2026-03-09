-- Seed para calculation_formulas
-- Fórmula padrão NR-1 com thresholds de risco e confiabilidade

-- Inserir fórmula padrão para cada partner existente
-- Ajuste o partner_id conforme necessário (1 é o padrão)
INSERT INTO calculation_formulas (
    partner_id, 
    version, 
    active,
    risk_low_max,
    risk_medium_max,
    reliability_acceptable_min,
    reliability_good_min,
    reliability_excellent_min,
    description,
    created_at
)
VALUES (
    1,                  -- partner_id (ajuste conforme seu partner)
    '1.0',              -- version
    true,               -- active
    1.5,                -- risk_low_max (scores <= 1.5 = baixo risco)
    2.5,                -- risk_medium_max (scores 1.5-2.5 = médio risco, >2.5 = alto)
    30,                 -- reliability_acceptable_min (>= 30% taxa de resposta)
    50,                 -- reliability_good_min (>= 50% taxa de resposta)
    70,                 -- reliability_excellent_min (>= 70% taxa de resposta)
    'Fórmula padrão baseada em estudos NR-1',
    NOW()
)
ON CONFLICT (partner_id, version) DO UPDATE SET
    active = EXCLUDED.active,
    risk_low_max = EXCLUDED.risk_low_max,
    risk_medium_max = EXCLUDED.risk_medium_max,
    reliability_acceptable_min = EXCLUDED.reliability_acceptable_min,
    reliability_good_min = EXCLUDED.reliability_good_min,
    reliability_excellent_min = EXCLUDED.reliability_excellent_min,
    description = EXCLUDED.description;

-- Verificar inserção
SELECT * FROM calculation_formulas WHERE partner_id = 1;
