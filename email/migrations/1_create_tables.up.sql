-- template contains email templates.
-- It is a separate table from message so that multiple messages can
-- easily be sent with the same contents.
CREATE TABLE "template" (
    -- id is the unique id for this template.
    id TEXT PRIMARY KEY,

    -- sender is the sender email to use.
    sender TEXT NOT NULL,

	-- subject is the subject line to use.
    subject TEXT NOT NULL,

    -- body_text and body_html are the plaintext and html bodies to use.
    body_text TEXT NOT NULL,
    body_html TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- user tracks users who have opted in for emails.
CREATE TABLE "user" (
    -- email_address is the email address of the person who opted in.
    email_address TEXT NOT NULL PRIMARY KEY,

    -- optin tracks the current opt-in status.
    optin BOOLEAN NOT NULL,

    -- optin_changed tracks the last time the user's opt-in status changed.
    optin_changed TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE "message" (
    -- id is a unique id for this particular email message.
    id BIGSERIAL PRIMARY KEY,
    -- email_address is the email address to send the message to.
    email_address TEXT NOT NULL REFERENCES "user" (email_address),

	-- template_id is a reference to the email template to send.
    template_id TEXT NOT NULL REFERENCES "template" (id),

    -- provider_id is the unique id the email provider generates for this particular message.
    -- It is non-null when the message has been successfully sent.
    provider_id TEXT NULL,

    -- scheduled_at is the time when the message is scheduled for delivery.
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- sent_at is the timestamp when the message was successfully sent.
    sent_at TIMESTAMP WITH TIME ZONE NULL
);

-- unsubscribe_event tracks when a user unsubscribes from future emails.
CREATE TABLE "unsubscribe_event" (
    -- email_address is the email address of the user who unsubscribed.
    email_address TEXT NOT NULL REFERENCES "user" (email_address),

    -- message_id is the message that led to the user unsubscribing.
    message_id BIGINT NOT NULL REFERENCES "message" (id),

	-- unsubscribed_at is when the user unsubscribed.
    unsubscribed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
