CREATE TABLE Applications
(
    application_id SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Endpoints
(
    endpoint_id         SERIAL PRIMARY KEY,
    application_id      INT NOT NULL REFERENCES Applications (application_id) ON DELETE CASCADE,
    url                 TEXT NOT NULL,
    monitoring_interval INT DEFAULT 30,
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UptimeLogs
(
    log_id        SERIAL PRIMARY KEY,
    endpoint_id   INT NOT NULL REFERENCES Endpoints (endpoint_id) ON DELETE CASCADE,
    status_code   INT,
    response_time INT,
    is_up         BOOLEAN,
    timestamp     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE NotificationChannels
(
    channel_id     SERIAL PRIMARY KEY,
    application_id INT NOT NULL REFERENCES Applications (application_id) ON DELETE CASCADE,
    type           VARCHAR(50) NOT NULL,
    details        JSON,
    is_active      BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Alerts
(
    alert_id    SERIAL PRIMARY KEY,
    channel_id  INT NOT NULL REFERENCES NotificationChannels (channel_id) ON DELETE CASCADE,
    endpoint_id INT NOT NULL REFERENCES Endpoints (endpoint_id) ON DELETE CASCADE,
    message     TEXT,
    sent_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
