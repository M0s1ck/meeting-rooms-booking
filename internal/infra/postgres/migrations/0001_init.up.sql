CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TYPE user_role AS ENUM ('admin', 'user');
CREATE TYPE booking_status AS ENUM ('active', 'cancelled');

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email citext NOT NULL UNIQUE,
    password_hash text NULL,
    role user_role NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_not_blank CHECK (btrim(email::text) != ''),
    CONSTRAINT users_password_hash_not_blank CHECK (
        password_hash IS NULL OR btrim(password_hash) != ''
    )
);

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    description text NULL,
    capacity integer NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT rooms_name_not_blank CHECK (btrim(name) != ''),
    CONSTRAINT rooms_capacity_positive CHECK (
        capacity IS NULL OR capacity > 0
    )
);

CREATE TRIGGER trg_rooms_set_updated_at
BEFORE UPDATE ON rooms
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE room_schedules (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id uuid NOT NULL UNIQUE REFERENCES rooms(id) ON DELETE RESTRICT,
    days_of_week smallint[] NOT NULL,
    start_time time NOT NULL,
    end_time time NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT room_schedules_time_order CHECK (end_time > start_time)
);

CREATE TABLE slots (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id uuid NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    start_at timestamptz NOT NULL,
    end_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT slots_time_order CHECK (end_at > start_at),
    CONSTRAINT ex_slots_no_overlap EXCLUDE USING gist (
        room_id WITH =,
        tstzrange(start_at, end_at, '[)') WITH &&)
);

CREATE INDEX idx_slots_start_at
    ON slots (start_at);

CREATE UNIQUE INDEX ux_slots_room_start_at
    ON slots (room_id, start_at);

CREATE TABLE bookings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id uuid NOT NULL REFERENCES slots(id) ON DELETE RESTRICT,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status booking_status NOT NULL DEFAULT 'active',
    conference_link text NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    cancelled_at timestamptz NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT bookings_status_cancelled_at_consistency CHECK (
        (status = 'active' AND cancelled_at IS NULL)
        OR (status = 'cancelled' AND cancelled_at IS NOT NULL)
    ),
    CONSTRAINT bookings_conference_link_not_blank CHECK (
        conference_link IS NULL OR btrim(conference_link) != ''
    )
);

CREATE TRIGGER trg_bookings_set_updated_at
BEFORE UPDATE ON bookings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE UNIQUE INDEX ux_bookings_one_active_per_slot
    ON bookings (slot_id)
    WHERE status = 'active';

CREATE INDEX idx_bookings_user_created_at
    ON bookings (user_id, created_at DESC, id DESC);

CREATE INDEX idx_bookings_created_at
    ON bookings (created_at DESC, id DESC);

CREATE INDEX idx_bookings_slot_id
    ON bookings (slot_id);
