CREATE TYPE wallet_type AS ENUM ('regular', 'bank');
CREATE TYPE transaction_status AS ENUM ('pending', 'processing', 'confirmed', 'failed');

CREATE TABLE wallet (
    public_key TEXT PRIMARY KEY NOT NULL,
    private_key TEXT UNIQUE NOT NULL,
    address TEXT UNIQUE NOT NULL,
    seed_phrase TEXT NOT NULL,
    kind wallet_type NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE transaction (
    hash TEXT PRIMARY KEY NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    status transaction_status NOT NULL,
    confirmations INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (from_address) REFERENCES wallet(address),
    FOREIGN KEY (to_address) REFERENCES wallet(address)
);
