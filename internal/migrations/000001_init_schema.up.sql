CREATE TABLE users(
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE events(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE shows(
    id UUID PRIMARY KEY,
    name TEXT,
    event_id UUID NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_event FOREIGN KEY(event_id)
        REFERENCES events(id)
        ON DELETE CASCADE
);

CREATE TABLE seats(
    ID UUID PRIMARY KEY,
    show_id UUID NOT NULL,
    seat_number TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'AVAILABLE',
    CONSTRAINT fk_key FOREIGN KEY(show_id)
        REFERENCES shows(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_seat_per_show UNIQUE(show_id, seat_number)
);