package management

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/auth0.v4"
	"gopkg.in/auth0.v4/internal/testing/expect"
)

func TestConnection(t *testing.T) {

	c := &Connection{
		Name:     auth0.Stringf("Test-Connection-%d", time.Now().Unix()),
		Strategy: auth0.String("auth0"),
	}

	var err error

	t.Run("Create", func(t *testing.T) {
		err = m.Connection.Create(c)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := c.Options.(*ConnectionOptions); !ok {
			t.Errorf("unexpected options type %T", c.Options)
		}
		t.Logf("%v\n", c)
	})

	t.Run("Read", func(t *testing.T) {
		c, err = m.Connection.Read(auth0.StringValue(c.ID))
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v\n", c)
	})

	t.Run("List", func(t *testing.T) {
		cs, err := m.Connection.List()
		if err != nil {
			t.Error(err)
		}
		for _, c := range cs.Connections {
			var ok bool

			switch c.GetStrategy() {
			case ConnectionStrategyAuth0:
				_, ok = c.Options.(*ConnectionOptions)
			case ConnectionStrategyGoogleOAuth2:
				_, ok = c.Options.(*ConnectionOptionsGoogleOAuth2)
			case ConnectionStrategyFacebook:
				_, ok = c.Options.(*ConnectionOptionsFacebook)
			case ConnectionStrategyApple:
				_, ok = c.Options.(*ConnectionOptionsApple)
			case ConnectionStrategyLinkedin:
				_, ok = c.Options.(*ConnectionOptionsLinkedin)
			case ConnectionStrategyGitHub:
				_, ok = c.Options.(*ConnectionOptionsGitHub)
			case ConnectionStrategyWindowsLive:
				_, ok = c.Options.(*ConnectionOptionsWindowsLive)
			case ConnectionStrategySalesforce, ConnectionStrategySalesforceCommunity, ConnectionStrategySalesforceSandbox:
				_, ok = c.Options.(*ConnectionOptionsSalesforce)
			case ConnectionStrategyEmail:
				_, ok = c.Options.(*ConnectionOptionsEmail)
			case ConnectionStrategySMS:
				_, ok = c.Options.(*ConnectionOptionsSMS)
			case ConnectionStrategyOIDC:
				_, ok = c.Options.(*ConnectionOptionsOIDC)
			case ConnectionStrategyAD:
				_, ok = c.Options.(*ConnectionOptionsAD)
			case ConnectionStrategyAzureAD:
				_, ok = c.Options.(*ConnectionOptionsAzureAD)
			case ConnectionStrategySAML:
				_, ok = c.Options.(*ConnectionOptionsSAML)
			default:
				_, ok = c.Options.(map[string]interface{})
			}

			if !ok {
				t.Errorf("unexpected options type %T", c.Options)
			}

			t.Logf("%s %s %T\n", c.GetID(), c.GetName(), c.Options)
		}
	})

	t.Run("Update", func(t *testing.T) {

		id := auth0.StringValue(c.ID)

		c.ID = nil       // read-only
		c.Name = nil     // read-only
		c.Strategy = nil // read-only

		c.Options = &ConnectionOptions{

			BruteForceProtection: auth0.Bool(true),
			ImportMode:           auth0.Bool(false), // try some zero values
			DisableSignup:        auth0.Bool(true),
			RequiresUsername:     auth0.Bool(false),

			CustomScripts: map[string]interface{}{"get_user": "function( { return callback(null) }"},
			Configuration: map[string]interface{}{"foo": "bar"},
		}

		err = m.Connection.Update(id, c)
		if err != nil {
			t.Error(err)
		}

		t.Logf("%v\n", c)
	})

	t.Run("Delete", func(t *testing.T) {
		err = m.Connection.Delete(auth0.StringValue(c.ID))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("ReadByName", func(t *testing.T) {
		cs, err := m.Connection.ReadByName("Username-Password-Authentication")
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v\n", cs)
	})

	t.Run("GoogleOAuth2", func(t *testing.T) {
		g := &Connection{
			Name:     auth0.Stringf("Test-Connection-%d", time.Now().Unix()),
			Strategy: auth0.String("google-oauth2"),
			Options: &ConnectionOptionsGoogleOAuth2{
				AllowedAudiences: []interface{}{
					"example.com",
					"api.example.com",
				},
				Profile:  auth0.Bool(true),
				Calendar: auth0.Bool(true),
				Youtube:  auth0.Bool(false),
			},
		}

		defer m.Connection.Delete(g.GetID())

		err := m.Connection.Create(g)
		if err != nil {
			t.Fatal(err)
		}

		o, ok := g.Options.(*ConnectionOptionsGoogleOAuth2)
		if !ok {
			t.Fatalf("unexpected type %T", o)
		}

		expect.Expect(t, o.GetProfile(), true)
		expect.Expect(t, o.GetCalendar(), true)
		expect.Expect(t, o.GetYoutube(), false)
		expect.Expect(t, o.Scopes(), []string{"email", "profile", "calendar"})

		t.Logf("%s\n", g)
	})

	t.Run("OIDC", func(t *testing.T) {
		o := &ConnectionOptionsOIDC{}
		expect.Expect(t, len(o.Scopes()), 0)

		o.SetScopes(true, "foo", "bar", "baz")
		expect.Expect(t, len(o.Scopes()), 3)
		expect.Expect(t, o.Scopes(), []string{"bar", "baz", "foo"})

		o.SetScopes(false, "baz")
		expect.Expect(t, len(o.Scopes()), 2)
		expect.Expect(t, o.Scopes(), []string{"bar", "foo"})
	})

	t.Run("Email", func(t *testing.T) {
		name := fmt.Sprintf("Test-Connection-Email-%d", time.Now().Unix())
		from := "{{application.name}} <test@example.com>"
		subject := "Email Login - {{application.name}}"
		syntax := "liquid"
		body := "<html><body>email contents</body></html>"
		scope := "openid profile"
		e := &Connection{
			Name:     auth0.String(name),
			Strategy: auth0.String("email"),
			Options: &ConnectionOptionsEmail{
				Email: &ConnectionOptionsEmailSettings{
					Syntax:  auth0.String(syntax),
					From:    auth0.String(from),
					Subject: auth0.String(subject),
					Body:    auth0.String(body),
				},
				OTP: &ConnectionOptionsOTP{
					TimeStep: auth0.Int(100),
					Length:   auth0.Int(4),
				},
				AuthParams: map[string]string{
					"scope": scope,
				},
				BruteForceProtection: auth0.Bool(true),
				DisableSignup:        auth0.Bool(true),
				Name:                 auth0.String(name),
			},
		}

		defer m.Connection.Delete(e.GetID())

		err := m.Connection.Create(e)
		if err != nil {
			t.Fatal(err)
		}

		o, ok := e.Options.(*ConnectionOptionsEmail)
		if !ok {
			t.Fatalf("unexpected type %T", o)
		}

		expect.Expect(t, o.GetEmail().GetSyntax(), syntax)
		expect.Expect(t, o.GetEmail().GetFrom(), from)
		expect.Expect(t, o.GetEmail().GetSubject(), subject)
		expect.Expect(t, o.GetEmail().GetBody(), body)
		expect.Expect(t, o.GetOTP().GetTimeStep(), 100)
		expect.Expect(t, o.GetOTP().GetLength(), 4)
		expect.Expect(t, o.AuthParams["scope"], scope)
		expect.Expect(t, o.GetBruteForceProtection(), true)
		expect.Expect(t, o.GetDisableSignup(), true)
		expect.Expect(t, o.GetName(), name)

		t.Logf("%s\n", e)
	})

	t.Run("SMS", func(t *testing.T) {
		name := fmt.Sprintf("Test-Connection-SMS-%d", time.Now().Unix())
		from := "+17777777777"
		template := "Your verification code is { code }}"
		syntax := "liquid"
		scope := "openid profile"
		twilioSid := "abc132asdfasdf56"
		twilioToken := "234127asdfsada23"
		messagingServiceSID := "273248090982390423"
		g := &Connection{
			Name:     auth0.String(name),
			Strategy: auth0.String("sms"),
			Options: &ConnectionOptionsSMS{
				From:     auth0.String(from),
				Template: auth0.String(template),
				Syntax:   auth0.String(syntax),
				OTP: &ConnectionOptionsOTP{
					TimeStep: auth0.Int(110),
					Length:   auth0.Int(5),
				},
				AuthParams: map[string]string{
					"scope": scope,
				},
				BruteForceProtection: auth0.Bool(true),
				DisableSignup:        auth0.Bool(true),
				Name:                 auth0.String(name),
				TwilioSID:            auth0.String(twilioSid),
				TwilioToken:          auth0.String(twilioToken),
				MessagingServiceSID:  auth0.String(messagingServiceSID),
			},
		}

		defer m.Connection.Delete(g.GetID())

		err := m.Connection.Create(g)
		if err != nil {
			t.Fatal(err)
		}

		o, ok := g.Options.(*ConnectionOptionsSMS)
		if !ok {
			t.Fatalf("unexpected type %T", o)
		}

		expect.Expect(t, o.GetTemplate(), template)
		expect.Expect(t, o.GetFrom(), from)
		expect.Expect(t, o.GetSyntax(), syntax)
		expect.Expect(t, o.GetOTP().GetTimeStep(), 110)
		expect.Expect(t, o.GetOTP().GetLength(), 5)
		expect.Expect(t, o.AuthParams["scope"], scope)
		expect.Expect(t, o.GetBruteForceProtection(), true)
		expect.Expect(t, o.GetDisableSignup(), true)
		expect.Expect(t, o.GetName(), name)
		expect.Expect(t, g.GetName(), name)
		expect.Expect(t, o.GetTwilioSID(), twilioSid)
		expect.Expect(t, o.GetTwilioToken(), twilioToken)
		expect.Expect(t, o.GetMessagingServiceSID(), messagingServiceSID)

		t.Logf("%s\n", g)
	})

	t.Run("SAML", func(t *testing.T) {
		g := &Connection{
			Name:     auth0.Stringf("Test-SAML-Connection-%d", time.Now().Unix()),
			Strategy: auth0.String("samlp"),
			Options: &ConnectionOptionsSAML{
				SignInEndpoint: auth0.String("https://saml.identity/provider"),
				// Sample certificate from https://golang.org/src/crypto/x509/example_test.go
				SigningCert: auth0.String(`-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRl
cm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAw
WjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
TW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFp
bC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q
5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC
7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAa
BgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUF
BwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcy
LmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2Nz
cDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAf
BgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisG
AQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29t
L0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+
gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283
TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq
0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoW
RvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOh
yE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==
-----END CERTIFICATE-----`),
				TenantDomain: auth0.String("example.con"),
			},
		}
		defer m.Connection.Delete(g.GetID())

		err := m.Connection.Create(g)
		if err != nil {
			t.Fatal(err)
		}

		o, ok := g.Options.(*ConnectionOptionsSAML)
		if !ok {
			t.Fatalf("unexpected type %T", o)
		}

		expect.Expect(t, o.GetSignInEndpoint(), "https://saml.identity/provider")
		expect.Expect(t, o.GetTenantDomain(), "example.con")
	})
}
