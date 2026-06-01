CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA pgtasks;

CREATE TYPE pgtasks.task_priority AS ENUM ('low', 'medium', 'high');

CREATE TYPE pgtasks.task_status AS ENUM ('pending', 'completed');

CREATE TYPE pgtasks.task_icon AS ENUM (
    'mark', 'home',
    'job', 'supermarket',
    'cafe', 'activity',
    'drive', 'flight',
    'star', 'flag',
    'hospital', 'outdoor'
);

CREATE TABLE pgtasks.task (
    id UUID NOT NULL,
    owner_id UUID NOT NULL,
    group_id UUID,
    assignees UUID[],
    title TEXT NOT NULL CHECK (LENGTH(title) BETWEEN 1 AND 200),
    description TEXT NOT NULL CHECK (LENGTH(description) BETWEEN 1 AND 10000),
    location GEOGRAPHY(Point, 4326),
    is_favorite BOOLEAN NOT NULL DEFAULT FALSE,
    priority pgtasks.task_priority NOT NULL DEFAULT 'low',
    status pgtasks.task_status NOT NULL DEFAULT 'pending',
    icon pgtasks.task_icon NOT NULL DEFAULT 'mark',
    deadline TIMESTAMPTZ,
    notify_at TIMESTAMPTZ[],
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    purge_at TIMESTAMPTZ,
    title_lang regconfig NOT NULL DEFAULT 'simple',
    description_lang regconfig NOT NULL DEFAULT 'simple',
    search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector(title_lang, title), 'A') ||
        setweight(to_tsvector(description_lang, description), 'B')
        ) STORED,
    PRIMARY KEY (id, owner_id)
);
