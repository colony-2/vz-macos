package vmapi

// Recover uses durable adapter records alongside core journals. Only known
// eligible undelivered bootstrap with identity, intent, native phase and trust
// evidence may resume. Never repeat ambiguous clone/run or arbitrary Exec.
// Crash after delivery before success recording may repeat an idempotent recipe:
// bootstrap is at least once. Reads never trigger recovery.
type RecoveryState string

const (
	RecoveryRecovered RecoveryState = "recovered"
	RecoverySkipped   RecoveryState = "skipped"
	RecoveryUncertain RecoveryState = "uncertain"
)

type RecoveryRecord struct {
	Instance          InstanceRef
	IntentFingerprint string
	State             RecoveryState
	Code              ErrorCode
}

// Records are bounded by Info. Overflow sets Complete false; no cursor is implied.
type RecoveryResult struct {
	Records  []RecoveryRecord
	Complete bool
}
