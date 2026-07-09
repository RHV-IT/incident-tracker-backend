@echo off
docker exec -it issuetracker_db psql -U tracker_user -d issuetracker -c "CREATE INDEX IF NOT EXISTS idx_incident_logs_incident_id ON incident_logs (incident_id);"