-- NR-1 Backend v2 - Refactored Database Schema
-- Multi-tenancy: Partner (Consultoria RH) -> Companies (Clientes)

-- Drop tables if they exist (in correct order due to foreign keys)
DROP TABLE IF EXISTS risk_metrics CASCADE;
DROP TABLE IF EXISTS assessment_versions CASCADE;
DROP TABLE IF EXISTS action_plan_templates CASCADE;
DROP TABLE IF EXISTS risk_categories CASCADE;
DROP TABLE IF EXISTS action_plans CASCADE;
DROP TABLE IF EXISTS invitations CASCADE;
DROP TABLE IF EXISTS questionnaire_assignments CASCADE;
DROP TABLE IF EXISTS analytics_reports CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS answers CASCADE;
DROP TABLE IF EXISTS employee_submissions CASCADE;
DROP TABLE IF EXISTS questions CASCADE;
DROP TABLE IF EXISTS assessment_templates CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS partners CASCADE;

-- ============================================
-- PARTNERS (Tenant Raiz - Consultoria RH)
-- ============================================
CREATE TABLE partners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) UNIQUE NOT NULL,
    email VARCHAR(255),
    mobile VARCHAR(20),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_partners_cnpj ON partners(cnpj);
CREATE INDEX idx_partners_active ON partners(active);

-- ============================================
-- COMPANIES (Clientes do Partner)
-- ============================================
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) UNIQUE NOT NULL,
    email VARCHAR(255),
    mobile VARCHAR(20),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_companies_partner_id ON companies(partner_id);
CREATE INDEX idx_companies_cnpj ON companies(cnpj);
CREATE INDEX idx_companies_active ON companies(active);

-- ============================================
-- DEPARTMENTS
-- ============================================
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE INDEX idx_departments_partner_id ON departments(partner_id);
CREATE INDEX idx_departments_company_id ON departments(company_id);
CREATE INDEX idx_departments_active ON departments(active);

-- ============================================
-- EMPLOYEES
-- ============================================
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    mobile VARCHAR(20),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_employees_partner_id ON employees(partner_id);
CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_employees_department_id ON employees(department_id);
CREATE INDEX idx_employees_active ON employees(active);

-- ============================================
-- USERS (Usuários do sistema)
-- ============================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    type INT NOT NULL, -- 1 = master, 2 = client
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    UNIQUE(partner_id, mobile)
);

CREATE INDEX idx_users_partner_id ON users(partner_id);
CREATE INDEX idx_users_mobile ON users(mobile);
CREATE INDEX idx_users_active ON users(active);

