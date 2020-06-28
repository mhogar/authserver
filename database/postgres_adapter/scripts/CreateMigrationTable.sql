CREATE TABLE IF NOT EXISTS public.migration (
    "timestamp" varchar(14) NOT NULL,
    CONSTRAINT migration_pk PRIMARY KEY ("timestamp")
);