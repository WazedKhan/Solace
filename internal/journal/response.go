package journal

func toMoodResponse(m *Mood) *MoodResponse {
	if m == nil {
		return nil
	}
	return &MoodResponse{
		ID:   m.ID,
		Name: m.Name,
	}
}

func toJournalResponse(j *JournalWithMood) JournalResponse {
	if j == nil {
		return JournalResponse{}
	}
	return JournalResponse{
		ID:          j.ID,
		Title:       j.Title,
		Description: j.Description,
		Status:      j.Status,
		ImageURL:    j.ImageURL,
		CreatedAt:   j.CreatedAt,
		UpdatedAt:   j.UpdatedAt,
		Mood:        toMoodResponse(j.Mood),
	}
}
