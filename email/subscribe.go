package email

import (
	"context"
	"encoding/base64"
	"errors"
	"log"

	"github.com/gorilla/securecookie"

	"encore.dev/storage/sqldb"
)

type SubscribeParams struct {
	Email string
}

// Subscribe subscribes to the email newsletter for a given email.
//encore:api public method=POST path=/email/subscribe
func Subscribe(ctx context.Context, p *SubscribeParams) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "user" (email_address, optin)
		VALUES ($1, true)
		ON CONFLICT (email_address) DO UPDATE
		SET optin = true, optin_changed = NOW()
	`, p.Email)
	return err
}

// getAllSubscribers returns all the subscribers to the email newsletter.
func getAllSubscribers(ctx context.Context) (emails []string, err error) {
	rows, err := sqldb.Query(ctx, `
		SELECT email_address
		FROM "user"
		WHERE optin
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, rows.Err()
}

type UnsubscribeParams struct {
	// Token is the unsubscribe token in to the email.
	Token string
}

// Unsubscribe unsubscribes the user from the email list.
//encore:api public method=POST path=/email/unsubscribe
func Unsubscribe(ctx context.Context, params *UnsubscribeParams) error {
	email, emailID, err := decodeUnsubscribeToken(params.Token)
	if err != nil {
		return err
	}
	_, err = sqldb.Exec(ctx, `
		UPDATE "user" SET optin = false, optin_changed = NOW()
		WHERE email_address = $1 AND optin
	`, email)
	if err != nil {
		return err
	}

	_, err = sqldb.Exec(ctx, `
		INSERT INTO "unsubscribe_event" (email_address, message_id)
		VALUES ($1, $2)
	`, email, emailID)
	return err
}

// isOptedIn reports whether a particular email address has opted in to emails.
func isOptedIn(ctx context.Context, email string) (bool, error) {
	var status bool
	err := sqldb.QueryRow(ctx, `
		SELECT optin FROM "user" WHERE email_address = $1
	`, email).Scan(&status)
	if errors.Is(err, sqldb.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return status, nil
}

// unsubscribeTokenData is the unencoded data that makes up the unsubscribe token.
type unsubscribeTokenData struct {
	Email     string
	MessageID int64
}

// encodeUnsubscribeToken encodes an email address and message id as a token using HMAC.
func encodeUnsubscribeToken(email string, messageID int64) (string, error) {
	return tokenCookie.Encode("token", unsubscribeTokenData{Email: email, MessageID: messageID})
}

// decodeUnsubscribeToken decodes a token into the email and message id that it contains.
func decodeUnsubscribeToken(token string) (email string, messageID int64, err error) {
	var data unsubscribeTokenData
	err = tokenCookie.Decode("token", token, &data)
	return data.Email, data.MessageID, nil
}

// tokenCookie is the securecookie for encoding email unsubscribe tokens.
var tokenCookie = func() *securecookie.SecureCookie {
	hashKey, err := base64.RawURLEncoding.DecodeString(secrets.TokenHashKey)
	if err != nil {
		log.Fatalln("bad TokenHashKey:", err)
	}
	return securecookie.New(hashKey, nil)
}()
