DROP TYPE IF EXISTS notification_status CASCADE;

ALTER TABLE IF EXISTS "emails_notification"
DROP
CONSTRAINT IF EXISTS "notification_email_unique";
DROP TABLE IF EXISTS "emails_notification";

DROP TABLE IF EXISTS shares;
DROP TABLE IF EXISTS emails;
