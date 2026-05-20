package contextguard

import (
	"fmt"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

const summarizeSystemPrompt = `You are summarizing a conversation to preserve context for continuing later.

Critical: This summary will be the ONLY context available when the conversation resumes. Assume all previous messages will be lost. Be thorough.

Required sections:

## Current State

- What was being discussed or worked on (exact user request if applicable)
- Current progress and what has been completed
- What was being addressed right now (incomplete work or open thread)
- What remains to be done or answered (specific, not vague)

## Key Information

- Facts, data, and specific details mentioned (names, dates, numbers, URLs, identifiers)
- User preferences, instructions, and constraints stated during the conversation
- Definitions, terminology, or domain knowledge established
- Any external resources, references, or sources mentioned

## Context & Decisions

- Decisions made during the conversation and why
- Alternatives that were considered and discarded (and why)
- Assumptions made
- Important clarifications or corrections that occurred
- Any blockers, risks, or open questions identified

## Exact Next Steps

Be specific. Don't write "continue with the task" — write exactly what should happen next, with enough detail that someone reading only this summary can pick up without asking questions.

Tone: Write as if briefing a colleague taking over mid-conversation. Include everything they would need to continue without asking questions. Write in the same language as the conversation.

Length: A dynamic word limit will be appended to this prompt at runtime based on the model's buffer size. Within that limit, err on the side of too much detail rather than too little. Critical context is worth the tokens.`

// loadSummary loads the summary from the agent's state, return "" if not found.
func loadSummary(ctx agent.CallbackContext) string {
	key := stateKeyPrefixSummary + ctx.AgentName()
	value, err := ctx.State().Get(key)
	if err != nil {
		return ""
	}
	return value.(string)
}

func injectSummary(req *model.LLMRequest, existingSummary string, lastCompactIdx int) {
	if len(req.Contents) > 0 && req.Contents[0] != nil {
		first := req.Contents[0]
		if first.Parts[0] != nil && first.Parts[0].Text != "" {
			return
		}
	}

	contentSummary := &genai.Content{
		Parts: []*genai.Part{
			{Text: existingSummary},
		},
		Role: "user",
	}

	if lastCompactIdx > 0 && lastCompactIdx <= len(req.Contents) {
		req.Contents = append(req.Contents[:lastCompactIdx], req.Contents[lastCompactIdx])
	} else {
		req.Contents = append(req.Contents, contentSummary)
	}
}

// loadContentsAtCompaction reads the Content count recorded at the last
// sliding-window compaction. Returns 0 if no compaction has happened yet.
func loadContentsAtCompaction(ctx agent.CallbackContext) int {
	key := stateKeyPrefixContentsAtCompaction + ctx.AgentName()
	val, err := ctx.State().Get(key)
	if err != nil {
		return 0
	}
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	}
	return 0
}

// computeBuffer compute buffer size for sliding window strategy.
func computeBuffer(contextWindow int) int {
	if contextWindow >= largeContextWindowThreshold {
		return largeContextWindowBuffer
	}
	return int(float64(contextWindow) * smallContextWindowRatio)
}

func summarize(ctx agent.CallbackContext, llm model.LLM, oldMessages []*genai.Content, existingSummary string, todos []TodoItem, buffer int) (string, error) {
	maxOutputTokens := int32(float64(buffer) * 0.5)
	maxWords := int32(float64(maxOutputTokens) * 0.75)

	systemPrompt := summarizeSystemPrompt + fmt.Sprintf("\n\nKeep the summary under %d words", maxWords)
	userPrompt := buildSummarizePrompt(oldMessages, existingSummary, todos)

	req := &model.LLMRequest{
		Model: llm.Name(),
		Contents: []*genai.Content{
			{
				Role:  "user",
				Parts: []*genai.Part{{Text: userPrompt}},
			},
		},
		Config: &genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemPrompt}}},
			MaxOutputTokens:   maxOutputTokens,
		},
	}
	var resultSummrized string
	for resp, err := range llm.GenerateContent(ctx, req, false) {
		if err != nil {
			return "", fmt.Errorf("summarization LLM call failed: %w", err)
		}
		if resp != nil && resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part != nil && part.Text != "" {
					resultSummrized += part.Text
				}
			}
		}
	}
	if resultSummrized == "" {
		return buildFallbackSummary(oldMessages, existingSummary), nil
	}
	return resultSummrized, nil
}

