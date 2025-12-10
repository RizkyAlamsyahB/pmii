CREATE TABLE "post_tags" (
  "post_id" int NOT NULL,
  "tag_id" int NOT NULL,
  PRIMARY KEY ("post_id", "tag_id")
);

ALTER TABLE "post_tags" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
ALTER TABLE "post_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id");
