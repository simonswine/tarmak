package wing

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type Wing struct {
	log   *logrus.Entry
	flags *viper.Viper
}

type Provider interface {
	ID() string
}

func New(flags *viper.Viper) *Wing {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	t := &Wing{
		log:   logger.WithField("app", "wing"),
		flags: flags,
	}
	return t
}

func (w *Wing) Must(err error) *Wing {
	if err != nil {
		w.log.Fatal(err)
	}
	return w
}
