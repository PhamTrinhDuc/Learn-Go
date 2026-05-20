package contextguard

var (
	PackageName                        = "Contextguard"
	StrategySlidingWindow              = "sliding_window"
	stateKeyPrefixSummary              = "__context_guard_summary_"
	stateKeyPrefixContentsAtCompaction = "__context_guard_contents_at_compaction_"

	largeContextWindowThreshold = 200_000
	largeContextWindowBuffer    = 20_000
	smallContextWindowRatio     = 0.20
)
