package settlement

import "errors"

var (
	errSettlementSourceNotReady   = errors.New("settlement source not ready")
	errSettlementTemporaryFailure = errors.New("temporary settlement failure")
)
