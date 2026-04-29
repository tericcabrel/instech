-- +goose Up
CREATE TABLE tools (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    category TEXT NOT NULL, -- language, framework, library
    sub_type TEXT, -- backend, frontend, fullstack, etc.
    prolang TEXT, -- python, javascript, etc.
    release_year INTEGER,
    dev_status TEXT, -- active, deprecated, etc.
    details TEXT,
    use_cases TEXT[], -- json array of strings ["UI", "SPA", "SSR", "Fullstack", "SEO", "API", "Backend"]
    tags TEXT[], -- json array of strings [component-based, declarative, functional, object-oriented].
    website TEXT,
    github TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_tools_slug ON tools (slug);
CREATE INDEX idx_tools_category ON tools (category);
CREATE INDEX idx_tools_sub_type ON tools (sub_type);
CREATE INDEX idx_tools_language ON tools (prolang);

CREATE TABLE relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_tool_id INTEGER NOT NULL,
    to_tool_id INTEGER NOT NULL,
    kind TEXT NOT NULL, --built_on, inspired_by, alternative_to, replaced_by, used_with
    metadata TEXT, -- json object containing the reason for the relationship which is the explanation
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_relationships_from_tool_id ON relationships (from_tool_id);
CREATE INDEX idx_relationships_to_tool_id ON relationships (to_tool_id);
CREATE INDEX idx_relationships_kind ON relationships (kind);

-- +goose Down
DROP TABLE tools;
DROP TABLE relationships;