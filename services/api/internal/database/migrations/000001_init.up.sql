CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE fabric_piece_status AS ENUM ('in_stock', 'used_up', 'archived');
CREATE TYPE project_status AS ENUM ('planned', 'making', 'completed', 'archived');
CREATE TYPE stock_movement_type AS ENUM ('received', 'consumed', 'adjusted', 'restored');

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    external_provider text NOT NULL,
    external_id text NOT NULL,
    nickname text NOT NULL DEFAULT '',
    avatar_url text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (external_provider, external_id)
);

CREATE TABLE assets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    storage_key text NOT NULL,
    thumbnail_key text NOT NULL DEFAULT '',
    mime_type text NOT NULL,
    width integer NOT NULL CHECK (width > 0),
    height integer NOT NULL CHECK (height > 0),
    size_bytes bigint NOT NULL CHECK (size_bytes >= 0),
    sha256 text NOT NULL,
    perceptual_hash text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (owner_id, sha256)
);

CREATE TABLE fabrics (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    name text NOT NULL,
    material text NOT NULL DEFAULT '',
    pattern text NOT NULL DEFAULT '',
    style_tags text[] NOT NULL DEFAULT '{}',
    colors jsonb NOT NULL DEFAULT '[]',
    cover_asset_id uuid REFERENCES assets(id),
    memo text NOT NULL DEFAULT '',
    recognition_status text NOT NULL DEFAULT 'confirmed',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE fabric_pieces (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    fabric_id uuid NOT NULL REFERENCES fabrics(id),
    label text NOT NULL DEFAULT '整块布',
    length_mm integer CHECK (length_mm IS NULL OR length_mm >= 0),
    width_mm integer CHECK (width_mm IS NULL OR width_mm >= 0),
    quantity numeric(12,3) NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    unit text NOT NULL DEFAULT 'piece',
    storage_location text NOT NULL DEFAULT '',
    purchase_channel text NOT NULL DEFAULT '',
    purchase_date date,
    total_price_cent bigint CHECK (total_price_cent IS NULL OR total_price_cent >= 0),
    status fabric_piece_status NOT NULL DEFAULT 'in_stock',
    source_label text NOT NULL DEFAULT '',
    memo text NOT NULL DEFAULT '',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE projects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    name text NOT NULL,
    kind text NOT NULL DEFAULT '',
    status project_status NOT NULL DEFAULT 'planned',
    size text NOT NULL DEFAULT '',
    memo text NOT NULL DEFAULT '',
    reference_project_id uuid REFERENCES projects(id),
    version bigint NOT NULL DEFAULT 1,
    started_at date,
    completed_at date,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE project_materials (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    fabric_id uuid REFERENCES fabrics(id),
    fabric_piece_id uuid REFERENCES fabric_pieces(id),
    role text NOT NULL,
    planned_amount_mm integer CHECK (planned_amount_mm IS NULL OR planned_amount_mm >= 0),
    actual_amount_mm integer CHECK (actual_amount_mm IS NULL OR actual_amount_mm >= 0),
    is_suggestion boolean NOT NULL DEFAULT false,
    suggestion_snapshot jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (fabric_id IS NOT NULL OR suggestion_snapshot IS NOT NULL)
);

CREATE TABLE stock_movements (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid NOT NULL REFERENCES users(id),
    fabric_piece_id uuid NOT NULL REFERENCES fabric_pieces(id),
    project_id uuid REFERENCES projects(id),
    movement_type stock_movement_type NOT NULL,
    quantity_delta numeric(12,3) NOT NULL,
    quantity_after numeric(12,3) NOT NULL CHECK (quantity_after >= 0),
    note text NOT NULL DEFAULT '',
    idempotency_key text,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (owner_id, idempotency_key)
);

CREATE TABLE entity_assets (
    asset_id uuid NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    entity_type text NOT NULL,
    entity_id uuid NOT NULL,
    purpose text NOT NULL DEFAULT 'gallery',
    sort_order integer NOT NULL DEFAULT 0,
    PRIMARY KEY (asset_id, entity_type, entity_id)
);

CREATE TABLE external_mappings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    provider text NOT NULL,
    entity_type text NOT NULL,
    entity_id uuid NOT NULL,
    external_id text NOT NULL,
    external_parent_id text NOT NULL DEFAULT '',
    synced_at timestamptz,
    sync_error text NOT NULL DEFAULT '',
    UNIQUE (provider, entity_type, entity_id),
    UNIQUE (provider, external_id)
);

CREATE INDEX fabrics_owner_updated_idx ON fabrics (owner_id, updated_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX fabric_pieces_fabric_idx ON fabric_pieces (fabric_id) WHERE deleted_at IS NULL;
CREATE INDEX projects_owner_updated_idx ON projects (owner_id, updated_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX project_materials_project_idx ON project_materials (project_id);
CREATE INDEX stock_movements_piece_created_idx ON stock_movements (fabric_piece_id, created_at DESC);

