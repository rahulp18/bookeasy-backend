CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    duration_minutes INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE shows(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue TEXT NOT NULL,
    event_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_event FOREIGN KEY(event_id)
        REFERENCES events(id)
        ON DELETE CASCADE
);

CREATE TABLE seats(
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seat_row TEXT NOT NULL,
    seat_number TEXT NOT NULL,
    UNIQUE(seat_row, seat_number)
);

CREATE TABLE show_seats(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'available',
    
    CONSTRAINT fk_show FOREIGN KEY(show_id)
        REFERENCES shows(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_seat FOREIGN KEY(seat_id)
        REFERENCES seats(id)
        ON DELETE CASCADE,
    CONSTRAINT seat_status_check CHECK (status IN('available','locked','booked')),
    UNIQUE(show_id, seat_id)
);
CREATE TABLE bookings(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    show_id UUID NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_show FOREIGN KEY(show_id)
        REFERENCES shows(id)
        ON DELETE CASCADE,
    CONSTRAINT booking_status_check CHECK (status IN('pending','confirmed','cancelled'))
);
CREATE TABLE booking_seats(
    booking_id UUID NOT NULL,
    show_seat_id UUID NOT NULL,
    PRIMARY KEY(booking_id, show_seat_id),
    CONSTRAINT fk_booking FOREIGN KEY(booking_id)
        REFERENCES bookings(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_show_seat FOREIGN KEY(show_seat_id)
        REFERENCES show_seats(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_show_seats_show_id ON show_seats(show_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_show_id ON bookings(show_id);
