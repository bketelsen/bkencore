CREATE TABLE "article" (
    slug TEXT NOT NULL PRIMARY KEY,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	modified_at timestamp with time zone NOT NULL DEFAULT now(),
    published BOOLEAN,
    title TEXT NOT NULL,
    summary TEXT,
    body TEXT,
    body_rendered TEXT
);

CREATE TABLE "email_subscriber" (
    email TEXT NOT NULL PRIMARY KEY,
    subscribed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
