package servicev1

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/rs/xid"
	simple "github.com/steambap/captcha"

	apiv1 "wailik.com/internal/captcha/api/v1"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/util"
)

const (
	ModeSimple = "simple"
)

type captchaCached struct {
	nonce  string
	salt   int
	answer interface{}
}

type captchaSrvc struct {
	service *service
}

func newCaptchaSrvc(s *service) *captchaSrvc {
	return &captchaSrvc{
		service: s,
	}
}

func (s *captchaSrvc) GenerateCaptcha(mode apiv1.ChallengeMode) (*apiv1.Challenge, error) {
	nonce := xid.New().String()
	switch mode {
	case ModeSimple:
		captcha, _ := simple.New(150, 50)
		buffer := new(bytes.Buffer)
		encoder := base64.NewEncoder(base64.StdEncoding, buffer)
		err := captcha.WriteImage(encoder)
		if err != nil {
			return nil, err
		}
		data := apiv1.ChallengeData{
			"image": buffer.String(),
		}
		log.Debugf("%+v", captcha.Text)
		log.Debugf("%+v", nonce)
		salt := util.RandomInt(9999, 1000)

		r := s.service.cache.Set(nonce, captchaCached{
			nonce:  nonce,
			salt:   salt,
			answer: captcha.Text,
		}, time.Second*60)

		log.Debugf("%+v", r)

		challenge := apiv1.Challenge{
			Nonce: nonce,
			Data:  data,
			Salt:  salt,
		}

		ts := time.Now().Unix()
		log.Debugf("%+v, %+v", ts, hash(ts, salt, captcha.Text))

		return &challenge, nil
	}

	return nil, nil
}

func (s *captchaSrvc) VerifyCaptcha(captcha apiv1.Captcha) bool {
	ts, _ := strconv.ParseInt(captcha.Timestamp, 10, 64)

	if !checkTimestamp(ts, time.Second*60) {
		log.Debugf("timestamp expired")

		return false
	}

	log.Debugf("nonce:%+v", captcha.Nonce)
	r, ok := s.service.cache.Get(captcha.Nonce)
	if !ok {
		log.Debugf("no answer for the captcha:%+v", r)

		return false
	}

	cached := (r.(captchaCached))

	if !validateSignature(captcha.Signature, captcha.Mode, ts, cached.salt, captcha.Try) {
		log.Debugf("signature error")

		return false
	}

	if validateCaptcha(captcha.Mode, captcha.Try, cached.answer) {
		s.service.cache.Del(captcha.Nonce)

		return true
	}
	log.Debugf("unsuccessful try")

	return false
}

func validateSignature(signature string, mode apiv1.ChallengeMode, timestamp int64, salt int, try map[string]interface{}) bool {
	switch mode {
	case ModeSimple:
		r := hash(timestamp, salt, (try["text"]).(string))

		return r == signature
	}

	return false
}

func checkTimestamp(timestamp int64, duration time.Duration) bool {
	now := time.Now()
	then := time.Unix(timestamp, 0)

	log.Debugf("current:%+v", now.Sub(then))
	log.Debugf("setting:%+v", duration)

	return now.Sub(then) < duration
}

func validateCaptcha(mode apiv1.ChallengeMode, try map[string]interface{}, answer interface{}) bool {
	switch mode {
	case ModeSimple:
		text := (try["text"]).(string)

		return text == answer.(string)
	}

	return false
}

func hash(timestamp int64, salt int, text string) string {
	saltStr := strconv.Itoa(salt)
	timestampStr := strconv.FormatInt(timestamp, 10)
	data := timestampStr + saltStr + text
	m := sha256.New()
	m.Write([]byte(data))
	r := m.Sum([]byte(""))

	return base64.RawStdEncoding.EncodeToString(r)
}
