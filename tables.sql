CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'reporter',
    department VARCHAR(100) NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE incidents (
    id SERIAL PRIMARY KEY,
    
    -- Principal Person Involved (Who it happened to)
    principal_name VARCHAR(255) NOT NULL,
    principal_gender VARCHAR(50) NOT NULL,
    principal_dob VARCHAR(50) NOT NULL,
    principal_type VARCHAR(100) NOT NULL, -- patient, staff, visiting consultant, other
    patient_id VARCHAR(100),
    patient_ward_dept VARCHAR(255),
    staff_job_title VARCHAR(255),
    staff_phone VARCHAR(50),
    staff_place_of_work VARCHAR(255),
    staff_site VARCHAR(255),
    people_involved TEXT NOT NULL,

    -- When and Where The Incident Occurred
    date_of_incident VARCHAR(50) NOT NULL,
    time_of_incident VARCHAR(50) NOT NULL,
    location_of_incident VARCHAR(255) NOT NULL,
    incident_ward_dept VARCHAR(255) NOT NULL,
    
    -- Witnesses
    witnesses TEXT,
    witness_type VARCHAR(100),
    witness_ward_dept VARCHAR(255),
    witness_job_title VARCHAR(255),
    witness_phone VARCHAR(50), -- Maps to your json:"witenssPhone" typo safely

    -- Factual Description of the Incident
    is_near_miss BOOLEAN NOT NULL DEFAULT FALSE,
    cause_group VARCHAR(255) NOT NULL,
    causes TEXT NOT NULL,
    prescribing_doctor VARCHAR(255),

    -- Treatment Received
    treatment_received VARCHAR(255) NOT NULL,

    -- Equipment Involved
    equipment_involved VARCHAR(100) NOT NULL, -- Keeps string alignment with Go struct
    equipment_model VARCHAR(255),
    equipment_sent_for_repair BOOLEAN NOT NULL DEFAULT FALSE,
    equipment_withdrawn BOOLEAN NOT NULL DEFAULT FALSE,
    equipment_retained BOOLEAN NOT NULL DEFAULT FALSE,
    equipment_number VARCHAR(100),
    is_medical_device VARCHAR(50),            -- Keeps string alignment with Go struct
    
    -- Reporter Details (Section G)
    reporter_name VARCHAR(255) NOT NULL,
    reporter_designation VARCHAR(255) NOT NULL,
    signature BOOLEAN NOT NULL DEFAULT FALSE,
    reporter_info VARCHAR(255) NOT NULL,
    reporter_date VARCHAR(50) NOT NULL,       -- Avoids SQL 'date' keyword conflicts

    -- severity level
    severity_level VARCHAR(50) NOT NULL DEFAULT 'near miss',
    incident_status VARCHAR(50) NOT NULL DEFAULT 'unresolved'
);

create TABLE incident_management (
    id SERIAL PRIMARY KEY,
    incident_id INT UNIQUE NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,

    impact_on_service TEXT NOT NULL,
    contributory_factors TEXT NOT NULL,
    actions_taken_outcomes TEXT NOT NULL,
    recommendations TEXT NOT NULL,
    lessons_learned TEXT NOT NULL,

    informed_patient BOOLEAN NOT NULL DEFAULT FALSE,
    informed_relative BOOLEAN NOT NULL DEFAULT FALSE,
    informed_senior_manager BOOLEAN NOT NULL DEFAULT FALSE,
    informed_pharmacist BOOLEAN NOT NULL DEFAULT FALSE,
    police_incident_number VARCHAR(100),
    informed_other TEXT,

    risk_severity INT NOT NULL,
    risk_likelihood INT NOT NULL,
    risk_rating INT NOT NULL,

    ohs_absence_over_3_days BOOLEAN,
    ohs_act_of_violence_or_danger BOOLEAN,
    ohs_hospitalization_over_24_hours BOOLEAN,
    ohs_staff_name VARCHAR(255),
    ohs_staff_dob VARCHAR(50),
    ohs_staff_address TEXT,

    manager_name VARCHAR(255) NOT NULL,
    manager_signature BOOLEAN NOT NULL DEFAULT FALSE, -- Aligns with your signature standard
    manager_designation VARCHAR(255) NOT NULL,
    manager_date VARCHAR(50) NOT NULL
);

create TABLE incident_logs (
    id SERIAL PRIMARY KEY,
    incident_id INT REFERENCES incidents(id) ON DELETE CASCADE,
    changed_by INT REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50),
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE comments (
    id serial PRIMARY KEY,
    incident_id INT REFERENCES incidents(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    comment TEXT
);

CREATE TABLE death_reports (
  id SERIAL PRIMARY KEY,
  ref VARCHAR(100),
  reported_date VARCHAR(50),
  opened_date VARCHAR(50),
  submitted_time VARCHAR(50),
  handler VARCHAR(255),
  manager VARCHAR(255),
  location VARCHAR(255),
  department VARCHAR(255),
  specialty VARCHAR(255),
  exact_location VARCHAR(255),
  coding VARCHAR(100),
  type VARCHAR(100) DEFAULT 'Clinincal Incident',
  category VARCHAR(255),
  sub_category VARCHAR(100),
  risk_grading VARCHAR(100),
  result VARCHAR(255),
  actual_harm VARCHAR(255),
  potential_harm VARCHAR(255),
  details TEXT,
  incident_date VARCHAR(50),
  incident_time VARCHAR(50),
  description TEXT,
  action_taken TEXT,
  patient_involved BOOLEAN NOT NULL DEFAULT FALSE,
  patient_told BOOLEAN NOT NULL DEFAULT FALSE,
  family_told BOOLEAN NOT NULL DEFAULT FALSE,
  what_family_told TEXT,
  incident_investigation TEXT,
  -- departmental review meeting
  review_meeting_date VARCHAR(50),
  quality_assurance_lead VARCHAR(255),
  doctor_notified BOOLEAN NOT NULL DEFAULT FALSE,
  meeting_discussion_points TEXT,
  meeting_action_points TEXT,
  level_of_investigation VARCHAR(100)
);

-- Seed Initial Super Admin
INSERT INTO users (name, email, password, role, department) 
VALUES ('super admin', 'admin@example.com', '$2a$10$UQgnunKYIsM.hTWtjYooG.SPNKBqywEbOKddh1tU4tJuDiqfcn5Dm', 'superadmin', 'it');

-- Index for Fast Dashboard Performance
CREATE INDEX IF NOT EXISTS idx_incidents_id_desc ON incidents (id DESC);

CREATE INDEX IF NOT EXISTS idx_incident_management_incident_id ON incident_management (incident_id);

CREATE INDEX IF NOT EXISTS idx_incident_logs_incident_id ON incident_logs (incident_id);

CREATE INDEX IF NOT EXISTS idx_comment ON comments (id);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_users_global_search_trgm
ON users
USING gin ((name || ' ' || email || ' ' || role || ' ' || department) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_death_reports_id_desc ON death_reports (id DESC);

CREATE INDEX IF NOT EXISTS idx_death_reports_ref ON death_reports (ref);

CREATE INDEX IF NOT EXISTS idx_death_reports_dept_subcat ON death_reports (department, sub_category);
