ALTER TABLE voice_channels RENAME TO voice_lobby_channels;

ALTER SEQUENCE voice_channels_id_seq RENAME TO voice_lobby_channels_id_seq;

ALTER TABLE voice_lobby_channels RENAME CONSTRAINT voice_channels_pkey TO voice_lobby_channels_pkey;

ALTER TABLE voice_lobby_channels RENAME CONSTRAINT uni_voice_channels_channel_id TO uni_voice_lobby_channels_channel_id;

ALTER TABLE voice_lobby_channels RENAME CONSTRAINT fk_voice_channels_server TO fk_voice_lobby_channels_server;

-- Voice channels created from a lobby, tracked so they can be deleted once empty (also across restarts).
CREATE TABLE voice_lobby_rooms (
    id bigserial,
    voice_lobby_channel_id bigint NOT NULL,
    channel_id text NOT NULL,
    created_at timestamptz,
    updated_at timestamptz,
    CONSTRAINT voice_lobby_rooms_pkey PRIMARY KEY (id),
    CONSTRAINT uni_voice_lobby_rooms_channel_id UNIQUE (channel_id),
    CONSTRAINT fk_voice_lobby_rooms_voice_lobby_channel FOREIGN KEY (voice_lobby_channel_id) REFERENCES voice_lobby_channels (id)
);