-- ============================================
-- ASSESSMENT TEMPLATES (Questionários)
-- ============================================
CREATE TABLE assessment_templates (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version INT DEFAULT 1,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_assessment_templates_partner_id ON assessment_templates(partner_id);
CREATE INDEX idx_assessment_templates_active ON assessment_templates(active);

-- ============================================
-- ASSESSMENT VERSIONS (Auditoria/Versionamento)
-- ============================================
CREATE TABLE assessment_versions (
    id SERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    version INT NOT NULL,
    changes TEXT, -- JSON descrevendo mudanças
    created_by BIGINT, -- user_id
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_assessment_versions_template_id ON assessment_versions(template_id);
CREATE INDEX idx_assessment_versions_partner_id ON assessment_versions(partner_id);

-- ============================================
-- QUESTIONS
-- ============================================
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    question TEXT NOT NULL,
    type VARCHAR(50) NOT NULL, -- scale, multiple_choice, text, yes_no, number
    category VARCHAR(100), -- NR-1: Sobrecarga, Assédio, Autonomia, etc
    options TEXT, -- JSON array
    score_values TEXT, -- JSON array
    weight DECIMAL(5,2) DEFAULT 1.0,
    required BOOLEAN DEFAULT true,
    order_num INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE
);

CREATE INDEX idx_questions_partner_id ON questions(partner_id);
CREATE INDEX idx_questions_questionnaire_id ON questions(questionnaire_id);
CREATE INDEX idx_questions_category ON questions(category);

-- ============================================
-- EMPLOYEE SUBMISSIONS (Submissões de Questionários)
-- ============================================
CREATE TABLE employee_submissions (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    invitation_token VARCHAR(255) UNIQUE,
    status VARCHAR(50) DEFAULT 'pending', -- pending, in_progress, completed
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_employee_submissions_partner_id ON employee_submissions(partner_id);
CREATE INDEX idx_employee_submissions_company_id ON employee_submissions(company_id);
CREATE INDEX idx_employee_submissions_department_id ON employee_submissions(department_id);
CREATE INDEX idx_employee_submissions_employee_id ON employee_submissions(employee_id);
CREATE INDEX idx_employee_submissions_token ON employee_submissions(invitation_token);
CREATE INDEX idx_employee_submissions_status ON employee_submissions(status);

-- ============================================
-- ANSWERS
-- ============================================
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    submission_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    value TEXT NOT NULL,
    score INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (submission_id) REFERENCES employee_submissions(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE INDEX idx_answers_partner_id ON answers(partner_id);
CREATE INDEX idx_answers_submission_id ON answers(submission_id);
CREATE INDEX idx_answers_question_id ON answers(question_id);

-- ============================================
-- RISK METRICS (Métricas Pré-calculadas)
-- ============================================
CREATE TABLE risk_metrics (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    
    -- Métricas
    total_employees INT DEFAULT 0,
    total_submissions INT DEFAULT 0,
    completed_submissions INT DEFAULT 0,
    response_rate DECIMAL(5,2) DEFAULT 0,
    average_score DECIMAL(5,2) DEFAULT 0,
    risk_level VARCHAR(50), -- low, medium, high, critical
    reliability VARCHAR(50), -- low, medium, high
    can_calculate_risk BOOLEAN DEFAULT false,
    
    -- Scores por categoria (JSON)
    category_scores TEXT, -- JSONB: {"Sobrecarga": 7.5, "Assédio": 3.2, ...}
    
    -- Auditoria
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    
    -- Unique constraint: 1 métrica por departamento/questionário
    UNIQUE(partner_id, company_id, department_id, questionnaire_id)
);

CREATE INDEX idx_risk_metrics_partner_id ON risk_metrics(partner_id);
CREATE INDEX idx_risk_metrics_company_id ON risk_metrics(company_id);
CREATE INDEX idx_risk_metrics_department_id ON risk_metrics(department_id);
CREATE INDEX idx_risk_metrics_questionnaire_id ON risk_metrics(questionnaire_id);
CREATE INDEX idx_risk_metrics_risk_level ON risk_metrics(risk_level);

-- ============================================
-- QUESTIONNAIRE ASSIGNMENTS
-- ============================================
CREATE TABLE questionnaire_assignments (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_questionnaire_assignments_partner_id ON questionnaire_assignments(partner_id);
CREATE INDEX idx_questionnaire_assignments_questionnaire_id ON questionnaire_assignments(questionnaire_id);
CREATE INDEX idx_questionnaire_assignments_department_id ON questionnaire_assignments(department_id);

-- ============================================
-- INVITATIONS
-- ============================================
CREATE TABLE invitations (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    sent BOOLEAN DEFAULT false,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_invitations_partner_id ON invitations(partner_id);
CREATE INDEX idx_invitations_employee_id ON invitations(employee_id);
CREATE INDEX idx_invitations_token ON invitations(token);

-- ============================================
-- RISK CATEGORIES (Por Snapshot)
-- ============================================
CREATE TABLE risk_categories (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    snapshot_id BIGINT,
    category VARCHAR(100) NOT NULL,
    average_score DECIMAL(5,2) DEFAULT 0,
    risk_level VARCHAR(50),
    question_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_risk_categories_partner_id ON risk_categories(partner_id);
CREATE INDEX idx_risk_categories_snapshot_id ON risk_categories(snapshot_id);
CREATE INDEX idx_risk_categories_category ON risk_categories(category);

-- ============================================
-- ACTION PLAN TEMPLATES
-- ============================================
CREATE TABLE action_plan_templates (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    category VARCHAR(100),
    min_risk_level VARCHAR(50),
    title_template TEXT NOT NULL,
    description_template TEXT,
    priority VARCHAR(50),
    default_due_days INT DEFAULT 30,
    auto_create BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_action_plan_templates_partner_id ON action_plan_templates(partner_id);
CREATE INDEX idx_action_plan_templates_category ON action_plan_templates(category);

-- ============================================
-- ACTION PLANS
-- ============================================
CREATE TABLE action_plans (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    snapshot_id BIGINT,
    title TEXT NOT NULL,
    description TEXT,
    priority VARCHAR(50),
    status VARCHAR(50) DEFAULT 'pending',
    due_date DATE,
    completed_at TIMESTAMP,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_action_plans_partner_id ON action_plans(partner_id);
CREATE INDEX idx_action_plans_company_id ON action_plans(company_id);
CREATE INDEX idx_action_plans_department_id ON action_plans(department_id);
CREATE INDEX idx_action_plans_status ON action_plans(status);

-- ============================================
-- CALCULATION FORMULAS (Fórmulas de Cálculo Versionadas)
-- ============================================
CREATE TABLE calculation_formulas (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    version VARCHAR(10) NOT NULL,
    active BOOLEAN DEFAULT false,
    
    -- Thresholds de risco
    risk_low_max DECIMAL(5,2) NOT NULL DEFAULT 1.5,
    risk_medium_max DECIMAL(5,2) NOT NULL DEFAULT 2.5,
    
    -- Thresholds de confiabilidade (response rate)
    reliability_acceptable_min DECIMAL(5,2) DEFAULT 30,
    reliability_good_min DECIMAL(5,2) DEFAULT 50,
    reliability_excellent_min DECIMAL(5,2) DEFAULT 70,
    
    -- Metadados
    description TEXT,
    changelog TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    activated_at TIMESTAMP,
    
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(partner_id, version)
);

CREATE INDEX idx_calculation_formulas_partner_id ON calculation_formulas(partner_id);
CREATE INDEX idx_calculation_formulas_active ON calculation_formulas(partner_id, active);

COMMENT ON TABLE calculation_formulas IS 'Fórmulas de cálculo versionadas por partner - permite evolução da metodologia sem perder histórico';
COMMENT ON COLUMN calculation_formulas.version IS 'Versão da fórmula (ex: 1.0, 2.0)';
COMMENT ON COLUMN calculation_formulas.active IS 'Apenas uma versão pode estar ativa por partner';

-- ============================================
-- ANALYTICS REPORTS (Snapshots)
-- ============================================
CREATE TABLE analytics_reports (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    report_data TEXT, -- JSON com DepartmentAnalytics + CalculationMetadata
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (questionnaire_id) REFERENCES assessment_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_analytics_reports_partner_id ON analytics_reports(partner_id);
CREATE INDEX idx_analytics_reports_department_id ON analytics_reports(department_id);
CREATE INDEX idx_analytics_reports_questionnaire_id ON analytics_reports(questionnaire_id);
CREATE INDEX idx_analytics_reports_created_at ON analytics_reports(created_at);

-- ============================================
-- ACTION PLAN TEMPLATES (Templates para auto-geração)
-- ============================================
CREATE TABLE action_plan_templates (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    category VARCHAR(100) NOT NULL,
    min_risk_level VARCHAR(20) NOT NULL CHECK (min_risk_level IN ('low', 'medium', 'high')),
    title_template TEXT NOT NULL,
    description_template TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    default_due_days INTEGER NOT NULL DEFAULT 30,
    auto_create BOOLEAN DEFAULT true,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_action_plan_templates_partner ON action_plan_templates(partner_id);
CREATE INDEX idx_action_plan_templates_category ON action_plan_templates(category);
CREATE INDEX idx_action_plan_templates_active ON action_plan_templates(active, auto_create);

COMMENT ON TABLE action_plan_templates IS 'Templates para auto-geração de planos de ação baseados em categorias de risco';
COMMENT ON COLUMN action_plan_templates.category IS 'Categoria NR-1: Sobrecarga, Autonomia, Relacionamento, Assédio, Reconhecimento, Jornada';
COMMENT ON COLUMN action_plan_templates.title_template IS 'Template do título com variáveis: {category}, {department_name}, {average_score}';

-- ============================================
-- RISK CATEGORIES (Risco por categoria)
-- ============================================
CREATE TABLE risk_categories (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    snapshot_id BIGINT,
    category VARCHAR(100) NOT NULL,
    average_score DECIMAL(5,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high')),
    question_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE
);

CREATE INDEX idx_risk_categories_partner ON risk_categories(partner_id);
CREATE INDEX idx_risk_categories_snapshot ON risk_categories(snapshot_id);
CREATE INDEX idx_risk_categories_category ON risk_categories(category);
CREATE INDEX idx_risk_categories_risk_level ON risk_categories(risk_level);

COMMENT ON TABLE risk_categories IS 'Armazena métricas de risco por categoria (ex: Sobrecarga, Assédio)';

-- ============================================
-- ACTION PLANS (Planos de ação)
-- ============================================
CREATE TABLE action_plans (
    id SERIAL PRIMARY KEY,
    partner_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    questionnaire_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    snapshot_id BIGINT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    risk_level VARCHAR(20) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high')),
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    category VARCHAR(100),
    responsible_name VARCHAR(255) NOT NULL,
    responsible_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    evidence_urls JSONB DEFAULT '[]',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_action_plans_partner ON action_plans(partner_id);
CREATE INDEX idx_action_plans_company ON action_plans(company_id);
CREATE INDEX idx_action_plans_department ON action_plans(department_id);
CREATE INDEX idx_action_plans_snapshot ON action_plans(snapshot_id);
CREATE INDEX idx_action_plans_status ON action_plans(status);
CREATE INDEX idx_action_plans_priority ON action_plans(priority);
CREATE INDEX idx_action_plans_category ON action_plans(category);
CREATE INDEX idx_action_plans_responsible ON action_plans(responsible_id);
CREATE INDEX idx_action_plans_due_date ON action_plans(due_date);

COMMENT ON TABLE action_plans IS 'Planos de ação para mitigar riscos identificados';
COMMENT ON COLUMN action_plans.evidence_urls IS 'URLs de evidências de execução (fotos, documentos, etc)';

-- ============================================
-- SEED DATA (Opcional - Partner de Teste)
-- ============================================
-- INSERT INTO partners (name, cnpj, email, mobile) 
-- VALUES ('RH Solutions Consultoria', '12.345.678/0001-90', 'contato@rhsolutions.com.br', '11987654321');
