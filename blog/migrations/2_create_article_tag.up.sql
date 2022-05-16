CREATE TABLE "article_tag" (
    slug TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY(slug,tag),
    constraint fk_blog
        FOREIGN KEY(slug)
            references article(slug),
    constraint fk_tag
        FOREIGN KEY(tag)
            references tag(slug)
);
