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
);

-- Seed Initial Super Admin
INSERT INTO users (name, email, password, role, department) 
VALUES ('super admin', 'admin@example.com', '$2a$10$UQgnunKYIsM.hTWtjYooG.SPNKBqywEbOKddh1tU4tJuDiqfcn5Dm', 'superadmin', 'it');

-- Index for Fast Dashboard Performance
CREATE INDEX IF NOT EXISTS idx_incidents_id_desc ON incidents (id DESC);