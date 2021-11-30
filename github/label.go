package github

import (
	"github.com/google/go-github/v37/github"
)

type Labels struct {
	QueueLabels    []*github.Label
	DependentLabel *github.Label
	ExcludeLabel   *github.Label
}

func CreateLabels(l []string, labelColor string, queueName string) []*github.Label {
	labels := make([]*github.Label, len(l))
	for i := 0; i < len(l); i++ {
		labels[i] = CreateLabel(l[i], labelColor, queueName)
	}
	return labels
}

func CreateLabel(l string, labelColor string, queueName string) *github.Label {
	label := queueName + ": " + l
	return &github.Label{ID: nil, URL: nil, Name: &label, Color: &labelColor, Description: &label, Default: nil, NodeID: nil}
}

// boolean return if a label exists in a list of labels
func ContainsLabelWithName(s []*github.Label, e *github.Label) bool {
	for _, a := range s {
		if *a.Name == *e.Name {
			return true
		}
	}
	return false
}
