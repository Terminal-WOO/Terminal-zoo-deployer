package main

import (
	"time"

	"github.com/ClappFormOrg/AI-CO/go/pkg/log"
)

type Specification struct {
	Logger                 log.Logger    `ignored:"true"`
	HTTPListenAddress      string        `default:":8080" split_words:"true"`
	HTTPWriteTimeout       time.Duration `default:"120s" split_words:"true"`
	HTTPIdleTimeout        time.Duration `default:"120s" split_words:"true"`
	HTTPReadTimeout        time.Duration `default:"10s" split_words:"true"`
	TerminationGracePeriod time.Duration `default:"5s" split_words:"true"`
	TLSKeyFile             string        `default:"" split_words:"true"`
	TLSCertFile            string        `default:"" split_words:"true"`
}
