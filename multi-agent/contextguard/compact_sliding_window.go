package contextguard

import (
	"fmt"
	"log/slog"
	"sync"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
)

// slidingWindowStrategy implements turn-count-based compaction. When the
// number of new Content entries since the last compaction exceeds maxTurns,
// all but a small recent window (30% of maxTurns, minimum 3) are summarized
// and replaced with a single summary message.
type slidingWindowStrategy struct {
	registry              ModelRegistry
	llm                   model.LLM
	maxTurns              int
	maxCompactionAttempts int
	mu                    sync.Mutex
}

// newSlidingWindowStrategy creates a sliding window strategy for a single agent.
func newSlidingWindowStrategy(registry ModelRegistry, llm model.LLM, maxTurns int, maxCompactAttempts int) *slidingWindowStrategy {
	return &slidingWindowStrategy{
		registry:              registry,
		llm:                   llm,
		maxTurns:              maxTurns,
		maxCompactionAttempts: maxCompactAttempts,
	}
}

// Name return strategy name.
func (sw *slidingWindowStrategy) Name() string {
	return StrategySlidingWindow
}

func (rw *slidingWindowStrategy) Compact(ctx agent.CallbackContext, req *model.LLMRequest) error {
	existingSummary := loadSummary(ctx)
	indexContentAtlastCompact := loadContentsAtCompaction(ctx)

	totalContents := len(req.Contents)
	turnsSinceCompact := totalContents - indexContentAtlastCompact

	if turnsSinceCompact <= rw.maxTurns {
		if existingSummary != "" {
			injectSummary(req, existingSummary, indexContentAtlastCompact)
		}
	}
	slog.Info("%s [%s]: turn limit exceeded, summarizing",
		"agent", ctx.AgentName(),
		"session", ctx.SessionID(),
		"totalContents", totalContents,
		"indexContentAtLastCompact", indexContentAtlastCompact,
		"turnsSinceCompact", turnsSinceCompact,
		"maxTurns", rw.maxTurns,
	)

	// lock to avoid 2 goroutines compacting same session at same time.
	rw.mu.Lock()
	defer rw.mu.Unlock()

	contextWindow := rw.registry.ContextWindow(req.Model)
	buffer := computeBuffer(contextWindow)
	threshold := contextWindow - buffer

	userContent := ctx.UserContent()
	todos := loadTodos(ctx)
	recentKeep := max(3, 0.3*float64(rw.maxTurns))

	for attemp := range rw.maxCompactionAttempts {
		splitIdx := safeSplitIndex(req.Contents, len(req.Contents)-int(recentKeep))
		oldMessages := req.Contents[:splitIdx]
		recentMessages := req.Contents[splitIdx:]

		if len(oldMessages) == 0 {
			slog.Warn("%s [%s]: nothing to compact (split at 0), aborting",
				"agent", ctx.AgentName(),
				"attempt", attemp+1,
			)
			break
		}

		summary, err := summarize(ctx, rw.llm, oldMessages, existingSummary, todos, buffer)
		if err != nil {
			slog.Error("ContextGuard [sliding_window]: summarization FAILED",
				"agent", ctx.AgentName(),
				"session", ctx.SessionID(),
				"error", err,
			)
			return fmt.Errorf("summarization failed: %w", err)
		}

		existingSummary = summary
		tokenEstimate := estimateContentTokens(oldMessages)
		persistSummary(ctx, summary, tokenEstimate)
		persistContentAtCompaction(ctx, totalContents)

		replaceSummary(req, summary, recentMessages)
		injectContinuation(req, userContent)

		newTokens := estimateTokens(req)
	}
}
