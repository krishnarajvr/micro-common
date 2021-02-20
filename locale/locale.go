package locale

import (
	"github.com/kataras/i18n"
)

//Locale - Locale
type Locale struct {
	lang string
	i18n *i18n.I18n
}

//Get the language message
func (l *Locale) Get(key string, args ...interface{}) string {
	return l.i18n.Tr(l.lang, key, args...)
}

//SetLang - Set the language key
func (l *Locale) SetLang(langKey string) {
	l.lang = langKey
}

//New - Initilize the language
func (l *Locale) New(langKey string, langPath string, languages ...string) *Locale {
	I18n, err := i18n.New(i18n.Glob(langPath), languages...)

	if err != nil {
		panic(err)
	}

	l.lang = langKey
	l.i18n = I18n

	return l
}
