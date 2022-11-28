DROP TABLE IF EXISTS emails;

ALTER TABLE IF EXISTS "emails_notification"
DROP
CONSTRAINT IF EXISTS "notification_email_unique";
DROP TABLE IF EXISTS "emails_notification";
