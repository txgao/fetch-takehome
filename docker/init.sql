

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS item (
    item_uuid UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    price DECIMAL NOT NULL,
    short_description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS receipt (
    receipt_uuid UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    total DECIMAL NOT NULL,
    purchase_date DATE NOT NULL,
    purchase_time TIME,
    retailer VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS receipt_items (
    item_uuid UUID NOT NULL,
    receipt_uuid UUID NOT NULL
);

ALTER TABLE ONLY receipt_items
ADD CONSTRAINT receipt_items_pkey PRIMARY KEY (item_uuid, receipt_uuid);
ALTER TABLE ONLY receipt_items
ADD CONSTRAINT receipt_items_receipt_uuid_fkey FOREIGN KEY (receipt_uuid) REFERENCES receipt;
ALTER TABLE ONLY receipt_items
ADD CONSTRAINT receipt_items_item_uuid_fkey FOREIGN KEY (item_uuid) REFERENCES item;