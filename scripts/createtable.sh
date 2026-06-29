docker exec -it issuetracker_db psql -U tracker_user -d issuetracker -c "
CREATE TABLE incident_logs (
    id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL,
    changed_by INT,
    action VARCHAR(50),
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);"