package srospan

import "go.opentelemetry.io/otel/attribute"

func SourceOwnerId(val string) attribute.KeyValue {
	return attribute.String("sro.source.owner.id", val)
}

func SourceOwnerUsername(val string) attribute.KeyValue {
	return attribute.String("sro.source.owner.username", val)
}

func SourceCharacterId(val int) attribute.KeyValue {
	return attribute.Int("sro.source.character.id", val)
}

func SourceCharacterName(val string) attribute.KeyValue {
	return attribute.String("sro.source.character.username", val)
}

func TargetOwnerId(val string) attribute.KeyValue {
	return attribute.String("sro.target.owner.id", val)
}

func TargetOwnerUsername(val string) attribute.KeyValue {
	return attribute.String("sro.target.owner.username", val)
}

func TargetCharacterId(val int) attribute.KeyValue {
	return attribute.Int("sro.target.character.id", val)
}

func TargetCharacterName(val string) attribute.KeyValue {
	return attribute.String("sro.target.character.username", val)
}
