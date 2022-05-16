CREATE TABLE "page_tag" (
    slug TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY(slug,tag),
    constraint fk_page
        FOREIGN KEY(slug)
            references page(slug),
    constraint fk_tagpage
        FOREIGN KEY(tag)
            references tag(slug)
);
