CREATE TABLE "page" (
    slug TEXT NOT NULL PRIMARY KEY,
    id TEXT NOT NULL,
    uuid TEXT NOT NULL,
    title TEXT NOT NULL,
    html TEXT,
    plaintext TEXT,
    feature_image TEXT,
    featured boolean NOT NULL DEFAULT false,
    status TEXT NOT NULL,
    visibility TEXT NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	updated_at timestamp with time zone NOT NULL DEFAULT now(),
	published_at timestamp with time zone,
    custom_excerpt TEXT,
    canonical_url TEXT,
    excerpt TEXT,
    reading_time INTEGER,
    og_image TEXT,
    og_title TEXT,
    og_description TEXT,
    twitter_image TEXT,
    twitter_title TEXT,
    twitter_description TEXT,
    meta_title TEXT,
    meta_description TEXT,
    feature_image_alt TEXT,
    feature_image_caption TEXT,
    primary_tag TEXT,
    url TEXT
);