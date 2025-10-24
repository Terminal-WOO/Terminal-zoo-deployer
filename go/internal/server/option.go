package server

import "github.com/ClappFormOrg/AI-CO/go/pkg/log"

type Option func(*Handler)

func WithLogger(logger log.Logger) Option {
	return func(h *Handler) {
		h.logger = log.NewComponentLogger(logger, Component)
	}
}

func WithSecretTLSKey(key []byte) Option {
	return func(h *Handler) {
		h.tlsKey = key
	}
}

func WithSecretTLSCert(crt []byte) Option {
	return func(h *Handler) {
		h.tlsCrt = crt
	}
}
