CREATE TABLE "tag" (
    tag TEXT NOT NULL PRIMARY KEY,
    summary TEXT
);

CREATE TABLE "article_tag" (
    slug TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY(slug,tag),
    constraint fk_blog
        FOREIGN KEY(slug)
            references article(slug),
    constraint fk_tag
        FOREIGN KEY(tag)
            references tag(tag)
);

CREATE TABLE "category" (
    category TEXT NOT NULL PRIMARY KEY,
    summary TEXT
);

ALTER TABLE "article" 
ADD COLUMN category TEXT NULL,
add  constraint fk_category
        FOREIGN KEY(category)
            references category(category);