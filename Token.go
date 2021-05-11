package youtube

import (
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

func GetToken(clientID string, channelID string) (*oauth2.Token, *errortools.Error) {
	accessToken := "ya29.a0AfH6SMC756lCOskXeOMqXDh9o29PlEbm0_l7GDUgpfvwzv5SeoGjNeqVXMAS2T9jYaIhBNFcaOc32yjZ6zxpJp8p8vFjjGD3r4mh9ck_EjPc5q0_qPnBN5I2wiIiUgAUiTLZLEsuXGPP87A6OPoTffJb31V4"
	scope := "https://www.googleapis.com/auth/youtube.readonly https://www.googleapis.com/auth/yt-analytics.readonly"
	refreshToken := "1//09WbbhLZbARhBCgYIARAAGAkSNwF-L9Ir2pwNgbDfyBCYvFb5e5QvQnRwY-_jPJ9KHaqC64cSccO4oeOgbC4UWTwNjuSMygh3hZU"
	expiry, _ := time.Parse("2006-01-02 15:04:05.9999999", "2021-05-11 21:45:28.4942788")

	return &oauth2.Token{
		&accessToken,
		&scope,
		nil,
		nil,
		&refreshToken,
		&expiry,
	}, nil
}

func SaveToken(clientID string, channelID string, token *oauth2.Token) *errortools.Error {
	return nil

	/*

		if token == nil {
			return nil
		}

		sqlUpdate := "SET AccessToken = SOURCE.AccessToken, Expiry = SOURCE.Expiry"

		tokenType := "NULLIF('','')"
		if token.TokenType != nil {
			if *token.TokenType != "" {
				tokenType = fmt.Sprintf("'%s'", *token.TokenType)
				sqlUpdate = fmt.Sprintf("%s, TokenType = SOURCE.TokenType", sqlUpdate)
			}
		}

		accessToken := "NULLIF('','')"
		if token.AccessToken != nil {
			if *token.AccessToken != "" {
				accessToken = fmt.Sprintf("'%s'", *token.AccessToken)
			}
		}

		refreshToken := "NULLIF('','')"
		if token.RefreshToken != nil {
			if *token.RefreshToken != "" {
				refreshToken = fmt.Sprintf("'%s'", *token.RefreshToken)
				sqlUpdate = fmt.Sprintf("%s, RefreshToken = SOURCE.RefreshToken", sqlUpdate)
			}
		}

		expiry := "TIMESTAMP(NULL)"
		if token.Expiry != nil {
			expiry = fmt.Sprintf("TIMESTAMP('%s')", (*token.Expiry).Format("2006-01-02T15:04:05"))
		}

		scope := "NULLIF('','')"
		if token.Scope != nil {
			if *token.Scope != "" {
				scope = fmt.Sprintf("'%s'", *token.Scope)
				sqlUpdate = fmt.Sprintf("%s, Scope = SOURCE.Scope", sqlUpdate)
			}
		}

		sql := "MERGE `" + tableRefreshToken + "` AS TARGET " +
			"USING  (SELECT '" +
			apiName + "' AS Api,'" +
			clientID + "' AS ClientID," +
			tokenType + " AS TokenType," +
			accessToken + " AS AccessToken," +
			refreshToken + " AS RefreshToken," +
			expiry + " AS Expiry," +
			scope + " AS Scope) AS SOURCE " +
			" ON TARGET.Api = SOURCE.Api " +
			" AND TARGET.ClientID = SOURCE.ClientID " +
			"WHEN MATCHED THEN " +
			"	UPDATE " + sqlUpdate +
			" WHEN NOT MATCHED BY TARGET THEN " +
			"	INSERT (Api, ClientID, TokenType, AccessToken, RefreshToken, Expiry, Scope) " +
			"	VALUES (SOURCE.Api, SOURCE.ClientID, SOURCE.TokenType, SOURCE.AccessToken, SOURCE.RefreshToken, SOURCE.Expiry, SOURCE.Scope)"

		return service.Run(sql, "saving token")*/
}
