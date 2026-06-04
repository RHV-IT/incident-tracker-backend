CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'reporter',
    department VARCHAR(100) NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE issues (
    id SERIAL PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL,
    position VARCHAR(100) NOT NULL,          
    contact_info VARCHAR(255) NOT NULL,
    date_of_incident DATE NOT NULL,
    time_of_incident TIME NOT NULL,
    location_of_incident VARCHAR(255) NOT NULL,
    type_of_incident VARCHAR(150) NOT NULL,
    people_involved TEXT NOT NULL,
    description_of_incident TEXT NOT NULL,   
    immediate_action_taken TEXT NOT NULL,    
    injury_or_damage TEXT NOT NULL,          
    severity_level VARCHAR(50) NOT NULL,     
    supervisor_notified VARCHAR(255) NOT NULL,
    recommended_preventive_action TEXT NOT NULL 
);