// BuildSummrizePrompt
func buildSummarizePrompt(oldMessages []*genai.Content, existingSummary string, todos []TodoItem) string {
	var sb strings.Builder
	sb.WriteString("Provide a detailed summary of the following conversation\n\n")
	if existingSummary != "" {
		sb.WriteString("Previous summary of context\n")
		sb.WriteString(existingSummary)
		sb.WriteString("\nEnd of previous summary")
		sb.WriteString("Incorporate the previous summary into your new summary, updating any information that has changed, correcting any inaccuracies, and adding any new details from the intervening conversation.")
	}
	sb.WriteString("Conversation to summarize\n")
	for _, content := range oldMessages {
		if content == nil {
			continue
		}
		role := content.Role
		parts := content.Parts
		for _, part := range parts {
			if part == nil {
				continue
			}

			if part.Text != "" {
				sb.WriteString(fmt.Sprintf("%s: %s\n", role, part.Text))
			}

			if part.FunctionCall != nil {
				sb.WriteString(fmt.Sprintf("%s : [called tool] %s \n", role, part.FunctionCall.Name))
			}

			if part.FunctionResponse != nil {
				sb.WriteString(fmt.Sprintf("%s : [tool] %s returned a result \n", role, part.FunctionResponse.Name))
			}
		}
	}
	sb.WriteString("End of conversation")

	if len(todos) > 0 {
		sb.WriteString("Current todos list\n")
		for _, t := range todos {
			fmt.Fprint(&sb, "- [%s] %s \n", t.Status, t.Content)
		}

		sb.WriteString("End of todos list")
		sb.WriteString("Include these tasks and their statuses in your summary under a dedicated \"## Todo List\" section. ")
		sb.WriteString("Instruct the resuming assistant to restore them using the `todos` tool to continue tracking progress.\n")
	}
	return sb.String()
}

func buildFallbackSummary(oldMessages []*genai.Content, existingSummary string) string {
	var sb strings.Builder
	if existingSummary != "" {
		sb.WriteString(existingSummary)
		sb.WriteString("\n\n")
	}

	for _, content := range oldMessages {
		if content == nil {
			continue
		}
		role := content.Role
		if role == "" {
			role = "unknown"
		}
		for _, part := range content.Parts {
			if part != nil && part.Text != "" {
				sb.WriteString(fmt.Sprintf("%s: ", role))
				if len(part.Text) > 200 {
					sb.WriteString(part.Text[:200])
				} else {
					sb.WriteString(part.Text)
				}
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

// ===================================== TODOS HELPER FUNCTION ==========================================

// TotoItem present for single task
type TodoItem struct {
	Content    string `json:"content"`
	Status     string `json:"status"`
	ActiveForm string `json:"active_form,omitempty"`
}

func loadTodos(ctx agent.CallbackContext) []TodoItem {
	return nil
}

// safeSplitIndex computes the largest index <= splitPoint that lies on a turn boundary.
// If splitPoint would split a turn, it is adjusted to the end of the previous turn.
func safeSplitIndex(contents []*genai.Content, splitPoint int) int {
	if splitPoint <= 0 || splitPoint >= len(contents) {
		return splitPoint
	}

	orgIdx := splitPoint

	splitPoint = walkBackToPairBoundary(contents, splitPoint)
	if splitPoint <= 0 {
		splitPoint = walkForwardToPairBoundary(contents, orgIdx)
	}

	if splitPoint <= 0 {
		splitPoint = 1
	}
	if splitPoint >= len(contents) {
		splitPoint = len(contents) - 1
	}
	return splitPoint
}

func walkBackToPairBoundary(contents []*genai.Content, idx int) int {
	for idx > 0 {
		c := contents[idx]
		if c == nil {
			return idx
		}

		if c.Role == "model" && contentHasFunctionCall(c) {
			idx++
			continue
		}

		if c.Role == "user" && contentHasFunctionResponse(c) {
			idx++
			continue
		}
		break
	}
	return idx
}

func walkForwardToPairBoundary(contents []*genai.Content, idx int) int {
	for idx > 0 {
		c := contents[idx]
		if c == nil {
			return idx
		}

		if c.Role == "model" && contentHasFunctionCall(c) {
			idx--
			continue
		}

		if c.Role == "user" && contentHasFunctionResponse(c) {
			idx--
			continue
		}
		break
	}
	return idx
}

// contentHashFunctionCall return true if content has function call.
func contentHasFunctionCall(c *genai.Content) bool {
	for _, part := range c.Parts {
		if part != nil && part.FunctionCall != nil {
			return true
		}
	}
	return false
}

// contentHashFunctionResponse return true if content has function response.
func contentHasFunctionResponse(c *genai.Content) bool {
	for _, part := range c.Parts {
		if part != nil && part.FunctionResponse != nil {
			return true
		}
	}
	return false
}
