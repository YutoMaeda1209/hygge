CREATE TABLE accounts (
    id bigserial,
    discord_id text,
    email_address text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT accounts_pkey PRIMARY KEY (id),
    CONSTRAINT uni_accounts_discord_id UNIQUE (discord_id)
);

CREATE TABLE sessions (
    token_hash text NOT NULL,
    account_id bigint NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz,
    CONSTRAINT sessions_pkey PRIMARY KEY (token_hash),
    CONSTRAINT fk_sessions_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE INDEX idx_sessions_account_id ON sessions (account_id);

CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

CREATE TABLE subscribe_types (
    id bigserial,
    label text NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT subscribe_types_pkey PRIMARY KEY (id),
    CONSTRAINT uni_subscribe_types_label UNIQUE (label)
);

CREATE TABLE functions (
    id bigserial,
    label text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT functions_pkey PRIMARY KEY (id)
);

CREATE TABLE grant_func_perms (
    subscribe_type_id bigint NOT NULL,
    function_id bigint NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT grant_func_perms_pkey PRIMARY KEY (subscribe_type_id, function_id),
    CONSTRAINT fk_grant_func_perms_subscribe_type FOREIGN KEY (subscribe_type_id) REFERENCES subscribe_types (id),
    CONSTRAINT fk_grant_func_perms_function FOREIGN KEY (function_id) REFERENCES functions (id)
);

CREATE TABLE servers (
    id bigserial,
    guild_id text NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT servers_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_servers_guild_id ON servers (guild_id);

CREATE TABLE subscriptions (
    server_id bigint NOT NULL,
    subscribe_type_id bigint NOT NULL,
    contractor_id bigint NOT NULL,
    started_at timestamptz NOT NULL,
    expires_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT subscriptions_pkey PRIMARY KEY (server_id),
    CONSTRAINT fk_subscriptions_server FOREIGN KEY (server_id) REFERENCES servers (id),
    CONSTRAINT fk_subscriptions_subscribe_type FOREIGN KEY (subscribe_type_id) REFERENCES subscribe_types (id),
    CONSTRAINT fk_subscriptions_contractor FOREIGN KEY (contractor_id) REFERENCES accounts (id)
);

CREATE TABLE server_managers (
    server_id bigint NOT NULL,
    account_id bigint NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT server_managers_pkey PRIMARY KEY (server_id, account_id),
    CONSTRAINT fk_server_managers_server FOREIGN KEY (server_id) REFERENCES servers (id),
    CONSTRAINT fk_server_managers_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE TABLE discord_users (
    id bigserial,
    discord_id text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT discord_users_pkey PRIMARY KEY (id),
    CONSTRAINT uni_discord_users_discord_id UNIQUE (discord_id)
);

CREATE TABLE user_join_servers (
    discord_user_id bigint NOT NULL,
    server_id bigint NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT user_join_servers_pkey PRIMARY KEY (discord_user_id, server_id),
    CONSTRAINT fk_user_join_servers_discord_user FOREIGN KEY (discord_user_id) REFERENCES discord_users (id),
    CONSTRAINT fk_user_join_servers_server FOREIGN KEY (server_id) REFERENCES servers (id)
);

CREATE TABLE contacts (
    id bigserial,
    issued_by bigint,
    server_id bigint,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT contacts_pkey PRIMARY KEY (id),
    CONSTRAINT fk_contacts_discord_user FOREIGN KEY (issued_by) REFERENCES discord_users (id),
    CONSTRAINT fk_contacts_server FOREIGN KEY (server_id) REFERENCES servers (id)
);

CREATE TABLE contact_convos (
    id bigserial,
    contact_id bigint,
    send_by bigint,
    content text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT contact_convos_pkey PRIMARY KEY (id),
    CONSTRAINT fk_contact_convos_contact FOREIGN KEY (contact_id) REFERENCES contacts (id),
    CONSTRAINT fk_contact_convos_discord_user FOREIGN KEY (send_by) REFERENCES discord_users (id)
);

CREATE TABLE reaction_roles (
    id bigserial,
    server_id bigint,
    message_id text,
    reaction_id text,
    role_id text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT reaction_roles_pkey PRIMARY KEY (id),
    CONSTRAINT fk_reaction_roles_server FOREIGN KEY (server_id) REFERENCES servers (id)
);

CREATE TABLE voice_channels (
    id bigserial,
    server_id bigint,
    channel_id text,
    new_channel_name text,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT voice_channels_pkey PRIMARY KEY (id),
    CONSTRAINT uni_voice_channels_channel_id UNIQUE (channel_id),
    CONSTRAINT fk_voice_channels_server FOREIGN KEY (server_id) REFERENCES servers (id)
);
